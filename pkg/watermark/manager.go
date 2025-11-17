package watermark

import (
	"encoding/json"
	"fmt"
	"image"
	"sync"
)

// Manager 水印管理器 - 提供高级水印处理功能
type Manager struct {
	processor       ProcessorInterface
	config          *ManagerConfig
	mu              sync.RWMutex
	processingStats *ProcessingStats
}

// ManagerConfig 管理器配置
type ManagerConfig struct {
	MaxConcurrentJobs   int    `json:"max_concurrent_jobs"`   // 最大并发处理任务数
	EnableStats         bool   `json:"enable_stats"`          // 是否启用统计
	WatermarkPath       string `json:"watermark_path"`        // 水印文件存储路径
	CacheEnabled        bool   `json:"cache_enabled"`         // 是否启用缓存
	MaxImageSize        int    `json:"max_image_size"`        // 最大文件尺寸
	EnableTextWatermark bool   `json:"enable_text_watermark"` // 是否启用文字水印
}

// ProcessingStats 处理统计信息
type ProcessingStats struct {
	TotalProcessed   int64 `json:"total_processed"`
	TotalSuccess     int64 `json:"total_success"`
	TotalFailed      int64 `json:"total_failed"`
	AverageTime      int64 `json:"average_time"`       // 平均处理时间（毫秒）
	TotalProcessTime int64 `json:"total_process_time"` // 总处理时间
}

func NewManager(config *ManagerConfig) *Manager {
	if config == nil {
		config = DefaultManagerConfig()
	}

	processor := NewProcessor()
	if config.WatermarkPath != "" {
		processor.SetWatermarkImagePath(config.WatermarkPath)
	}
	processor.SetEnableTextWatermark(config.EnableTextWatermark)

	return &Manager{
		processor:       processor,
		config:          config,
		processingStats: &ProcessingStats{},
	}
}

// DefaultManagerConfig 默认管理器配置
func DefaultManagerConfig() *ManagerConfig {
	return &ManagerConfig{
		MaxConcurrentJobs:   10,
		EnableStats:         true,
		WatermarkPath:       "internal/static", // 静态资源目录
		CacheEnabled:        false,
		MaxImageSize:        4096,
		EnableTextWatermark: false,
	}
}

// ProcessWithConfig 使用配置字符串处理图像
func (m *Manager) ProcessWithConfig(file image.Image, configJSON string) (*WatermarkResult, error) {
	if configJSON == "" {
		return &WatermarkResult{Success: true, ProcessedImage: file}, nil
	}

	config := &WatermarkConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	return m.Process(file, config)
}

// Process 处理单个图像
func (m *Manager) Process(file image.Image, config *WatermarkConfig) (*WatermarkResult, error) {
	if file == nil {
		return nil, fmt.Errorf("输入图像不能为空")
	}

	if config == nil || !config.Enabled {
		return &WatermarkResult{Success: true, ProcessedImage: file}, nil
	}

	if m.config.EnableStats {
		m.mu.Lock()
		m.processingStats.TotalProcessed++
		m.mu.Unlock()
	}

	processedImg, err := m.processor.ProcessImage(file, config)

	if m.config.EnableStats {
		m.updateStats(err == nil)
	}

	if err != nil {
		return &WatermarkResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	return &WatermarkResult{
		Success:        true,
		ProcessedImage: processedImg,
	}, nil
}

// ProcessBytes 处理图像字节数据
func (m *Manager) ProcessBytes(imgData []byte, config *WatermarkConfig) (*WatermarkResult, error) {
	if len(imgData) == 0 {
		return nil, fmt.Errorf("输入图像数据不能为空")
	}

	if config == nil || !config.Enabled {
		return &WatermarkResult{Success: true, ProcessedData: imgData}, nil
	}

	if m.config.EnableStats {
		m.mu.Lock()
		m.processingStats.TotalProcessed++
		m.mu.Unlock()
	}

	processedData, err := m.processor.ProcessImageBytes(imgData, config)

	if m.config.EnableStats {
		m.updateStats(err == nil)
	}

	if err != nil {
		return &WatermarkResult{
			Success:      false,
			ErrorMessage: err.Error(),
		}, err
	}

	return &WatermarkResult{
		Success:       true,
		ProcessedData: processedData,
	}, nil
}

// ProcessBytesWithConfig 使用配置字符串处理图像字节数据
func (m *Manager) ProcessBytesWithConfig(imgData []byte, configJSON string) (*WatermarkResult, error) {
	if len(imgData) == 0 {
		return nil, fmt.Errorf("输入图像数据不能为空")
	}

	if configJSON == "" {
		return &WatermarkResult{Success: true, ProcessedData: imgData}, nil
	}

	config := &WatermarkConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	// 无论前端还是后端模式，都需要合成水印
	// - frontend 模式：使用前端生成的水印图片（generatedFile）
	// - backend 模式：使用后端文件系统的水印图片（fileURL）
	return m.ProcessBytes(imgData, config)
}

// ValidateConfig 验证配置
func (m *Manager) ValidateConfig(config *WatermarkConfig) error {
	return m.processor.ValidateConfig(config)
}

// ValidateConfigJSON 验证JSON配置字符串
func (m *Manager) ValidateConfigJSON(configJSON string) error {
	if configJSON == "" {
		return nil // 空配置是有效的
	}

	config := &WatermarkConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return fmt.Errorf("JSON格式错误: %w", err)
	}

	return m.ValidateConfig(config)
}

func (m *Manager) GetStats() *ProcessingStats {
	if !m.config.EnableStats {
		return nil
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	return &ProcessingStats{
		TotalProcessed:   m.processingStats.TotalProcessed,
		TotalSuccess:     m.processingStats.TotalSuccess,
		TotalFailed:      m.processingStats.TotalFailed,
		AverageTime:      m.processingStats.AverageTime,
		TotalProcessTime: m.processingStats.TotalProcessTime,
	}
}

// ResetStats 重置统计信息
func (m *Manager) ResetStats() {
	if !m.config.EnableStats {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.processingStats = &ProcessingStats{}
}

// updateStats 更新统计信息
func (m *Manager) updateStats(success bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if success {
		m.processingStats.TotalSuccess++
	} else {
		m.processingStats.TotalFailed++
	}

	if m.processingStats.TotalProcessed > 0 {
		m.processingStats.AverageTime = m.processingStats.TotalProcessTime / m.processingStats.TotalProcessed
	}
}

func (m *Manager) SetConfig(config *ManagerConfig) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.config = config

	if config.WatermarkPath != "" {
		m.processor.SetWatermarkImagePath(config.WatermarkPath)
	}
	m.processor.SetEnableTextWatermark(config.EnableTextWatermark)
}

func (m *Manager) GetConfig() *ManagerConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return &ManagerConfig{
		MaxConcurrentJobs:   m.config.MaxConcurrentJobs,
		EnableStats:         m.config.EnableStats,
		WatermarkPath:       m.config.WatermarkPath,
		CacheEnabled:        m.config.CacheEnabled,
		MaxImageSize:        m.config.MaxImageSize,
		EnableTextWatermark: m.config.EnableTextWatermark,
	}
}

// IsConfigValid 检查配置是否有效
func (m *Manager) IsConfigValid(configJSON string) bool {
	return m.ValidateConfigJSON(configJSON) == nil
}

func (m *Manager) GetSupportedFormats() []string {
	return m.processor.GetSupportedFormats()
}

func GetDefaultConfigJSON() string {
	data, _ := json.Marshal(DefaultWatermarkConfig)
	return string(data)
}

// ParseConfigFromJSON 从JSON字符串解析配置
func ParseConfigFromJSON(configJSON string) (*WatermarkConfig, error) {
	if configJSON == "" {
		return &WatermarkConfig{Enabled: false}, nil
	}

	config := &WatermarkConfig{}
	if err := json.Unmarshal([]byte(configJSON), config); err != nil {
		return nil, fmt.Errorf("配置解析失败: %w", err)
	}

	return config, nil
}
