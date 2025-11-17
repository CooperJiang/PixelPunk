package models

import (
	"encoding/json"
	"fmt"
	"pixelpunk/pkg/common"
	"strconv"

	"gorm.io/gorm"
)

/* Setting 全局设置模型 */
type Setting struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Key         string `gorm:"size:100;not null;uniqueIndex:idx_setting_key" json:"key"` // 设置键名
	Value       string `gorm:"type:text;not null" json:"value"`                          // 设置值(存储为JSON字符串)
	Type        string `gorm:"size:20;not null" json:"type"`                             // 值类型: string, number, boolean, json, array
	Group       string `gorm:"size:50;not null;index" json:"group"`                      // 设置分组: system, security, mail, backup, etc. ai
	Description string `gorm:"size:500" json:"description"`                              // 设置描述
	IsSystem    bool   `gorm:"default:false" json:"is_system"`                           // 是否为系统设置(系统设置不可删除)
}

/* SettingGroup 设置分组常量 */
const (
	SettingGroupSystem   = "system"   // 系统设置
	SettingGroupSecurity = "security" // 安全设置
	SettingGroupMail     = "mail"     // 邮件设置
	SettingGroupBackup   = "backup"   // 备份设置
	SettingGroupStorage  = "storage"  // 存储设置
	SettingGroupUI       = "ui"       // 界面设置
)

/* SettingValueType 设置值类型常量 */
const (
	SettingTypeString  = "string"  // 字符串类型
	SettingTypeNumber  = "number"  // 数字类型
	SettingTypeBoolean = "boolean" // 布尔类型
	SettingTypeJSON    = "json"    // JSON对象类型
	SettingTypeArray   = "array"   // 数组类型
)

func (Setting) TableName() string {
	return "setting"
}

/* GetStringValue 获取字符串值 */
func (s *Setting) GetStringValue() string {
	if s.Type == SettingTypeString {
		if s.isHTMLContentKey() {
			return s.Value
		}

		var val string
		if err := json.Unmarshal([]byte(s.Value), &val); err != nil {
			return s.Value // 如果解析失败，返回原始值
		}
		return val
	}
	return ""
}

/* GetNumberValue 获取数字值 */
func (s *Setting) GetNumberValue() float64 {
	if s.Type == SettingTypeNumber {
		var val float64
		json.Unmarshal([]byte(s.Value), &val)
		return val
	}
	return 0
}

/* GetBoolValue 获取布尔值 */
func (s *Setting) GetBoolValue() bool {
	if s.Type == SettingTypeBoolean {
		var val bool
		json.Unmarshal([]byte(s.Value), &val)
		return val
	}
	return false
}

/* GetJSONValue 获取JSON对象 */
func (s *Setting) GetJSONValue() map[string]interface{} {
	if s.Type == SettingTypeJSON {
		var val map[string]interface{}
		json.Unmarshal([]byte(s.Value), &val)
		return val
	}
	return nil
}

/* GetArrayValue 获取数组值 */
func (s *Setting) GetArrayValue() []interface{} {
	if s.Type == SettingTypeArray {
		var val []interface{}
		json.Unmarshal([]byte(s.Value), &val)
		return val
	}
	return nil
}

/* BeforeCreate 创建前处理值 */
func (s *Setting) BeforeCreate(tx *gorm.DB) error {
	return s.normalizeValue()
}

/* BeforeUpdate 更新前处理值 */
func (s *Setting) BeforeUpdate(tx *gorm.DB) error {
	return s.normalizeValue()
}

func (s *Setting) isHTMLContentKey() bool {
	htmlKeys := []string{
		"announcement_content",
		"mail_template_content",
		"custom_html_header",
		"custom_html_footer",
		"privacy_policy_content",   // 隐私政策HTML内容
		"terms_of_service_content", // 服务条款HTML内容
	}

	for _, key := range htmlKeys {
		if s.Key == key {
			return true
		}
	}
	return false
}

func (s *Setting) normalizeValue() error {
	var isValidJSON bool
	var parsedValue interface{}

	err := json.Unmarshal([]byte(s.Value), &parsedValue)
	isValidJSON = (err == nil)

	switch s.Type {
	case SettingTypeString, "text": // 支持text类型
		if s.isHTMLContentKey() {
			if isValidJSON {
				if strVal, ok := parsedValue.(string); ok {
					s.Value = strVal
				}
			}
			return nil
		}

		if isValidJSON {
			if _, ok := parsedValue.(string); !ok {
				s.Value = fmt.Sprintf("\"%v\"", parsedValue)
			}
			return nil
		}

		bytes, err := json.Marshal(s.Value)
		if err != nil {
			return err
		}
		s.Value = string(bytes)
		return nil

	case SettingTypeNumber, SettingTypeBoolean, SettingTypeJSON, SettingTypeArray:
		if isValidJSON {
			switch s.Type {
			case SettingTypeNumber:
				if _, ok := parsedValue.(float64); !ok {
					if _, ok := parsedValue.(json.Number); !ok {
						if strVal, ok := parsedValue.(string); ok {
							if _, err := strconv.ParseFloat(strVal, 64); err == nil {
								break
							}
						}
						s.Value = "0"
					}
				}
			case SettingTypeBoolean:
				if _, ok := parsedValue.(bool); !ok {
					if strVal, ok := parsedValue.(string); ok {
						if _, err := strconv.ParseBool(strVal); err == nil {
							break
						}
					}
					s.Value = "false"
				}
			case SettingTypeArray:
				if _, ok := parsedValue.([]interface{}); !ok {
					s.Value = "[]"
				}
			case SettingTypeJSON:
				bytes, _ := json.Marshal(parsedValue)
				s.Value = string(bytes)
				return nil
			}

			bytes, _ := json.Marshal(parsedValue)
			s.Value = string(bytes)
			return nil
		}

		var newValue interface{}

		switch s.Type {
		case SettingTypeNumber:
			if f, err := strconv.ParseFloat(s.Value, 64); err == nil {
				newValue = f
			} else {
				newValue = 0
			}
		case SettingTypeBoolean:
			if b, err := strconv.ParseBool(s.Value); err == nil {
				newValue = b
			} else {
				newValue = false
			}
		case SettingTypeArray:
			newValue = []interface{}{}
		case SettingTypeJSON:
			newValue = map[string]interface{}{}
		default:
			newValue = s.Value
		}

		bytes, err := json.Marshal(newValue)
		if err != nil {
			return err
		}
		s.Value = string(bytes)
		return nil

	default:
		s.Type = SettingTypeString
		return s.normalizeValue() // 递归调用处理字符串类型
	}
}
