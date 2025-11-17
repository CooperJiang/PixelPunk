package cache

import (
	"context"
	"errors"
	"fmt"
	"pixelpunk/pkg/config"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisCache *RedisCache

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

// InitRedis 初始化Redis缓存
func InitRedis() error {
	cfg := config.GetConfig().Redis

	// 如果未配置Redis主机，则不使用Redis
	if cfg.Host == "" {
		return errors.New("Redis未配置")
	}

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	options := &redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	client := redis.NewClient(options)

	ctx := context.Background()

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctxWithTimeout).Result()
	if err != nil {
		client.Close() // 关闭失败的连接
		return err
	}

	redisCache = &RedisCache{
		client: client,
		ctx:    ctx,
	}

	defaultCache = redisCache

	return nil
}

// Set 设置缓存
func (c *RedisCache) Set(key string, value string, expiration time.Duration) error {
	return c.client.Set(c.ctx, key, value, expiration).Err()
}

// Get 获取缓存
func (c *RedisCache) Get(key string) (string, error) {
	return c.client.Get(c.ctx, key).Result()
}

// Del 删除缓存
func (c *RedisCache) Del(key string) error {
	return c.client.Del(c.ctx, key).Err()
}

// Exists 检查键是否存在
func (c *RedisCache) Exists(key string) bool {
	result, _ := c.client.Exists(c.ctx, key).Result()
	return result > 0
}

// TTL 获取过期时间
func (c *RedisCache) TTL(key string) (time.Duration, error) {
	return c.client.TTL(c.ctx, key).Result()
}

// Expire 设置过期时间
func (c *RedisCache) Expire(key string, expiration time.Duration) error {
	return c.client.Expire(c.ctx, key, expiration).Err()
}

// Close 关闭Redis连接
func (c *RedisCache) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}

// RedisAvailable 检查Redis是否可用
func RedisAvailable() bool {
	return redisCache != nil && redisCache.client != nil
}

// 性能监控和指标收集使用
func GetRedisClient() *redis.Client {
	if !RedisAvailable() {
		return nil
	}
	return redisCache.client
}

func GetRedisContext() context.Context {
	if !RedisAvailable() {
		return context.Background()
	}
	return redisCache.ctx
}
