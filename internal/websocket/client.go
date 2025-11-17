package websocket

import (
	"encoding/json"
	"sync"
	"time"

	"pixelpunk/pkg/logger"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端连接
type Client struct {
	// WebSocket连接
	conn *websocket.Conn

	send chan []byte

	ID      string // 客户端唯一ID
	UserID  uint   // 用户ID
	IsAdmin bool   // 是否为管理员

	isConnected bool
	lastPing    time.Time
	mutex       sync.RWMutex

	pendingAcks map[string]*Message // 等待确认的消息
	ackMutex    sync.RWMutex
}

func NewClient(conn *websocket.Conn, userID uint, isAdmin bool) *Client {
	return &Client{
		conn:        conn,
		send:        make(chan []byte, 256),
		ID:          generateClientID(userID),
		UserID:      userID,
		IsAdmin:     isAdmin,
		isConnected: true,
		lastPing:    time.Now(),
		pendingAcks: make(map[string]*Message),
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(msg *Message) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return ErrClientDisconnected
	}

	// 如果需要确认，添加到待确认列表
	if msg.RequireAck {
		c.ackMutex.Lock()
		c.pendingAcks[msg.ID] = msg
		c.ackMutex.Unlock()
	}

	data, err := msg.ToJSON()
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
		return nil
	default:
		// 发送缓冲区满，断开连接
		c.Close()
		return ErrSendBufferFull
	}
}

// SendJSON 发送JSON数据
func (c *Client) SendJSON(data interface{}) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if !c.isConnected {
		return ErrClientDisconnected
	}

	return c.conn.WriteJSON(data)
}

// Close 关闭连接
func (c *Client) Close() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.isConnected {
		c.isConnected = false
		close(c.send)
		c.conn.Close()
	}
}

// IsConnected 检查连接状态
func (c *Client) IsConnected() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isConnected
}

func (c *Client) UpdatePing() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastPing = time.Now()
}

func (c *Client) GetLastPing() time.Time {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.lastPing
}

// HandleAck 处理消息确认
func (c *Client) HandleAck(messageID string) {
	c.ackMutex.Lock()
	defer c.ackMutex.Unlock()
	delete(c.pendingAcks, messageID)
}

func (c *Client) GetPendingAcks() []*Message {
	c.ackMutex.RLock()
	defer c.ackMutex.RUnlock()

	messages := make([]*Message, 0, len(c.pendingAcks))
	for _, msg := range c.pendingAcks {
		messages = append(messages, msg)
	}
	return messages
}

// ReadPump 读取消息循环
func (c *Client) ReadPump(manager *Manager) {
	defer func() {
		manager.UnregisterClient(c)
		c.Close()
	}()

	c.conn.SetReadLimit(manager.config.MaxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(manager.config.ReadTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.UpdatePing()
		c.conn.SetReadDeadline(time.Now().Add(manager.config.ReadTimeout))
		return nil
	})

	for {
		_, messageData, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Error("WebSocket读取错误: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(messageData, &msg); err != nil {
			logger.Warn("解析WebSocket消息失败: %v", err)
			continue
		}

		c.handleMessage(&msg)
	}
}

// WritePump 写入消息循环
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				logger.Error("WebSocket写入错误: %v", err)
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				logger.Error("WebSocket ping失败: %v", err)
				return
			}
		}
	}
}

// handleMessage 处理接收到的消息
func (c *Client) handleMessage(msg *Message) {
	switch msg.Type {
	case MessageTypePong:
		c.UpdatePing()

	case MessageTypePing:
		pongMsg := NewMessage(MessageTypePong, map[string]interface{}{
			"timestamp": time.Now().Unix(),
		})
		c.SendMessage(pongMsg)

	default:
		// 其他消息类型暂时忽略
	}
}

// generateClientID 生成客户端ID
func generateClientID(userID uint) string {
	return "client_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}
