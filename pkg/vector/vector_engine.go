package vector

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"sync"

	"gorm.io/gorm"
)

// VectorConfig 向量配置结构体
type VectorConfig struct {
	Enabled             bool    `json:"enabled"`
	Provider            string  `json:"provider"`
	Model               string  `json:"model"`
	APIKey              string  `json:"api_key"`
	BaseURL             string  `json:"base_url"`
	Timeout             int     `json:"timeout"`
	SimilarityThreshold float32 `json:"similarity_threshold"`
	MaxResults          int     `json:"max_results"`
	QdrantEmbeddedMode  bool    `json:"qdrant_embedded_mode"`
	QdrantHost          string  `json:"qdrant_host"`
	QdrantPort          int     `json:"qdrant_port"`
}

// getVectorConfigFromDB 从数据库获取向量配置
func getVectorConfigFromDB() (*VectorConfig, error) {
	vectorSettings, err := setting.GetSettingsByGroupAsMap("vector")
	if err != nil {
		logger.Error("获取向量设置失败: %v", err)
		return nil, fmt.Errorf("获取向量设置失败: %v", err)
	}

	config := &VectorConfig{
		Enabled:             false,
		Provider:            "openai",
		Model:               "text-embedding-3-small",
		APIKey:              "",
		BaseURL:             "https://api.openai.com/v1",
		Timeout:             30,
		SimilarityThreshold: 0.3,
		MaxResults:          50,
		// Qdrant默认设置 - 使用独立服务
		QdrantEmbeddedMode: false,
		QdrantHost:         "localhost",
		QdrantPort:         6333,
	}

	if val, ok := vectorSettings.Settings["vector_enabled"]; ok {
		if enabled, ok := val.(bool); ok {
			config.Enabled = enabled
		}
	}

	if val, ok := vectorSettings.Settings["vector_provider"]; ok {
		if provider, ok := val.(string); ok && provider != "" {
			config.Provider = provider
		}
	}

	if val, ok := vectorSettings.Settings["vector_model"]; ok {
		if model, ok := val.(string); ok && model != "" {
			config.Model = model
		}
	}

	if val, ok := vectorSettings.Settings["vector_api_key"]; ok {
		if apiKey, ok := val.(string); ok {
			config.APIKey = apiKey
		}
	}

	if val, ok := vectorSettings.Settings["vector_base_url"]; ok {
		if baseURL, ok := val.(string); ok && baseURL != "" {
			config.BaseURL = baseURL
		}
	}

	if val, ok := vectorSettings.Settings["vector_timeout"]; ok {
		if timeout, ok := val.(float64); ok {
			config.Timeout = int(timeout)
		}
	}

	if val, ok := vectorSettings.Settings["vector_similarity_threshold"]; ok {
		if threshold, ok := val.(float64); ok {
			config.SimilarityThreshold = float32(threshold)
		}
	}

	if val, ok := vectorSettings.Settings["vector_max_results"]; ok {
		if maxResults, ok := val.(float64); ok {
			config.MaxResults = int(maxResults)
		}
	}

	if val, ok := vectorSettings.Settings["qdrant_embedded_mode"]; ok {
		if embeddedMode, ok := val.(bool); ok {
			config.QdrantEmbeddedMode = embeddedMode
		}
	}

	if val, ok := vectorSettings.Settings["qdrant_host"]; ok {
		if host, ok := val.(string); ok {
			config.QdrantHost = host
		}
	}

	if val, ok := vectorSettings.Settings["qdrant_port"]; ok {
		if port, ok := val.(float64); ok {
			config.QdrantPort = int(port)
		}
	}

	return config, nil
}

// VectorEngine 向量引擎
type VectorEngine struct {
	db        *gorm.DB
	storage   VectorStorage
	embedding EmbeddingProvider
	enabled   bool
	mutex     sync.RWMutex
}

// VectorService 向量服务接口
type VectorService interface {
	IsEnabled() bool
	Initialize() error

	ProcessFile(fileID, description string) error
	BatchProcessFiles(items []VectorItem) error

	SearchFiles(query string, limit int, userID uint, threshold float32) ([]VectorSearchResult, error)

	DeleteVector(fileID string) error
	GetUserVectorCount(userID uint) (int64, error)
	GetStorageStats() (*VectorStorageStats, error)
	VectorExists(fileID string) (bool, error) // 新增：检查向量是否存在

	HealthCheck() error
}

var (
	globalVectorEngine *VectorEngine
	engineOnce         sync.Once
)

func NewVectorEngine(db *gorm.DB) *VectorEngine {
	engine := &VectorEngine{
		db:      db,
		enabled: false, // 初始状态为禁用，懒加载时检查配置
	}

	// 存储层将在初始化时使用Qdrant向量数据库
	return engine
}

func GetGlobalVectorEngine() *VectorEngine {
	return globalVectorEngine
}

func GetEngine() *VectorEngine {
	return GetGlobalVectorEngine()
}

// InitGlobalVectorEngine 初始化全局向量引擎（嵌入式模式，已弃用）
func InitGlobalVectorEngine(db *gorm.DB) error {
	engineOnce.Do(func() {
		globalVectorEngine = NewVectorEngine(db)
	})

	if globalVectorEngine != nil && globalVectorEngine.enabled {
		return globalVectorEngine.Initialize()
	}

	return nil
}

// InitQdrantVectorEngine 初始化Qdrant向量引擎（直连模式，使用动态配置）
func InitQdrantVectorEngine(qdrantURL string, timeout int) error {
	engineOnce.Do(func() {
		qdrantClient := NewQdrantClient(qdrantURL, timeout)

		if err := qdrantClient.InitCollection(); err != nil {
			logger.Error("初始化Qdrant集合失败: %v", err)
		}

		// 使用动态OpenAI客户端（无需初始化配置，每次调用时动态读取）
		dynamicClient := NewDynamicOpenAIClient()

		db := database.GetDB()
		if db == nil {
			logger.Error("数据库连接不可用，向量引擎初始化失败")
			return
		}

		globalVectorEngine = &VectorEngine{
			db:        db,
			storage:   qdrantClient,
			embedding: dynamicClient, // 动态客户端，自动读取最新配置
			enabled:   true,
		}
	})

	return nil
}

// ensureInitialized 确保向量引擎已正确初始化（简化版，动态客户端无需懒加载）
func (ve *VectorEngine) ensureInitialized() error {
	if ve == nil {
		return fmt.Errorf("向量引擎未初始化")
	}

	ve.mutex.RLock()
	defer ve.mutex.RUnlock()

	if ve.storage == nil {
		return fmt.Errorf("向量存储未初始化")
	}

	if !ve.enabled {
		return fmt.Errorf("向量搜索功能已禁用")
	}

	if ve.embedding == nil {
		return fmt.Errorf("embedding客户端未初始化")
	}

	// 对于QdrantClient，检查连接是否可用
	if qdrantClient, ok := ve.storage.(*QdrantClient); ok {
		if err := qdrantClient.HealthCheck(); err != nil {
			logger.Warn("Qdrant连接检查失败: %v", err)
			return fmt.Errorf("qdrant不可用: %v", err)
		}
	}

	return nil
}

// Close 关闭向量引擎，清理资源
func (ve *VectorEngine) Close() error {
	if ve == nil {
		return nil
	}

	ve.mutex.Lock()
	defer ve.mutex.Unlock()

	// Qdrant HTTP客户端不需要显式关闭

	ve.storage = nil
	ve.embedding = nil
	ve.enabled = false

	return nil
}

// IsEnabled 检查向量搜索是否启用
func (ve *VectorEngine) IsEnabled() bool {
	if ve == nil {
		return false
	}

	if err := ve.ensureInitialized(); err != nil {
		return false
	}

	ve.mutex.RLock()
	defer ve.mutex.RUnlock()
	return ve.enabled
}

// Initialize 初始化向量引擎
func (ve *VectorEngine) Initialize() error {
	if !ve.IsEnabled() {
		return fmt.Errorf("向量搜索功能未启用")
	}

	if err := ve.db.AutoMigrate(&models.FileVector{}); err != nil {
		logger.Error("向量表迁移失败: %v", err)
		return fmt.Errorf("数据库迁移失败: %v", err)
	}

	if openaiClient, ok := ve.embedding.(*OpenAIEmbeddingClient); ok {
		if err := openaiClient.ValidateAPIKey(); err != nil {
			logger.Warn("OpenAI API密钥验证失败: %v", err)
			// 不要因为API密钥验证失败就禁用功能，可能是网络问题
		}
	}

	return nil
}

// CloneVectorFrom 从已有文件复制向量到新文件（用于重复文件零成本复用）
// description 为空时将尝试从Qdrant payload 或 AIInfo 获取
func (ve *VectorEngine) CloneVectorFrom(originalID, newID, description string) error {
	if ve == nil {
		return nil
	}
	if err := ve.ensureInitialized(); err != nil {
		return fmt.Errorf("向量引擎未就绪: %v", err)
	}

	// 仅在存储为 Qdrant 时支持直接拷贝
	qc, ok := ve.storage.(*QdrantClient)
	if !ok {
		return nil // 非Qdrant存储暂不处理，后续由队列兜底
	}

	vec, payload, err := qc.FetchVectorWithPayload(originalID)
	if err != nil {
		return err
	}

	desc := description
	if desc == "" {
		if v, ok := payload["description"].(string); ok {
			desc = v
		}
	}
	if desc == "" {
		// 回退从 AIInfo 表读取
		var ai models.FileAIInfo
		_ = ve.db.Where("file_id = ?", newID).First(&ai).Error
		if ai.Description != "" {
			desc = ai.Description
		}
		if desc == "" {
			desc = ai.SearchContent
		}
	}
	model := "text-embedding-3-small"
	if v, ok := payload["model"].(string); ok && v != "" {
		model = v
	}

	if err := ve.storage.StoreVector(newID, vec, desc, model); err != nil {
		return err
	}

	var existing models.FileVector
	if err := ve.db.Where("file_id = ?", newID).First(&existing).Error; err == nil {
		_ = ve.db.Model(&existing).Updates(map[string]interface{}{
			"status":        common.VectorStatusCompleted,
			"description":   desc,
			"model":         model,
			"dimension":     len(vec),
			"error_message": "",
		}).Error
	} else {
		iv := &models.FileVector{
			FileID:      newID,
			Description: desc,
			Model:       model,
			Dimension:   len(vec),
			Status:      common.VectorStatusCompleted,
		}
		_ = ve.db.Create(iv).Error
	}
	return nil
}

// ProcessFile 处理单个文件的向量化（命名保留，内部统一用 fileID）
func (ve *VectorEngine) ProcessFile(fileID, description string) error {
	if err := ve.ensureInitialized(); err != nil {
		logger.Warn("向量引擎初始化失败: %v", err)
		return fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	if description == "" {
		logger.Error("文件描述为空，无法进行向量化")
		return fmt.Errorf("文件描述为空")
	}

	vector, err := ve.embedding.GenerateEmbedding(description)
	if err != nil {
		logger.Error("OpenAI向量生成失败 [%s]: %v", fileID, err)
		return fmt.Errorf("向量化失败: %v", err)
	}

	model := ve.getCurrentModel()

	if err := ve.storage.StoreVector(fileID, vector, description, model); err != nil {
		logger.Error("数据库存储向量失败 [%s]: %v", fileID, err)
		return fmt.Errorf("存储失败: %v", err)
	}

	return nil
}

// BatchProcessFiles 批量处理文件向量化
func (ve *VectorEngine) BatchProcessFiles(items []VectorItem) error {
	if err := ve.ensureInitialized(); err != nil {
		return fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	if len(items) == 0 {
		return nil
	}

	descriptions := make([]string, len(items))
	for i, item := range items {
		descriptions[i] = item.Description
	}

	vectors, err := ve.embedding.BatchGenerateEmbeddings(descriptions)
	if err != nil {
		logger.Error("批量向量化失败: %v", err)
		return fmt.Errorf("批量向量化失败: %v", err)
	}

	model := ve.getCurrentModel()

	vectorItems := make([]VectorItem, len(items))
	for i, item := range items {
		vectorItems[i] = VectorItem{
			FileID:      item.FileID,
			Vector:      vectors[i],
			Description: item.Description,
			Model:       model,
		}
	}

	if err := ve.storage.BatchStoreVectors(vectorItems); err != nil {
		logger.Error("批量存储向量失败: %v", err)
		return fmt.Errorf("批量存储失败: %v", err)
	}

	return nil
}

// SearchSimilarByFileID 通过文件ID搜索相似文件
func (ve *VectorEngine) SearchSimilarByFileID(fileID string, limit int, userID uint, threshold float32) ([]VectorSearchResult, error) {
	if err := ve.ensureInitialized(); err != nil {
		return nil, fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	var baseVector models.FileVector
	if err := ve.db.Where("file_id = ? AND status = ?", fileID, common.VectorStatusCompleted).First(&baseVector).Error; err != nil {
		return nil, fmt.Errorf("文件向量信息不存在或未完成处理")
	}

	// 使用Qdrant进行向量相似度搜索
	if qdrantClient, ok := ve.storage.(*QdrantClient); ok {
		// 直接使用 fileID 从 Qdrant 搜索相似向量
		return qdrantClient.SearchSimilarByID(fileID, limit, userID, threshold, baseVector.Model)
	}

	// 如果不是Qdrant存储，降级使用描述文本搜索
	logger.Warn("非Qdrant存储，降级使用文本搜索")
	return ve.SearchFiles(baseVector.Description, limit, userID, threshold)
}

// SearchFiles 搜索相似文件
func (ve *VectorEngine) SearchFiles(query string, limit int, userID uint, threshold float32) ([]VectorSearchResult, error) {
	if err := ve.ensureInitialized(); err != nil {
		logger.Error("向量引擎未初始化: %v", err)
		return nil, fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	if query == "" {
		return nil, fmt.Errorf("搜索查询为空")
	}

	queryVector, err := ve.embedding.GenerateEmbedding(query)
	if err != nil {
		logger.Error("查询向量化失败: %v", err)
		return nil, fmt.Errorf("查询向量化失败: %v", err)
	}

	model := ve.getCurrentModel()

	results, err := ve.storage.SearchSimilar(queryVector, limit, userID, threshold, model)
	if err != nil {
		logger.Error("向量搜索失败: %v", err)
		return nil, fmt.Errorf("搜索失败: %v", err)
	}

	return results, nil
}

func (ve *VectorEngine) DeleteVector(fileID string) error {
	// 如果向量功能未启用，不执行删除操作（不算错误）
	if err := ve.ensureInitialized(); err != nil {
		return nil
	}

	return ve.storage.DeleteVector(fileID)
}

func (ve *VectorEngine) GetUserVectorCount(userID uint) (int64, error) {
	if err := ve.ensureInitialized(); err != nil {
		return 0, fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	return ve.storage.GetVectorCount(userID)
}

func (ve *VectorEngine) GetStorageStats() (*VectorStorageStats, error) {
	if err := ve.ensureInitialized(); err != nil {
		return nil, fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	return ve.storage.GetStorageStats()
}

// VectorExists 检查向量是否存在
func (ve *VectorEngine) VectorExists(fileID string) (bool, error) {
	if err := ve.ensureInitialized(); err != nil {
		return false, fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	return ve.storage.VectorExists(fileID)
}

func (ve *VectorEngine) GetAllFileIDs(limit int) ([]string, error) {
	if err := ve.ensureInitialized(); err != nil {
		return nil, fmt.Errorf("向量搜索功能不可用: %v", err)
	}
	// 仅在 Qdrant 存储下支持
	if qdrantClient, ok := ve.storage.(*QdrantClient); ok {
		return qdrantClient.GetAllFileIDs(limit)
	}
	return nil, fmt.Errorf("当前存储不支持遍历所有向量id")
}

// HealthCheck 健康检查
func (ve *VectorEngine) HealthCheck() error {
	if err := ve.ensureInitialized(); err != nil {
		return fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	if ve.storage == nil {
		return fmt.Errorf("向量存储未初始化")
	}

	// 对于 QdrantClient，检查连接健康状态
	if qdrantClient, ok := ve.storage.(*QdrantClient); ok {
		if err := qdrantClient.HealthCheck(); err != nil {
			return fmt.Errorf("qdrant连接不健康: %v", err)
		}
	}

	// 注意：在直连模式下，ve.db 可以为 nil，ve.embedding 也可能不需要

	return nil
}

// ReloadConfig 重新加载向量配置（动态客户端无需操作，配置会自动读取）
func (ve *VectorEngine) ReloadConfig() error {
	if ve == nil {
		return fmt.Errorf("向量引擎未初始化")
	}

	// 动态客户端会在每次调用时自动读取最新配置，无需手动重新加载
	return nil
}


// Disable 禁用向量搜索功能
func (ve *VectorEngine) Disable() {
	if ve == nil {
		return
	}

	ve.mutex.Lock()
	defer ve.mutex.Unlock()

	ve.enabled = false
}

// getCurrentModel 获取当前向量模型
func (ve *VectorEngine) getCurrentModel() string {
	vectorConfig, err := getVectorConfigFromDB()
	if err != nil {
		logger.Warn("获取向量模型配置失败，使用默认值: %v", err)
		return "text-embedding-3-small"
	}

	if vectorConfig.Model != "" {
		return vectorConfig.Model
	}

	return "text-embedding-3-small" // 默认模型
}

// ProcessFileWithModel 使用指定模型处理单个文件的向量化
func (ve *VectorEngine) ProcessFileWithModel(fileID, description, model string) error {
	if err := ve.ensureInitialized(); err != nil {
		logger.Warn("向量引擎初始化失败: %v", err)
		return fmt.Errorf("向量搜索功能不可用: %v", err)
	}

	if description == "" {
		logger.Error("文件描述为空，无法进行向量化")
		return fmt.Errorf("文件描述为空")
	}

	if model == "" {
		model = "text-embedding-3-small" // 默认模型
	}

	vector, err := ve.embedding.GenerateEmbedding(description)
	if err != nil {
		logger.Error("向量生成失败 [%s] (模型: %s): %v", fileID, model, err)
		return fmt.Errorf("向量化失败: %v", err)
	}

	if err := ve.storage.StoreVector(fileID, vector, description, model); err != nil {
		logger.Error("数据库存储向量失败 [%s] (模型: %s): %v", fileID, model, err)
		return fmt.Errorf("存储失败: %v", err)
	}

	return nil
}

// AddFileToVectorQueue 添加文件到向量处理队列
func AddFileToVectorQueue(file models.File) {
	if !IsVectorEnabled() {
		return
	}

	if file.Description == "" {
		return
	}

	// 异步处理，避免阻塞上传流程
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("AddFileToVectorQueue goroutine panic: %v, 文件ID: %s", r, file.ID)
			}
		}()

		db := database.GetDB()
		if db == nil {
			logger.Error("无法获取数据库连接，跳过向量队列添加")
			return
		}

		var existingVector models.FileVector
		err := db.Where("file_id = ?", file.ID).First(&existingVector).Error
		if err == nil {
			// 记录已存在，不需要重复添加
			return
		}

		newVector := models.FileVector{
			FileID:      file.ID,
			Status:      common.VectorStatusPending,
			Description: file.Description,
		}

		if err := db.Create(&newVector).Error; err != nil {
			logger.Error("创建向量记录失败: %v", err)
			return
		}

		// 通知队列状态变更（新服务聚合）
		// 避免循环依赖：统计推送由 vector_queue_service 负责
	}()
}

// IsVectorEnabled 检查向量功能是否启用
func IsVectorEnabled() bool {
	engine := GetGlobalVectorEngine()
	return engine != nil && engine.IsEnabled()
}
