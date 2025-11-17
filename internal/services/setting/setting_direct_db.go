package setting

import (
	"encoding/json"
	"fmt"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
)

// GetSettingDirectFromDB 直接从数据库获取单个配置值（绕过缓存）
// 用于AI/向量等需要实时配置的场景
func GetSettingDirectFromDB(group, key string, defaultValue interface{}) interface{} {
	db := database.GetDB()
	if db == nil {
		logger.Warn("数据库连接不可用，使用默认值: group=%s, key=%s", group, key)
		return defaultValue
	}

	type SettingRow struct {
		Value string
		Type  string
	}

	var row SettingRow
	err := db.Table("setting").
		Where("`group` = ? AND `key` = ?", group, key).
		Select("value, type").
		First(&row).Error

	if err != nil {
		return defaultValue
	}

	// 根据type解析value
	switch row.Type {
	case "boolean":
		var boolVal bool
		if err := json.Unmarshal([]byte(row.Value), &boolVal); err == nil {
			return boolVal
		}
	case "number":
		var numVal float64
		if err := json.Unmarshal([]byte(row.Value), &numVal); err == nil {
			return numVal
		}
	case "string", "text":
		// 字符串也需要JSON解析,因为数据库存储的是带引号的JSON字符串
		var strVal string
		if err := json.Unmarshal([]byte(row.Value), &strVal); err == nil {
			return strVal
		}
		return row.Value
	}

	return defaultValue
}

// GetBoolDirectFromDB 直接从数据库获取布尔值配置
func GetBoolDirectFromDB(group, key string, defaultValue bool) bool {
	val := GetSettingDirectFromDB(group, key, defaultValue)
	if boolVal, ok := val.(bool); ok {
		return boolVal
	}
	return defaultValue
}

// GetIntDirectFromDB 直接从数据库获取整数配置
func GetIntDirectFromDB(group, key string, defaultValue int) int {
	val := GetSettingDirectFromDB(group, key, float64(defaultValue))
	if floatVal, ok := val.(float64); ok {
		return int(floatVal)
	}
	return defaultValue
}

// GetFloatDirectFromDB 直接从数据库获取浮点数配置
func GetFloatDirectFromDB(group, key string, defaultValue float64) float64 {
	val := GetSettingDirectFromDB(group, key, defaultValue)
	if floatVal, ok := val.(float64); ok {
		return floatVal
	}
	return defaultValue
}

// GetStringDirectFromDB 直接从数据库获取字符串配置
func GetStringDirectFromDB(group, key string, defaultValue string) string {
	val := GetSettingDirectFromDB(group, key, defaultValue)
	if strVal, ok := val.(string); ok {
		return strVal
	}
	return defaultValue
}

// GetMultipleSettingsDirectFromDB 直接从数据库批量获取配置（绕过缓存）
// 返回map[key]interface{}
func GetMultipleSettingsDirectFromDB(group string, keys []string) (map[string]interface{}, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接不可用")
	}

	type SettingRow struct {
		Key   string
		Value string
		Type  string
	}

	var rows []SettingRow
	err := db.Table("setting").
		Where("`group` = ? AND `key` IN (?)", group, keys).
		Select("key, value, type").
		Find(&rows).Error

	if err != nil {
		return nil, fmt.Errorf("从数据库读取配置失败: %v", err)
	}

	result := make(map[string]interface{})
	for _, row := range rows {
		var parsedValue interface{}

		switch row.Type {
		case "boolean":
			var boolVal bool
			if err := json.Unmarshal([]byte(row.Value), &boolVal); err == nil {
				parsedValue = boolVal
			}
		case "number":
			var numVal float64
			if err := json.Unmarshal([]byte(row.Value), &numVal); err == nil {
				parsedValue = numVal
			}
		case "string", "text":
			// 字符串也需要JSON解析,因为数据库存储的是带引号的JSON字符串
			var strVal string
			if err := json.Unmarshal([]byte(row.Value), &strVal); err == nil {
				parsedValue = strVal
			} else {
				// 如果JSON解析失败,则直接返回原始值(向后兼容)
				parsedValue = row.Value
			}
		default:
			parsedValue = row.Value
		}

		result[row.Key] = parsedValue
	}

	return result, nil
}
