package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"time"
)

// DynamicAIClient 动态AI客户端（无需初始化，每次调用时读取最新配置）
type DynamicAIClient struct{}

func NewDynamicAIClient() AIClient {
	return &DynamicAIClient{}
}

// getConfigFromDB 从数据库动态读取配置（绕过缓存，直接查数据库）
func (c *DynamicAIClient) getConfigFromDB() (*Config, error) {
	// 直接从数据库读取配置（不使用缓存）
	// 原因：缓存可能在启动时保存了空配置，导致后续配置更新无法生效
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	type SettingRow struct {
		Key   string
		Value string
		Type  string
	}

	var settings []SettingRow
	if err := db.Table("setting").
		Where("`group` = ?", "ai").
		Where("`key` IN (?)", []string{
			"ai_enabled", "ai_provider", "ai_api_key", "ai_proxy",
			"ai_model", "ai_max_tokens", "ai_temperature", "ai_timeout",
		}).
		Select("`key`, `value`, `type`").
		Find(&settings).Error; err != nil {
		return nil, fmt.Errorf("查询AI配置失败: %v", err)
	}

	config := &Config{
		Enabled:     false,
		Provider:    "openai",
		APIKey:      "",
		BaseURL:     "https://api.openai.com/v1",
		Model:       "gpt-4o",
		MaxTokens:   4000,
		Temperature: 0.1,
		Timeout:     30 * time.Second,
	}

	for _, s := range settings {
		switch s.Key {
		case "ai_enabled":
			var enabled bool
			if err := json.Unmarshal([]byte(s.Value), &enabled); err == nil {
				config.Enabled = enabled
			}
		case "ai_provider":
			// 字符串需要JSON解析
			var provider string
			if err := json.Unmarshal([]byte(s.Value), &provider); err == nil && provider != "" {
				config.Provider = provider
			}
		case "ai_api_key":
			// 字符串需要JSON解析
			var key string
			if err := json.Unmarshal([]byte(s.Value), &key); err == nil {
				config.APIKey = key
			}
		case "ai_proxy":
			// 字符串需要JSON解析
			var proxy string
			if err := json.Unmarshal([]byte(s.Value), &proxy); err == nil && proxy != "" {
				// 标准化URL：自动添加/v1，去除多余斜杠等
				config.BaseURL = utils.NormalizeOpenAIBaseURL(proxy)
			}
		case "ai_model":
			// 字符串需要JSON解析
			var model string
			if err := json.Unmarshal([]byte(s.Value), &model); err == nil && model != "" {
				config.Model = model
			}
		case "ai_max_tokens":
			var maxTokens float64
			if err := json.Unmarshal([]byte(s.Value), &maxTokens); err == nil && maxTokens > 0 {
				config.MaxTokens = int(maxTokens)
			}
		case "ai_temperature":
			var temperature float64
			if err := json.Unmarshal([]byte(s.Value), &temperature); err == nil {
				config.Temperature = float32(temperature)
			}
		case "ai_timeout":
			var timeoutSeconds float64
			if err := json.Unmarshal([]byte(s.Value), &timeoutSeconds); err == nil && timeoutSeconds > 0 {
				config.Timeout = time.Duration(timeoutSeconds) * time.Second
			}
		}
	}

	return config, nil
}

// getProvider 动态创建provider（每次调用时创建临时实例）
func (c *DynamicAIClient) getProvider() (AIProvider, *Config, error) {
	config, err := c.getConfigFromDB()
	if err != nil {
		return nil, nil, fmt.Errorf("读取AI配置失败: %v", err)
	}

	provider, err := createProvider(config)
	if err != nil {
		return nil, nil, fmt.Errorf("创建AI提供商失败: %v", err)
	}

	return provider, config, nil
}

// AnalyzeFile 分析文件（动态读取配置）
func (c *DynamicAIClient) AnalyzeFile(ctx context.Context, req *FileAnalysisRequest) (*AIResponse, error) {
	provider, config, err := c.getProvider()
	if err != nil {
		return &AIResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}, err
	}

	if !config.Enabled {
		return &AIResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if config.APIKey == "" {
		return &AIResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	result, err := provider.AnalyzeFile(ctx, req)
	if err != nil {
		logger.Error("AI文件分析失败: %v", err)
		return &AIResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI分析失败: %v", err),
		}, err
	}

	if !result.Success {
		logger.Warn("AI文件分析失败: %s", result.ErrMsg)
	}

	return result, nil
}

// CategorizeFile 文件分类（动态读取配置）
func (c *DynamicAIClient) CategorizeFile(ctx context.Context, req *FileCategorizationRequest) (*FileCategorizationResponse, error) {
	provider, config, err := c.getProvider()
	if err != nil {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}, err
	}

	if !config.Enabled {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if config.APIKey == "" {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	if len(req.Categories) == 0 {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "没有可用的分类选项",
		}, nil
	}

	result, err := provider.CategorizeFile(ctx, req)
	if err != nil {
		logger.Error("AI文件分类失败: %v", err)
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI分类失败: %v", err),
		}, err
	}

	if !result.Success {
		logger.Warn("AI文件分类失败: %s", result.ErrMsg)
	}

	return result, nil
}

// TagFile 文件标注（动态读取配置）
func (c *DynamicAIClient) TagFile(ctx context.Context, req *FileTaggingRequest) (*FileAnalysisResponse, error) {
	provider, config, err := c.getProvider()
	if err != nil {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}, err
	}

	if !config.Enabled {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if config.APIKey == "" {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	result, err := provider.TagFile(ctx, req)
	if err != nil {
		logger.Error("AI文件标注失败: %v", err)
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI标注失败: %v", err),
		}, err
	}

	return result, nil
}

// GenerateEmbedding 生成文本向量（动态读取配置）
func (c *DynamicAIClient) GenerateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	provider, config, err := c.getProvider()
	if err != nil {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  err.Error(),
		}, err
	}

	if !config.Enabled {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if config.APIKey == "" {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	result, err := provider.GenerateEmbedding(ctx, req)
	if err != nil {
		logger.Error("文本向量化失败: %v", err)
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("向量化失败: %v", err),
		}, err
	}

	if !result.Success {
		logger.Warn("文本向量化失败: %s", result.ErrMsg)
	}

	return result, nil
}

// TestConnection 测试连接（动态读取配置）
func (c *DynamicAIClient) TestConnection(ctx context.Context) (*TestResult, error) {
	provider, config, err := c.getProvider()
	if err != nil {
		return &TestResult{
			Success: false,
			Message: err.Error(),
		}, err
	}

	if config.APIKey == "" {
		return &TestResult{
			Success: false,
			Message: "AI API密钥未配置",
		}, nil
	}

	result, err := provider.TestConnection(ctx)
	if err != nil {
		logger.Error("AI连接测试失败: %v", err)
		return &TestResult{
			Success: false,
			Message: fmt.Sprintf("连接测试失败: %v", err),
		}, err
	}

	if !result.Success {
		logger.Warn("AI连接测试失败: %s", result.Message)
	}

	return result, nil
}

func (c *DynamicAIClient) GetProviderInfo() *ProviderInfo {
	provider, _, err := c.getProvider()
	if err != nil {
		logger.Warn("获取提供商信息失败: %v", err)
		return &ProviderInfo{
			Name:        "unknown",
			DisplayName: "Unknown Provider",
			Models:      []string{},
			Features:    []string{},
		}
	}

	return provider.GetProviderInfo()
}
