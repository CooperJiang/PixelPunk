package websocket

import (
	"sync"
	"time"

	"pixelpunk/pkg/logger"
)

// Manager WebSocket连接管理器
type Manager struct {
	clients    map[string]*Client
	clientsMux sync.RWMutex

	register   chan *Client
	unregister chan *Client

	broadcast chan *Message

	stopChan chan struct{}

	stats *Stats

	config *Config
}

// Config 配置
type Config struct {
	MaxConnections  int           // 最大连接数
	PingInterval    time.Duration // ping间隔
	PongTimeout     time.Duration // pong超时
	WriteTimeout    time.Duration // 写超时
	ReadTimeout     time.Duration // 读超时
	MaxMessageSize  int64         // 最大消息大小
	SendBufferSize  int           // 发送缓冲区大小
	RateLimitPerSec int           // 每秒消息速率限制
}

// Stats 统计信息
type Stats struct {
	ActiveConnections int64 // 活跃连接数
	TotalMessages     int64 // 总消息数
	MessagePerSecond  int64 // 每秒消息数
	mutex             sync.RWMutex
}

func NewManager(config *Config) *Manager {
	if config == nil {
		config = DefaultConfig()
	}

	return &Manager{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message, 1000),
		stopChan:   make(chan struct{}),
		stats:      &Stats{},
		config:     config,
	}
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		MaxConnections:  1000,
		PingInterval:    30 * time.Second,
		PongTimeout:     60 * time.Second,
		WriteTimeout:    10 * time.Second,
		ReadTimeout:     60 * time.Second,
		MaxMessageSize:  1024 * 1024, // 1MB
		SendBufferSize:  256,
		RateLimitPerSec: 20,
	}
}

// Start 启动管理器
func (m *Manager) Start() {
	go m.run()
	go m.pingClients()
	go m.cleanupClients()
}

// Stop 停止管理器
func (m *Manager) Stop() {
	close(m.stopChan)

	m.clientsMux.RLock()
	for _, client := range m.clients {
		client.Close()
	}
	m.clientsMux.RUnlock()
}

// RegisterClient 注册客户端
func (m *Manager) RegisterClient(client *Client) {
	select {
	case m.register <- client:
	case <-m.stopChan:
	}
}

// UnregisterClient 注销客户端
func (m *Manager) UnregisterClient(client *Client) {
	select {
	case m.unregister <- client:
	case <-m.stopChan:
	}
}

// BroadcastMessage 广播消息
func (m *Manager) BroadcastMessage(msg *Message) {
	select {
	case m.broadcast <- msg:
	case <-m.stopChan:
	}
}

// SendToClient 发送消息给指定客户端
func (m *Manager) SendToClient(clientID string, msg *Message) error {
	m.clientsMux.RLock()
	client, exists := m.clients[clientID]
	m.clientsMux.RUnlock()

	if !exists {
		return ErrClientNotFound
	}

	return client.SendMessage(msg)
}

// SendToAdmins 发送消息给所有管理员
func (m *Manager) SendToAdmins(msg *Message) {
	m.clientsMux.RLock()
	defer m.clientsMux.RUnlock()

	for _, client := range m.clients {
		if client.IsAdmin && client.IsConnected() {
			go func(c *Client) {
				if err := c.SendMessage(msg); err != nil {
					logger.Warn("发送消息给管理员失败: %v", err)
				}
			}(client)
		}
	}
}

func (m *Manager) GetStats() *Stats {
	m.stats.mutex.RLock()
	defer m.stats.mutex.RUnlock()

	return &Stats{
		ActiveConnections: m.stats.ActiveConnections,
		TotalMessages:     m.stats.TotalMessages,
		MessagePerSecond:  m.stats.MessagePerSecond,
	}
}

// run 主运行循环
func (m *Manager) run() {
	for {
		select {
		case client := <-m.register:
			m.handleRegister(client)

		case client := <-m.unregister:
			m.handleUnregister(client)

		case message := <-m.broadcast:
			m.handleBroadcast(message)

		case <-m.stopChan:
			return
		}
	}
}

// handleRegister 处理客户端注册
func (m *Manager) handleRegister(client *Client) {
	m.clientsMux.RLock()
	currentConnections := len(m.clients)
	m.clientsMux.RUnlock()

	if currentConnections >= m.config.MaxConnections {
		logger.Warn("连接数已达上限，拒绝新连接: %s", client.ID)
		client.Close()
		return
	}

	m.clientsMux.Lock()
	m.clients[client.ID] = client
	m.clientsMux.Unlock()

	m.stats.mutex.Lock()
	m.stats.ActiveConnections++
	m.stats.mutex.Unlock()

	welcomeMsg := NewMessage(MessageTypeSystemStatus, map[string]interface{}{
		"status":      "connected",
		"client_id":   client.ID,
		"server_time": time.Now().Unix(),
	})
	client.SendMessage(welcomeMsg)
}

// handleUnregister 处理客户端注销
func (m *Manager) handleUnregister(client *Client) {
	m.clientsMux.Lock()
	if _, exists := m.clients[client.ID]; exists {
		delete(m.clients, client.ID)
		client.Close()
	}
	m.clientsMux.Unlock()

	m.stats.mutex.Lock()
	m.stats.ActiveConnections--
	m.stats.mutex.Unlock()

}

// handleBroadcast 处理广播消息
func (m *Manager) handleBroadcast(msg *Message) {
	m.clientsMux.RLock()
	clients := make([]*Client, 0, len(m.clients))
	for _, client := range m.clients {
		if client.IsConnected() {
			clients = append(clients, client)
		}
	}
	m.clientsMux.RUnlock()

	for _, client := range clients {
		go func(c *Client) {
			if err := c.SendMessage(msg); err != nil {
				logger.Warn("广播消息失败: %v", err)
			}
		}(client)
	}

	m.stats.mutex.Lock()
	m.stats.TotalMessages++
	m.stats.mutex.Unlock()
}

// pingClients 定期ping客户端
func (m *Manager) pingClients() {
	ticker := time.NewTicker(m.config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.sendPingToAll()
		case <-m.stopChan:
			return
		}
	}
}

// sendPingToAll 向所有客户端发送ping
func (m *Manager) sendPingToAll() {
	pingMsg := NewMessage(MessageTypePing, map[string]interface{}{
		"timestamp": time.Now().Unix(),
	})

	m.clientsMux.RLock()
	for _, client := range m.clients {
		if client.IsConnected() {
			go func(c *Client) {
				if err := c.SendMessage(pingMsg); err != nil {
					m.UnregisterClient(c)
				}
			}(client)
		}
	}
	m.clientsMux.RUnlock()
}

// cleanupClients 清理断开的客户端
func (m *Manager) cleanupClients() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.doCleanup()
		case <-m.stopChan:
			return
		}
	}
}

// doCleanup 执行清理
func (m *Manager) doCleanup() {
	now := time.Now()
	toRemove := make([]*Client, 0)

	m.clientsMux.RLock()
	for _, client := range m.clients {
		// 检查最后ping时间，超时的客户端标记为删除
		if now.Sub(client.GetLastPing()) > m.config.PongTimeout {
			toRemove = append(toRemove, client)
		}
	}
	m.clientsMux.RUnlock()

	for _, client := range toRemove {
		m.UnregisterClient(client)
	}
}
