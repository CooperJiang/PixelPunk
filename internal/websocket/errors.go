package websocket

import "errors"

// WebSocket相关错误定义
var (
	ErrClientDisconnected = errors.New("客户端已断开连接")
	ErrSendBufferFull     = errors.New("发送缓冲区已满")
	ErrInvalidMessage     = errors.New("无效的消息格式")
	ErrUnauthorized       = errors.New("未授权的连接")
	ErrRateLimitExceeded  = errors.New("消息频率超限")
	ErrClientNotFound     = errors.New("客户端不存在")
)
