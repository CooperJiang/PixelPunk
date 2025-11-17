package utils

import (
	"strings"
)

// NormalizeOpenAIBaseURL 标准化OpenAI Base URL
// 处理用户可能输入的各种格式：
// - https://api.openai.com       -> https://api.openai.com/v1
// - https://api.openai.com/      -> https://api.openai.com/v1
// - https://api.openai.com/v1    -> https://api.openai.com/v1
// - https://api.openai.com/v1/   -> https://api.openai.com/v1
// - https://api.996444.cn        -> https://api.996444.cn/v1
// - https://api.996444.cn/       -> https://api.996444.cn/v1
// - https://api.996444.cn/v1     -> https://api.996444.cn/v1
// - https://api.996444.cn/v1/    -> https://api.996444.cn/v1
func NormalizeOpenAIBaseURL(baseURL string) string {
	if baseURL == "" {
		return "https://api.openai.com/v1"
	}

	baseURL = strings.TrimSpace(baseURL)

	baseURL = strings.TrimRight(baseURL, "/")

	// 检查是否已经包含 /v1
	if strings.HasSuffix(baseURL, "/v1") {
		return baseURL
	}

	// 如果不包含 /v1，添加它
	return baseURL + "/v1"
}

// ValidateOpenAIBaseURL 验证OpenAI Base URL是否有效
// 返回标准化后的URL和是否有效的标志
func ValidateOpenAIBaseURL(baseURL string) (string, bool) {
	if baseURL == "" {
		return "https://api.openai.com/v1", true
	}

	normalized := NormalizeOpenAIBaseURL(baseURL)

	// 基本验证：必须是https开头
	if !strings.HasPrefix(normalized, "http://") && !strings.HasPrefix(normalized, "https://") {
		return normalized, false
	}

	if len(normalized) < 10 {
		return normalized, false
	}

	return normalized, true
}
