package ai

import (
	"context"
	"fmt"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"sync"
	"time"
)

// AIClient 统一AI服务接口
type AIClient interface {
	AnalyzeFile(ctx context.Context, req *FileAnalysisRequest) (*AIResponse, error)

	CategorizeFile(ctx context.Context, req *FileCategorizationRequest) (*FileCategorizationResponse, error)

	// 文件标注（支持标签列表）
	TagFile(ctx context.Context, req *FileTaggingRequest) (*FileAnalysisResponse, error)

	GenerateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)

	TestConnection(ctx context.Context) (*TestResult, error)

	GetProviderInfo() *ProviderInfo
}

// AIProvider AI服务提供商接口
type AIProvider interface {
	AnalyzeFile(ctx context.Context, req *FileAnalysisRequest) (*AIResponse, error)

	CategorizeFile(ctx context.Context, req *FileCategorizationRequest) (*FileCategorizationResponse, error)

	// 文件标注（支持标签列表）
	TagFile(ctx context.Context, req *FileTaggingRequest) (*FileAnalysisResponse, error)

	GenerateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error)

	TestConnection(ctx context.Context) (*TestResult, error)

	GetProviderInfo() *ProviderInfo
}

// UnifiedAIClient 统一AI客户端实现
type UnifiedAIClient struct {
	provider AIProvider
	config   *Config
}

func NewAIClient() (AIClient, error) {
	config, err := GetAIConfig()
	if err != nil {
		return nil, fmt.Errorf("获取AI配置失败: %v", err)
	}

	provider, err := createProvider(config)
	if err != nil {
		return nil, fmt.Errorf("创建AI提供商失败: %v", err)
	}

	return &UnifiedAIClient{
		provider: provider,
		config:   config,
	}, nil
}

// NewAIClientWithConfig 使用指定配置创建AI客户端
func NewAIClientWithConfig(config *Config) (AIClient, error) {
	provider, err := createProvider(config)
	if err != nil {
		return nil, fmt.Errorf("创建AI提供商失败: %v", err)
	}

	return &UnifiedAIClient{
		provider: provider,
		config:   config,
	}, nil
}

// createProvider 根据配置创建对应的AI提供商
func createProvider(config *Config) (AIProvider, error) {
	switch config.Provider {
	case "openai":
		return NewOpenAIProvider(config), nil
	default:
		return nil, fmt.Errorf("不支持的AI提供商: %s", config.Provider)
	}
}

// AnalyzeImage 分析文件
func (c *UnifiedAIClient) AnalyzeFile(ctx context.Context, req *FileAnalysisRequest) (*AIResponse, error) {
	if !c.config.Enabled {
		return &AIResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if c.config.APIKey == "" {
		return &AIResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	result, err := c.provider.AnalyzeFile(ctx, req)
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

// CategorizeImage 文件分类
func (c *UnifiedAIClient) CategorizeFile(ctx context.Context, req *FileCategorizationRequest) (*FileCategorizationResponse, error) {
	if !c.config.Enabled {
		return &FileCategorizationResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if c.config.APIKey == "" {
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

	result, err := c.provider.CategorizeFile(ctx, req)
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

// GenerateEmbedding 生成文本向量
func (c *UnifiedAIClient) GenerateEmbedding(ctx context.Context, req *EmbeddingRequest) (*EmbeddingResponse, error) {
	if !c.config.Enabled {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if c.config.APIKey == "" {
		return &EmbeddingResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	result, err := c.provider.GenerateEmbedding(ctx, req)
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

// TagImage 文件标注（支持标签列表）
func (c *UnifiedAIClient) TagFile(ctx context.Context, req *FileTaggingRequest) (*FileAnalysisResponse, error) {
	if !c.config.Enabled {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  "AI服务未启用",
		}, nil
	}

	if c.config.APIKey == "" {
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  "AI API密钥未配置",
		}, nil
	}

	// 调用具体提供商的标注方法
	result, err := c.provider.TagFile(ctx, req)
	if err != nil {
		logger.Error("AI文件标注失败: %v", err)
		return &FileAnalysisResponse{
			Success: false,
			ErrMsg:  fmt.Sprintf("AI标注失败: %v", err),
		}, err
	}

	return result, nil
}

// TestConnection 测试连接
func (c *UnifiedAIClient) TestConnection(ctx context.Context) (*TestResult, error) {
	if c.config.APIKey == "" {
		return &TestResult{
			Success: false,
			Message: "AI API密钥未配置",
		}, nil
	}

	result, err := c.provider.TestConnection(ctx)
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

func (c *UnifiedAIClient) GetProviderInfo() *ProviderInfo {
	return c.provider.GetProviderInfo()
}

// 全局AI客户端实例（使用动态配置模式）
var (
	globalClient AIClient
	clientOnce   sync.Once
)

func GetDefaultClient() AIClient {
	clientOnce.Do(func() {
		// 使用动态客户端，每次调用时自动读取最新配置
		globalClient = NewDynamicAIClient()
	})
	return globalClient
}

// RefreshDefaultClient 刷新默认AI客户端（兼容性保留，动态客户端无需刷新）
func RefreshDefaultClient() error {
	// 动态客户端每次调用时自动读取最新配置，无需手动刷新
	return nil
}

// 兼容性函数 - 保持与现有代码的兼容性

// AnalyzeImageByURL 通过URL分析文件 - 兼容现有 callOpenAI 函数
func AnalyzeImageByURL(imageURL, prompt string) (*AIResponse, error) {
	client := GetDefaultClient()

	req := &FileAnalysisRequest{
		ImageURL: imageURL,
		Prompt:   prompt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return client.AnalyzeFile(ctx, req)
}

// AnalyzeFileByURL 通过URL分析文件（新命名，等价于 AnalyzeImageByURL）
func AnalyzeFileByURL(fileURL, prompt string) (*AIResponse, error) {
	return AnalyzeImageByURL(fileURL, prompt)
}

// AnalyzeImageByBase64 通过base64数据分析文件 - 兼容现有 callOpenAIWithBase64 函数
func AnalyzeImageByBase64(base64Data, imageFormat, prompt string) (*AIResponse, error) {
	client := GetDefaultClient()

	req := &FileAnalysisRequest{
		ImageData: base64Data,
		Format:    imageFormat,
		Prompt:    prompt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	return client.AnalyzeFile(ctx, req)
}

// AnalyzeFileByBase64 通过base64分析文件（新命名，等价于 AnalyzeImageByBase64）
func AnalyzeFileByBase64(base64Data, fileFormat, prompt string) (*AIResponse, error) {
	client := GetDefaultClient()

	req := &FileAnalysisRequest{
		ImageData: base64Data,
		Format:    fileFormat,
		Prompt:    prompt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()
	return client.AnalyzeFile(ctx, req)
}

// TestAIConfiguration 测试AI配置 - 兼容现有函数
func TestAIConfiguration() (map[string]interface{}, error) {
	client := GetDefaultClient()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := client.TestConnection(ctx)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"success": result.Success,
		"message": result.Message,
	}

	if result.Details != nil {
		for k, v := range result.Details {
			response[k] = v
		}
	}

	if !result.Success {
		response["error"] = "connection_failed"
	}

	return response, nil
}

// TestAIConfigurationWithParams 使用指定参数测试AI配置 - 兼容现有函数
func TestAIConfigurationWithParams(params map[string]interface{}) (map[string]interface{}, error) {
	config := &Config{
		Enabled:     true,
		Provider:    "openai",
		APIKey:      getStringFromParams(params, "ai_api_key", ""),
		BaseURL:     utils.NormalizeOpenAIBaseURL(getStringFromParams(params, "ai_proxy", "https://api.openai.com/v1")),
		Model:       getStringFromParams(params, "ai_model", "gpt-4o"),
		MaxTokens:   4000,
		Temperature: 0.1,
	}

	if config.APIKey == "" {
		return map[string]interface{}{
			"success": false,
			"message": "AI API密钥未配置",
			"error":   "missing_api_key",
		}, nil
	}

	client, err := NewAIClientWithConfig(config)
	if err != nil {
		return map[string]interface{}{
			"success": false,
			"message": "创建AI客户端失败",
			"error":   "client_creation_failed",
		}, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := client.TestConnection(ctx)
	if err != nil {
		return nil, err
	}

	response := map[string]interface{}{
		"success": result.Success,
		"message": result.Message,
	}

	if result.Details != nil {
		for k, v := range result.Details {
			response[k] = v
		}
	}

	if !result.Success {
		response["error"] = "connection_failed"
	}

	return response, nil
}

// getStringFromParams 从参数map中获取字符串值
func getStringFromParams(params map[string]interface{}, key, defaultValue string) string {
	if value, ok := params[key]; ok {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

// 兼容性函数 - 便于现有代码调用向量化功能

// GenerateEmbeddingByText 通过文本生成向量 - 兼容现有向量化调用
func GenerateEmbeddingByText(text string) ([]float32, error) {
	client := GetDefaultClient()

	req := &EmbeddingRequest{
		Text:  text,
		Model: "text-embedding-3-small", // 使用默认模型
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	response, err := client.GenerateEmbedding(ctx, req)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf("向量化失败: %s", response.ErrMsg)
	}

	return response.Embedding, nil
}

// GenerateEmbeddingWithModel 使用指定模型生成向量
func GenerateEmbeddingWithModel(text, model string) (*EmbeddingResponse, error) {
	client := GetDefaultClient()

	req := &EmbeddingRequest{
		Text:  text,
		Model: model,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return client.GenerateEmbedding(ctx, req)
}

// 兼容性函数 - 便于现有代码调用分类功能

// CategorizeImageByURL 通过URL进行文件分类
func CategorizeImageByURL(imageURL string, categories []CategoryInfo) (*FileCategorizationResponse, error) {
	client := GetDefaultClient()

	req := &FileCategorizationRequest{
		ImageURL:   imageURL,
		Categories: categories,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return client.CategorizeFile(ctx, req)
}

// CategorizeFileByURL 通过URL进行文件分类（新命名，等价于 CategorizeImageByURL）
func CategorizeFileByURL(fileURL string, categories []CategoryInfo) (*FileCategorizationResponse, error) {
	return CategorizeImageByURL(fileURL, categories)
}

// CategorizeImageByBase64 通过base64数据进行文件分类
func CategorizeImageByBase64(base64Data, imageFormat string, categories []CategoryInfo) (*FileCategorizationResponse, error) {
	client := GetDefaultClient()

	req := &FileCategorizationRequest{
		ImageData:  base64Data,
		Format:     imageFormat,
		Categories: categories,
	}

	// 使用60秒超时的context,防止AI API调用卡死
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return client.CategorizeFile(ctx, req)
}

// CategorizeFileByBase64 通过base64进行文件分类（新命名，等价于 CategorizeImageByBase64）
func CategorizeFileByBase64(base64Data, fileFormat string, categories []CategoryInfo) (*FileCategorizationResponse, error) {
	return CategorizeImageByBase64(base64Data, fileFormat, categories)
}

// TagImageWithBase64AndTags 通过base64数据和标签列表进行文件标注
func TagImageWithBase64AndTags(base64Data, imageFormat, prompt string, availableTags []TagInfo) (*FileAnalysisResponse, error) {
	client := GetDefaultClient()

	req := &FileTaggingRequest{
		ImageData:     base64Data,
		Format:        imageFormat,
		Prompt:        prompt,
		AvailableTags: availableTags,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	return client.TagFile(ctx, req)
}

// TagFileWithBase64AndTags 通过base64和标签列表标注文件（新命名，等价于 TagImageWithBase64AndTags）
func TagFileWithBase64AndTags(base64Data, fileFormat, prompt string, availableTags []TagInfo) (*FileAnalysisResponse, error) {
	return TagImageWithBase64AndTags(base64Data, fileFormat, prompt, availableTags)
}
