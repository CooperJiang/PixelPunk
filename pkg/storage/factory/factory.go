package factory

import (
	"fmt"
	"sync"

	"pixelpunk/pkg/storage/adapter"
	"pixelpunk/pkg/storage/config"
)

// StorageFactory 存储工厂
type StorageFactory struct {
	adapters map[string]adapter.AdapterFactory
	mu       sync.RWMutex
}

func NewStorageFactory() *StorageFactory {
	return &StorageFactory{
		adapters: make(map[string]adapter.AdapterFactory),
	}
}

// RegisterAdapter 注册适配器
func (f *StorageFactory) RegisterAdapter(storageType string, factory adapter.AdapterFactory) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.adapters[storageType] = factory
}

// CreateAdapter 创建适配器实例
func (f *StorageFactory) CreateAdapter(storageType string, configData map[string]interface{}) (adapter.StorageAdapter, error) {
	f.mu.RLock()
	factory, exists := f.adapters[storageType]
	f.mu.RUnlock()

	if !exists {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("unsupported storage type: %s", storageType),
			nil,
		)
	}

	adapterInstance := factory()
	if adapterInstance == nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to create adapter instance for type: %s", storageType),
			nil,
		)
	}

	if err := adapterInstance.Initialize(configData); err != nil {
		return nil, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to initialize %s adapter", storageType),
			err,
		)
	}

	return adapterInstance, nil
}

func (f *StorageFactory) GetSupportedTypes() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var types []string
	for storageType := range f.adapters {
		types = append(types, storageType)
	}
	return types
}

// HasAdapter 检查是否支持指定的存储类型
func (f *StorageFactory) HasAdapter(storageType string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	_, exists := f.adapters[storageType]
	return exists
}

// UnregisterAdapter 注销适配器
func (f *StorageFactory) UnregisterAdapter(storageType string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.adapters, storageType)
}

// ValidateConfig 验证配置
func (f *StorageFactory) ValidateConfig(storageType string, configData map[string]interface{}) error {
	f.mu.RLock()
	factory, exists := f.adapters[storageType]
	f.mu.RUnlock()

	if !exists {
		return adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("unsupported storage type: %s", storageType),
			nil,
		)
	}

	tempAdapter := factory()
	if tempAdapter == nil {
		return adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to create temp adapter for validation: %s", storageType),
			nil,
		)
	}

	if err := tempAdapter.Initialize(configData); err != nil {
		return adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("invalid config for %s adapter", storageType),
			err,
		)
	}

	return nil
}

// CreateAdapterWithMapConfig 使用MapConfig创建适配器
func (f *StorageFactory) CreateAdapterWithMapConfig(storageType string, cfg *config.MapConfig) (adapter.StorageAdapter, error) {
	return f.CreateAdapter(storageType, cfg.GetAll())
}

func (f *StorageFactory) GetAdapterCapabilities(storageType string) (adapter.Capabilities, error) {
	f.mu.RLock()
	factory, exists := f.adapters[storageType]
	f.mu.RUnlock()

	if !exists {
		return adapter.Capabilities{}, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("unsupported storage type: %s", storageType),
			nil,
		)
	}

	tempAdapter := factory()
	if tempAdapter == nil {
		return adapter.Capabilities{}, adapter.NewStorageError(
			adapter.ErrorTypeInternal,
			fmt.Sprintf("failed to create temp adapter: %s", storageType),
			nil,
		)
	}

	return tempAdapter.GetCapabilities(), nil
}

// 全局工厂实例
var globalFactory = NewStorageFactory()

// RegisterGlobalAdapter 注册全局适配器
func RegisterGlobalAdapter(storageType string, factory adapter.AdapterFactory) {
	globalFactory.RegisterAdapter(storageType, factory)
}

// CreateGlobalAdapter 使用全局工厂创建适配器
func CreateGlobalAdapter(storageType string, configData map[string]interface{}) (adapter.StorageAdapter, error) {
	return globalFactory.CreateAdapter(storageType, configData)
}

func GetGlobalFactory() *StorageFactory {
	return globalFactory
}

func GetGlobalSupportedTypes() []string {
	return globalFactory.GetSupportedTypes()
}
