package dto

// TestAIConfigDTO AI配置测试请求DTO
// 用于测试AI服务配置是否正确可用
type TestAIConfigDTO struct {
	APIKey      string  `json:"ai_api_key" binding:"required" validate:"required"`
	APIProxy    string  `json:"ai_proxy"`
	Model       string  `json:"ai_model" binding:"required" validate:"required"`
	Temperature float64 `json:"ai_temperature"`
	MaxTokens   int     `json:"ai_max_tokens"`
}

func (d *TestAIConfigDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"APIKey.required": "API密钥不能为空",
		"Model.required":  "AI模型不能为空",
	}
}
