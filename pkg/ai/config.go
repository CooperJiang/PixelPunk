package ai

import (
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"time"
)

func GetAIConfig() (*Config, error) {
	// 使用现有的配置读取逻辑
	aiSettings, err := setting.GetSettingsByGroupAsMap("ai")
	if err != nil {
		logger.Error("获取AI设置失败: %v", err)
		return nil, err
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

	// 解析设置值 - 沿用现有字段名
	if val, ok := aiSettings.Settings["ai_enabled"]; ok {
		if enabled, ok := val.(bool); ok {
			config.Enabled = enabled
		}
	}

	if val, ok := aiSettings.Settings["ai_provider"]; ok {
		if provider, ok := val.(string); ok && provider != "" {
			config.Provider = provider
		}
	}

	if val, ok := aiSettings.Settings["ai_api_key"]; ok {
		if apiKey, ok := val.(string); ok {
			config.APIKey = apiKey
		}
	}

	if val, ok := aiSettings.Settings["ai_proxy"]; ok {
		if proxy, ok := val.(string); ok && proxy != "" {
			config.BaseURL = utils.NormalizeOpenAIBaseURL(proxy)
		}
	}

	if val, ok := aiSettings.Settings["ai_model"]; ok {
		if model, ok := val.(string); ok && model != "" {
			config.Model = model
		} else {
		}
	} else {
	}

	if val, ok := aiSettings.Settings["ai_max_tokens"]; ok {
		if maxTokens, ok := val.(float64); ok {
			config.MaxTokens = int(maxTokens)
		}
	}

	if val, ok := aiSettings.Settings["ai_temperature"]; ok {
		if temp, ok := val.(float64); ok {
			config.Temperature = float32(temp)
		}
	}

	if val, ok := aiSettings.Settings["ai_timeout"]; ok {
		if timeout, ok := val.(float64); ok {
			config.Timeout = time.Duration(timeout) * time.Second
		}
	}

	return config, nil
}
