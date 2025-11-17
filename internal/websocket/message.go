package websocket

import (
	"encoding/json"
	"time"
)

// MessageType 消息类型
type MessageType string

const (
	// 消息类型常量
	MessageTypeQueueStats   MessageType = "queue_stats"
	MessageTypeVectorStats  MessageType = "vector_stats"
	MessageTypeLogs         MessageType = "logs"
	MessageTypeAnnouncement MessageType = "announcement"
	MessageTypeSystemStatus MessageType = "system_status"
	MessageTypeError        MessageType = "error"
	MessageTypePing         MessageType = "ping"
	MessageTypePong         MessageType = "pong"
)

// MessagePriority 消息优先级
type MessagePriority string

const (
	PriorityHigh   MessagePriority = "high"
	PriorityNormal MessagePriority = "normal"
	PriorityLow    MessagePriority = "low"
)

// Message WebSocket消息结构
type Message struct {
	ID         string          `json:"id"`                    // 消息唯一ID
	Type       MessageType     `json:"type"`                  // 消息类型
	Priority   MessagePriority `json:"priority"`              // 优先级
	Timestamp  int64           `json:"timestamp"`             // 时间戳
	Source     string          `json:"source,omitempty"`      // 数据源
	Data       interface{}     `json:"data,omitempty"`        // 消息内容
	RequireAck bool            `json:"require_ack,omitempty"` // 是否需要确认
}

func NewMessage(msgType MessageType, data interface{}) *Message {
	return &Message{
		ID:        generateMessageID(),
		Type:      msgType,
		Priority:  PriorityNormal,
		Timestamp: time.Now().Unix(),
		Data:      data,
	}
}

func NewHighPriorityMessage(msgType MessageType, data interface{}) *Message {
	msg := NewMessage(msgType, data)
	msg.Priority = PriorityHigh
	msg.RequireAck = true
	return msg
}

// ToJSON 转换为JSON
func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// FromJSON 从JSON解析
func FromJSON(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}

// generateMessageID 生成消息ID
func generateMessageID() string {
	return "msg_" + time.Now().Format("20060102150405") + "_" + randomString(6)
}

// randomString 生成随机字符串
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
