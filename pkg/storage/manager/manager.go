package manager

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"sync"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/storage/factory"
)

// ChannelRepository 渠道仓库接口
// 用于从数据库获取存储渠道配置信息
type ChannelRepository interface {
	GetChannel(channelID string) (*models.StorageChannel, error)
	GetChannelConfig(channelID string) (map[string]interface{}, error)
	GetDefaultChannel() (*models.StorageChannel, error)
	GetActiveChannels() ([]*models.StorageChannel, error)
}

// StorageManager 存储管理器
// 负责管理多个存储适配器实例和路由存储操作
type StorageManager struct {
	factory      *factory.StorageFactory
	adapters     map[string]adapter.StorageAdapter
	channelRepo  ChannelRepository
	mu           sync.RWMutex
	cacheTimeout time.Duration
	lastUpdate   map[string]time.Time
}

func NewStorageManager(channelRepo ChannelRepository) *StorageManager {
	return &StorageManager{
		factory:      factory.NewStorageFactory(),
		adapters:     make(map[string]adapter.StorageAdapter),
		channelRepo:  channelRepo,
		cacheTimeout: 5 * time.Minute, // 适配器缓存5分钟
		lastUpdate:   make(map[string]time.Time),
	}
}

// NewStorageManagerWithFactory 使用指定工厂创建存储管理器
func NewStorageManagerWithFactory(channelRepo ChannelRepository, f *factory.StorageFactory) *StorageManager {
	return &StorageManager{
		factory:      f,
		adapters:     make(map[string]adapter.StorageAdapter),
		channelRepo:  channelRepo,
		cacheTimeout: 5 * time.Minute,
		lastUpdate:   make(map[string]time.Time),
	}
}

func (m *StorageManager) GetAdapter(channelID string) (adapter.StorageAdapter, error) {
	m.mu.RLock()

	if adapterInstance, exists := m.adapters[channelID]; exists {
		if lastUpdate, ok := m.lastUpdate[channelID]; ok {
			if time.Since(lastUpdate) < m.cacheTimeout {
				m.mu.RUnlock()
				return adapterInstance, nil
			}
		}
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if adapterInstance, exists := m.adapters[channelID]; exists {
		if lastUpdate, ok := m.lastUpdate[channelID]; ok {
			if time.Since(lastUpdate) < m.cacheTimeout {
				return adapterInstance, nil
			}
		}
	}

	return m.createAdapter(channelID)
}

func (m *StorageManager) GetDefaultAdapter() (adapter.StorageAdapter, error) {
	channel, err := m.channelRepo.GetDefaultChannel()
	if err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			"failed to get default channel",
			err,
		)
	}

	return m.GetAdapter(channel.ID)
}

func (m *StorageManager) GetBestAdapter() (adapter.StorageAdapter, error) {
	channel, err := m.channelRepo.GetDefaultChannel()
	if err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			"failed to get best channel",
			err,
		)
	}

	return m.GetAdapter(channel.ID)
}

func (m *StorageManager) GetDefaultChannelID() (string, error) {
	channel, err := m.channelRepo.GetDefaultChannel()
	if err != nil {
		return "", fmt.Errorf("failed to get default channel: %w", err)
	}

	return channel.ID, nil
}

// Upload 上传文件到指定渠道
func (m *StorageManager) Upload(ctx context.Context, channelID string, req *adapter.UploadRequest) (*adapter.UploadResult, error) {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}

	return adapterInstance.Upload(ctx, req)
}

// UploadWithDefault 使用默认适配器上传
func (m *StorageManager) UploadWithDefault(ctx context.Context, req *adapter.UploadRequest) (*adapter.UploadResult, error) {
	adapterInstance, err := m.GetDefaultAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to get default adapter: %w", err)
	}

	return adapterInstance.Upload(ctx, req)
}

// UploadWithBest 使用最佳适配器上传
func (m *StorageManager) UploadWithBest(ctx context.Context, req *adapter.UploadRequest) (*adapter.UploadResult, error) {
	adapterInstance, err := m.GetBestAdapter()
	if err != nil {
		return nil, fmt.Errorf("failed to get best adapter: %w", err)
	}

	return adapterInstance.Upload(ctx, req)
}

// Delete 删除文件
func (m *StorageManager) Delete(ctx context.Context, channelID, path string) error {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}

	return adapterInstance.Delete(ctx, path)
}

func (m *StorageManager) GetURL(channelID, path string, options *adapter.URLOptions) (string, error) {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return "", fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}

	return adapterInstance.GetURL(path, options)
}

// HealthCheck 健康检查
func (m *StorageManager) HealthCheck(ctx context.Context, channelID string) error {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}

	return adapterInstance.HealthCheck(ctx)
}

// HealthCheckAll 检查所有活跃渠道的健康状态
func (m *StorageManager) HealthCheckAll(ctx context.Context) map[string]error {
	channels, err := m.channelRepo.GetActiveChannels()
	if err != nil {
		return map[string]error{
			"repository": fmt.Errorf("failed to get active channels: %w", err),
		}
	}

	results := make(map[string]error)
	for _, channel := range channels {
		if err := m.HealthCheck(ctx, channel.ID); err != nil {
			results[channel.ID] = err
		} else {
			results[channel.ID] = nil // 成功
		}
	}

	return results
}

// RefreshAdapter 刷新适配器（强制重新创建）
func (m *StorageManager) RefreshAdapter(channelID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.adapters, channelID)
	delete(m.lastUpdate, channelID)

	_, err := m.createAdapter(channelID)
	return err
}

func (m *StorageManager) RemoveAdapter(channelID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.adapters, channelID)
	delete(m.lastUpdate, channelID)
}

func (m *StorageManager) GetChannelType(channelID string) (string, error) {
	channel, err := m.channelRepo.GetChannel(channelID)
	if err != nil {
		return "", fmt.Errorf("failed to get channel: %w", err)
	}

	return channel.Type, nil
}

func (m *StorageManager) GetAdapterCapabilities(channelID string) (adapter.Capabilities, error) {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return adapter.Capabilities{}, fmt.Errorf("failed to get adapter: %w", err)
	}

	return adapterInstance.GetCapabilities(), nil
}

// createAdapter 创建适配器（内部方法，调用时需要持有锁）
func (m *StorageManager) createAdapter(channelID string) (adapter.StorageAdapter, error) {
	channel, err := m.channelRepo.GetChannel(channelID)
	if err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to get channel %s", channelID),
			err,
		)
	}

	config, err := m.channelRepo.GetChannelConfig(channelID)
	if err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to get channel config for %s", channelID),
			err,
		)
	}

	adapterInstance, err := m.factory.CreateAdapter(channel.Type, config)
	if err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to create adapter for channel %s", channelID),
			err,
		)
	}

	m.adapters[channelID] = adapterInstance
	m.lastUpdate[channelID] = time.Now()

	return adapterInstance, nil
}

func (m *StorageManager) SetCacheTimeout(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cacheTimeout = timeout
}

// ClearCache 清空所有缓存
func (m *StorageManager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.adapters = make(map[string]adapter.StorageAdapter)
	m.lastUpdate = make(map[string]time.Time)
}

func (m *StorageManager) GetCachedAdapterCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.adapters)
}

// GetBase64 获取文件的Base64编码
func (m *StorageManager) GetBase64(ctx context.Context, channelID, path string) (string, error) {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return "", fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}
	rc, err := adapterInstance.ReadFile(ctx, path)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// GetThumbnailBase64 获取缩略图的Base64编码
func (m *StorageManager) GetThumbnailBase64(ctx context.Context, channelID, path string) (string, error) {
	// 统一实现：直接读取给定路径（应为对象键，且指向缩略图）
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return "", fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}
	rc, err := adapterInstance.ReadFile(ctx, path)
	if err != nil {
		return "", err
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// ReadFile 读取文件内容
func (m *StorageManager) ReadFile(ctx context.Context, channelID, path string) (io.ReadCloser, error) {
	adapterInstance, err := m.GetAdapter(channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get adapter for channel %s: %w", channelID, err)
	}

	return adapterInstance.ReadFile(ctx, path)
}
