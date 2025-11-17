package queue

import (
	"time"
)

// TaggingTask 表示一条AI打标任务（最小信息：文件ID）
type TaggingTask struct {
	FileID string
	// 可扩展字段：优先级、重试次数等
}

// MessageType 支持多队列类型
type MessageType string

const (
	MessageTypeAITagging MessageType = "ai_tagging"
	MessageTypeVector    MessageType = "vector"
)

// Metrics 队列运行时指标（用于WS推送与监控）
type Metrics struct {
	QueueLength  int
	InFlight     int
	DelayedCount int
	DLQCount     int

	ActiveWorkers         int
	ConfiguredConcurrency int
	Paused                bool
	ProcessingRatePerMin  float64
}

// Queue 通用队列接口（支持 Redis/DB 实现）
type Queue interface {
	// EnqueueUnique 幂等入队：同一 fileID 不重复入队
	EnqueueUnique(fileID string, priority int) error

	// Fetch 取出一条任务并获得租约 lease（可见性超时/过期后可被再取）
	Fetch(lease time.Duration) (*TaggingTask, AckFunc, NackFunc, error)

	// Metrics 获取队列指标（实现可内部采样/缓存）
	Metrics() (*Metrics, error)

	// Close 关闭队列资源
	Close() error
}

// AckFunc 确认完成（成功）
type AckFunc func() error

// NackFunc 失败处理：可选择延迟重试或丢到DLQ
// delay <= 0 表示立即重试；toDLQ=true 表示进入死信，不再自动重试
type NackFunc func(delay time.Duration, toDLQ bool, lastError string) error
