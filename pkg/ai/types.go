package ai

import "time"

// TokenUsage AI调用的token使用情况
type TokenUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// AIResponse 通用AI响应结构
type AIResponse struct {
	Success      bool        `json:"success"`
	Data         string      `json:"data,omitempty"`
	ErrMsg       string      `json:"errMsg,omitempty"`
	ImageURL     string      `json:"imageUrl,omitempty"`
	Usage        *TokenUsage `json:"usage,omitempty"`
	HttpDuration int64       `json:"http_duration,omitempty"` // HTTP调用耗时（毫秒）
}

// Config AI配置结构 - 复用现有配置读取方式
type Config struct {
	Enabled     bool    `json:"ai_enabled"`
	Provider    string  `json:"ai_provider"`
	APIKey      string  `json:"ai_api_key"`
	BaseURL     string  `json:"ai_proxy"`
	Model       string  `json:"ai_model"`
	MaxTokens   int     `json:"ai_max_tokens"`
	Temperature float32 `json:"ai_temperature"`
	Timeout     time.Duration
}

// FileAnalysisRequest 文件分析请求（主类型）
type FileAnalysisRequest struct {
	ImageURL  string `json:"image_url"`
	ImageData string `json:"image_data"` // base64数据
	Format    string `json:"format"`     // 文件格式
	Prompt    string `json:"prompt"`
}

// FileTaggingRequest 文件标注请求（支持标签列表）（主类型）
type FileTaggingRequest struct {
	ImageURL      string    `json:"image_url"`
	ImageData     string    `json:"image_data"` // base64数据
	Format        string    `json:"format"`     // 文件格式
	Prompt        string    `json:"prompt"`
	AvailableTags []TagInfo `json:"available_tags"` // 可用标签列表
}

// TestResult 连接测试结果
type TestResult struct {
	Success      bool                   `json:"success"`
	Message      string                 `json:"message"`
	Details      map[string]interface{} `json:"details,omitempty"`
	ResponseTime time.Duration          `json:"response_time"`
}

// ProviderInfo 提供商信息
type ProviderInfo struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Models      []string `json:"models"`
	Features    []string `json:"features"`
}

// EmbeddingRequest 文本向量化请求
type EmbeddingRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"` // 可选，使用特定模型
}

// EmbeddingResponse 文本向量化响应
type EmbeddingResponse struct {
	Success   bool        `json:"success"`
	Embedding []float32   `json:"embedding,omitempty"`
	ErrMsg    string      `json:"errMsg,omitempty"`
	Usage     *TokenUsage `json:"usage,omitempty"`
	Model     string      `json:"model,omitempty"`
	Dimension int         `json:"dimension,omitempty"`
}

// FileCategorizationRequest 文件分类请求（主类型）
type FileCategorizationRequest struct {
	ImageURL   string         `json:"image_url"`
	ImageData  string         `json:"image_data"` // base64数据
	Format     string         `json:"format"`     // 文件格式
	Categories []CategoryInfo `json:"categories"` // 可选分类列表
}

// CategoryInfo 分类信息
type CategoryInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"` // system_template/user
}

// FileCategorizationResponse 文件分类响应（主类型）
type FileCategorizationResponse struct {
	Success             bool        `json:"success"`
	CategoryID          uint        `json:"category_id,omitempty"`
	CategoryName        string      `json:"category_name,omitempty"`
	CategoryDescription string      `json:"category_description,omitempty"`
	ErrMsg              string      `json:"errMsg,omitempty"`
	Usage               *TokenUsage `json:"usage,omitempty"`
}

// TagInfo 标签信息
type TagInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"`      // category_system/user_custom/system_popular
	UsageCount  int    `json:"usage_count"` // 使用次数
}

// FileAnalysisResponse 文件分析响应（主类型）
type FileAnalysisResponse struct {
	Success     bool        `json:"success"`
	Tags        []string    `json:"tags"`
	Description string      `json:"description"`
	ErrMsg      string      `json:"errMsg,omitempty"`
	Usage       *TokenUsage `json:"usage,omitempty"`
}
