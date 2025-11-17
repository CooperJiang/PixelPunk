package queue

import (
	"context"
	"fmt"
	"time"

	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"github.com/redis/go-redis/v9"
)

// RedisQueue 使用 Redis List + Lua 原子脚本实现
type RedisQueue struct {
	cli *redis.Client
	ctx context.Context
	pfx string // key 前缀：<app.ns>:<kind>

	kQueue      string // 列表：主队列
	kProcessing string // 列表：处理中
	kProcZ      string // zset：处理中租约过期时间
	kDelayedZ   string // zset：延迟重试
	kDLQ        string // 列表：死信
	kEnqueued   string // set：去重
}

func NewRedisQueue() *RedisQueue {
	rc := cache.GetRedisClient()
	if rc == nil {
		return nil
	}
	ctx := cache.GetRedisContext()
	ns := cache.GetNamespace()
	// 默认用于打标队列；向量队列在外层设置 pfx
	pfx := fmt.Sprintf("%s:%s", ns, "ai:tagging")
	q := &RedisQueue{
		cli: rc,
		ctx: ctx,
		pfx: pfx,
	}
	q.kQueue = pfx + ":queue"
	q.kProcessing = pfx + ":processing:list"
	q.kProcZ = pfx + ":processing:z"
	q.kDelayedZ = pfx + ":delayed"
	q.kDLQ = pfx + ":dlq"
	q.kEnqueued = pfx + ":enqueued"
	return q
}

// WithPrefix 允许自定义业务前缀（例如 vector）
func (q *RedisQueue) WithPrefix(kind string) *RedisQueue {
	q.pfx = fmt.Sprintf("%s:%s", cache.GetNamespace(), kind)
	q.kQueue = q.pfx + ":queue"
	q.kProcessing = q.pfx + ":processing:list"
	q.kProcZ = q.pfx + ":processing:z"
	q.kDelayedZ = q.pfx + ":delayed"
	q.kDLQ = q.pfx + ":dlq"
	q.kEnqueued = q.pfx + ":enqueued"
	return q
}

func (q *RedisQueue) keyExists() bool { return q.cli != nil }

// EnqueueUnique：SADD去重成功才LPUSH（改进版：检查心跳避免重复入队）
func (q *RedisQueue) EnqueueUnique(fileID string, priority int) error {
	if !q.keyExists() {
		return fmt.Errorf("redis not available")
	}

	// 先检查File表心跳时间，避免重复入队正在处理的任务
	var file struct {
		AITaggingStatus      string
		AILastHeartbeatAt    *time.Time
		AIProcessingWorkerID string
	}
	db := database.GetDB()
	if db != nil {
		if err := db.Table("file").
			Select("ai_tagging_status, ai_last_heartbeat_at, ai_processing_worker_id").
			Where("id = ?", fileID).
			Take(&file).Error; err == nil {
			// 如果文件正在处理中（pending + 心跳时间在2分钟内），跳过入队
			if file.AITaggingStatus == "pending" &&
				file.AILastHeartbeatAt != nil &&
				time.Since(*file.AILastHeartbeatAt) < 2*time.Minute {
				return nil
			}
		}
	}

	added, err := q.cli.SAdd(q.ctx, q.kEnqueued, fileID).Result()
	if err != nil {
		return err
	}
	if added == 1 {
		return q.cli.LPush(q.ctx, q.kQueue, fileID).Err()
	}
	return nil
}

// Lua：原子 RPOPLPUSH + ZADD(lease)
var luaFetch = redis.NewScript(`
local q = KEYS[1]
local processingList = KEYS[2]
local processingZ = KEYS[3]
local leaseTs = ARGV[1]
local v = redis.call('RPOPLPUSH', q, processingList)
if v then
  redis.call('ZADD', processingZ, leaseTs, v)
end
return v
`)

func (q *RedisQueue) Fetch(lease time.Duration) (*TaggingTask, AckFunc, NackFunc, error) {
	if !q.keyExists() {
		return nil, nil, nil, fmt.Errorf("redis not available")
	}
	leaseUntil := time.Now().Add(lease).Unix()
	v, err := luaFetch.Run(q.ctx, q.cli, []string{q.kQueue, q.kProcessing, q.kProcZ}, leaseUntil).Result()
	if err != nil && err != redis.Nil {
		return nil, nil, nil, err
	}
	if v == nil {
		return nil, nil, nil, redis.Nil
	}
	id, _ := v.(string)
	task := &TaggingTask{FileID: id}

	// Ack：LREM processing:list + ZREM processing:z + SREM enqueued
	ack := func() error {
		pipe := q.cli.TxPipeline()
		pipe.LRem(q.ctx, q.kProcessing, 0, id)
		pipe.ZRem(q.ctx, q.kProcZ, id)
		pipe.SRem(q.ctx, q.kEnqueued, id)
		_, e := pipe.Exec(q.ctx)
		return e
	}

	// Nack：从processing移除；toDLQ则LPUSH到DLQ；否则ZADD delayed
	nack := func(delay time.Duration, toDLQ bool, lastError string) error {
		pipe := q.cli.TxPipeline()
		pipe.LRem(q.ctx, q.kProcessing, 0, id)
		pipe.ZRem(q.ctx, q.kProcZ, id)
		if toDLQ {
			pipe.LPush(q.ctx, q.kDLQ, id)
			// 不立即从去重集合删除，保留不再自动进入主队列
		} else {
			when := time.Now().Add(delay).Unix()
			pipe.ZAdd(q.ctx, q.kDelayedZ, redis.Z{Score: float64(when), Member: id})
		}
		_, e := pipe.Exec(q.ctx)
		if e != nil {
			return e
		}
		return nil
	}

	return task, ack, nack, nil
}

// 由外部定时器调用：搬运到期的 delayed 与超时 processing 到主队列
func (q *RedisQueue) ReapOnce(now time.Time) error {
	if !q.keyExists() {
		return nil
	}
	for i := 0; i < 100; i++ {
		zs, err := q.cli.ZRangeByScoreWithScores(q.ctx, q.kDelayedZ, &redis.ZRangeBy{Min: "-inf", Max: fmt.Sprintf("%d", now.Unix()), Count: 100}).Result()
		if err != nil {
			return err
		}
		if len(zs) == 0 {
			break
		}
		pipe := q.cli.TxPipeline()
		for _, z := range zs {
			id := z.Member.(string)
			pipe.ZRem(q.ctx, q.kDelayedZ, id)
			pipe.LPush(q.ctx, q.kQueue, id)
		}
		if _, err := pipe.Exec(q.ctx); err != nil {
			return err
		}
	}
	for i := 0; i < 100; i++ {
		zs, err := q.cli.ZRangeByScoreWithScores(q.ctx, q.kProcZ, &redis.ZRangeBy{Min: "-inf", Max: fmt.Sprintf("%d", now.Unix()), Count: 100}).Result()
		if err != nil {
			return err
		}
		if len(zs) == 0 {
			break
		}
		pipe := q.cli.TxPipeline()
		for _, z := range zs {
			id := z.Member.(string)
			pipe.ZRem(q.ctx, q.kProcZ, id)
			pipe.LRem(q.ctx, q.kProcessing, 0, id)
			pipe.LPush(q.ctx, q.kQueue, id)
		}
		if _, err := pipe.Exec(q.ctx); err != nil {
			return err
		}
	}
	return nil
}

func (q *RedisQueue) Metrics() (*Metrics, error) {
	if !q.keyExists() {
		return nil, fmt.Errorf("redis not available")
	}
	pipe := q.cli.Pipeline()
	ql := pipe.LLen(q.ctx, q.kQueue)
	infl := pipe.LLen(q.ctx, q.kProcessing)
	delayed := pipe.ZCard(q.ctx, q.kDelayedZ)
	dlq := pipe.LLen(q.ctx, q.kDLQ)
	if _, err := pipe.Exec(q.ctx); err != nil {
		return nil, err
	}
	return &Metrics{QueueLength: int(ql.Val()), InFlight: int(infl.Val()), DelayedCount: int(delayed.Val()), DLQCount: int(dlq.Val())}, nil
}

func (q *RedisQueue) Close() error { return nil }

// StartReaper 启动后台reaper（可由服务方管理其生命周期）
func (q *RedisQueue) StartReaper(interval time.Duration, stop <-chan struct{}) {
	if !q.keyExists() {
		return
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := q.ReapOnce(time.Now()); err != nil {
					logger.Warn("redis queue reap failed: %v", err)
				}
			case <-stop:
				return
			}
		}
	}()
}
