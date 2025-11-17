package setting

import (
	"testing"
)

func TestNormalizeSiteBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "带 https:// 和尾部斜杠",
			input:    "https://cdn.example.com/",
			expected: "https://cdn.example.com",
		},
		{
			name:     "带 http:// 和尾部斜杠",
			input:    "http://cdn.example.com/",
			expected: "http://cdn.example.com",
		},
		{
			name:     "只带尾部斜杠",
			input:    "cdn.example.com/",
			expected: "cdn.example.com",
		},
		{
			name:     "标准格式",
			input:    "cdn.example.com",
			expected: "cdn.example.com",
		},
		{
			name:     "带前后空格",
			input:    "  cdn.example.com  ",
			expected: "cdn.example.com",
		},
		{
			name:     "带 https 前缀空格和尾部斜杠",
			input:    " https://cdn.example.com/ ",
			expected: "https://cdn.example.com",
		},
		{
			name:     "只有 https://",
			input:    "https://",
			expected: "https:/", // TrimSuffix 会去掉末尾斜杠
		},
		{
			name:     "空字符串",
			input:    "",
			expected: "",
		},
		{
			name:     "带端口号",
			input:    "http://192.168.1.100:8080/",
			expected: "http://192.168.1.100:8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSiteBaseURL(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeSiteBaseURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeSettingValue(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "规范化 site_base_url",
			key:      "site_base_url",
			value:    "https://cdn.example.com/",
			expected: "https://cdn.example.com",
		},
		{
			name:     "非 site_base_url 键值不变",
			key:      "other_key",
			value:    "https://example.com/",
			expected: "https://example.com/",
		},
		{
			name:     "site_base_url 非字符串值不变",
			key:      "site_base_url",
			value:    123,
			expected: 123,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeSettingValue(tt.key, tt.value)
			if result != tt.expected {
				t.Errorf("normalizeSettingValue(%q, %v) = %v, want %v", tt.key, tt.value, result, tt.expected)
			}
		})
	}
}
