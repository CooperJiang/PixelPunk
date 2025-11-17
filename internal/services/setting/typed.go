package setting

import (
	"strconv"
	"strings"
)

/* GetBool 返回布尔配置，读取失败或不存在时返回默认值 */
func GetBool(group, key string, def bool) bool {
	m, err := GetSettingsByGroupAsMap(group)
	if err != nil || m == nil {
		return def
	}
	if v, ok := m.Settings[key]; ok {
		switch t := v.(type) {
		case bool:
			return t
		case string:
			s := strings.TrimSpace(strings.ToLower(t))
			if s == "true" || s == "1" || s == "yes" || s == "on" {
				return true
			}
			if s == "false" || s == "0" || s == "no" || s == "off" {
				return false
			}
		case float64:
			return t != 0
		}
	}
	return def
}

/* GetInt 返回整型配置，读取失败或不存在时返回默认值 */
func GetInt(group, key string, def int) int {
	m, err := GetSettingsByGroupAsMap(group)
	if err != nil || m == nil {
		return def
	}
	if v, ok := m.Settings[key]; ok {
		switch t := v.(type) {
		case float64:
			return int(t)
		case int:
			return t
		case string:
			if n, err := strconv.Atoi(strings.TrimSpace(t)); err == nil {
				return n
			}
		}
	}
	return def
}

/* GetFloat 返回浮点配置，读取失败或不存在时返回默认值 */
func GetFloat(group, key string, def float64) float64 {
	m, err := GetSettingsByGroupAsMap(group)
	if err != nil || m == nil {
		return def
	}
	if v, ok := m.Settings[key]; ok {
		switch t := v.(type) {
		case float64:
			return t
		case string:
			if f, err := strconv.ParseFloat(strings.TrimSpace(t), 64); err == nil {
				return f
			}
		case int:
			return float64(t)
		}
	}
	return def
}

/* GetString 返回字符串配置，读取失败或不存在时返回默认值 */
func GetString(group, key, def string) string {
	m, err := GetSettingsByGroupAsMap(group)
	if err != nil || m == nil {
		return def
	}
	if v, ok := m.Settings[key]; ok {
		switch t := v.(type) {
		case string:
			s := strings.TrimSpace(t)
			if s != "" {
				return s
			}
		case float64:
			return strconv.Itoa(int(t))
		case bool:
			if t {
				return "true"
			}
			return "false"
		}
	}
	return def
}
