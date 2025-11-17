package config

import (
	"fmt"
	"strconv"
)

// MapConfig 基于map的配置实现
type MapConfig struct {
	data map[string]interface{}
}

func NewMapConfig(data map[string]interface{}) *MapConfig {
	if data == nil {
		data = make(map[string]interface{})
	}
	return &MapConfig{data: data}
}

// Get 获取原始值
func (c *MapConfig) Get(key string) interface{} {
	return c.data[key]
}

func (c *MapConfig) GetString(key string) string {
	if val, ok := c.data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", val)
	}
	return ""
}

func (c *MapConfig) GetInt(key string) int {
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case int:
			return v
		case int64:
			return int(v)
		case float64:
			return int(v)
		case string:
			if i, err := strconv.Atoi(v); err == nil {
				return i
			}
		}
	}
	return 0
}

// GetInt64 获取int64值
func (c *MapConfig) GetInt64(key string) int64 {
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case int64:
			return v
		case int:
			return int64(v)
		case float64:
			return int64(v)
		case string:
			if i, err := strconv.ParseInt(v, 10, 64); err == nil {
				return i
			}
		}
	}
	return 0
}

func (c *MapConfig) GetBool(key string) bool {
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case bool:
			return v
		case string:
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		case int:
			return v != 0
		case float64:
			return v != 0
		}
	}
	return false
}

// GetFloat64 获取浮点数值
func (c *MapConfig) GetFloat64(key string) float64 {
	if val, ok := c.data[key]; ok {
		switch v := val.(type) {
		case float64:
			return v
		case float32:
			return float64(v)
		case int:
			return float64(v)
		case int64:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
	}
	return 0
}

// Set 设置值
func (c *MapConfig) Set(key string, value interface{}) {
	c.data[key] = value
}

// Has 检查键是否存在
func (c *MapConfig) Has(key string) bool {
	_, exists := c.data[key]
	return exists
}

func (c *MapConfig) GetAll() map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range c.data {
		result[k] = v
	}
	return result
}

// Clone 克隆配置
func (c *MapConfig) Clone() *MapConfig {
	return NewMapConfig(c.GetAll())
}

// Merge 合并配置
func (c *MapConfig) Merge(other *MapConfig) {
	for k, v := range other.data {
		c.data[k] = v
	}
}

func (c *MapConfig) GetStringWithDefault(key, defaultValue string) string {
	if val := c.GetString(key); val != "" {
		return val
	}
	return defaultValue
}

func (c *MapConfig) GetIntWithDefault(key string, defaultValue int) int {
	if c.Has(key) {
		return c.GetInt(key)
	}
	return defaultValue
}

func (c *MapConfig) GetBoolWithDefault(key string, defaultValue bool) bool {
	if c.Has(key) {
		return c.GetBool(key)
	}
	return defaultValue
}
