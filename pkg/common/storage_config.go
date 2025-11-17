package common

import (
	"fmt"
)

// StorageSettings 存储设置结构体，避免循环导入
type StorageSettings struct {
	GuestAllowedStorageDurations []string
	GuestDefaultStorageDuration  string
	EnableGuestUpload            bool
	GuestDefaultAccessLevel      string

	UserAllowedStorageDurations []string
	UserDefaultStorageDuration  string
}

// StorageConfig 存储配置管理器
type StorageConfig struct {
	settings *StorageSettings
}

func NewStorageConfig(settings *StorageSettings) *StorageConfig {
	return &StorageConfig{
		settings: settings,
	}
}

func (sc *StorageConfig) GetAllowedDurations(isGuest bool) []string {
	var allowedDurations []string

	if isGuest {
		// 游客模式：使用guest配置（强制不超过7天）
		if len(sc.settings.GuestAllowedStorageDurations) > 0 {
			allowedDurations = sc.settings.GuestAllowedStorageDurations
		} else {
			// 备用默认选项（最多7天）
			allowedDurations = []string{"3d", "7d"}
		}
		// 过滤 >7天或永久/非法的选项
		filtered := make([]string, 0, len(allowedDurations))
		for _, d := range allowedDurations {
			if d == "permanent" {
				continue
			}
			dur := ParseStorageDuration(d)
			if dur > 0 && dur <= 7*24*60*60*1e9 { // 7天
				filtered = append(filtered, d)
			}
		}
		allowedDurations = filtered
	} else {
		// 已登录用户：使用upload配置
		if len(sc.settings.UserAllowedStorageDurations) > 0 {
			allowedDurations = sc.settings.UserAllowedStorageDurations
		} else {
			allowedDurations = []string{"3d", "7d", "30d"}
		}
		// 已登录用户总是可以选择永久存储（如果不存在的话）
		hasPermanent := false
		for _, duration := range allowedDurations {
			if duration == "permanent" {
				hasPermanent = true
				break
			}
		}
		if !hasPermanent {
			allowedDurations = append([]string{"permanent"}, allowedDurations...)
		}
	}

	return allowedDurations
}

func (sc *StorageConfig) GetDefaultDuration(isGuest bool) string {
	if isGuest {
		// 游客模式：使用guest配置
		if sc.settings.GuestDefaultStorageDuration != "" {
			return sc.settings.GuestDefaultStorageDuration
		}
		// 备用默认选项：游客使用最短时长
		allowedDurations := sc.GetAllowedDurations(true)
		if len(allowedDurations) > 0 {
			return allowedDurations[0]
		}
		return "3d"
	} else {
		// 已登录用户：使用upload配置
		if sc.settings.UserDefaultStorageDuration != "" {
			return sc.settings.UserDefaultStorageDuration
		}
		// 备用默认选项：已登录用户使用永久存储
		return "permanent"
	}
}

// ValidateStorageDuration 验证存储时长是否在允许范围内
func (sc *StorageConfig) ValidateStorageDuration(duration string, isGuest bool) error {
	allowedDurations := sc.GetAllowedDurations(isGuest)

	for _, allowed := range allowedDurations {
		if duration == allowed {
			return nil
		}
	}

	var allowedStr string
	for i, allowed := range allowedDurations {
		if i > 0 {
			allowedStr += "、"
		}
		allowedStr += allowed
	}

	if isGuest {
		return fmt.Errorf("存储时长必须是 %s（游客仅支持限时存储）", allowedStr)
	} else {
		return fmt.Errorf("存储时长必须是 %s", allowedStr)
	}
}

// IsGuestUploadEnabled 检查游客上传是否启用
func (sc *StorageConfig) IsGuestUploadEnabled() bool {
	return sc.settings.EnableGuestUpload
}

func (sc *StorageConfig) GetGuestDefaultAccessLevel() string {
	if sc.settings.GuestDefaultAccessLevel != "" {
		return sc.settings.GuestDefaultAccessLevel
	}
	return "private" // 默认私密
}

func (sc *StorageConfig) GetGuestAccessLevelText() string {
	level := sc.GetGuestDefaultAccessLevel()
	switch level {
	case "private":
		return "私密"
	case "protected":
		return "受保护"
	case "public":
		return "公开"
	default:
		return "私密"
	}
}

func (sc *StorageConfig) GetStorageDurationOptions(isGuest bool) []StorageOption {
	allowedDurations := sc.GetAllowedDurations(isGuest)
	options := make([]StorageOption, 0, len(allowedDurations))

	durationLabels := map[string]string{
		"permanent": "🔒 永久",
		"1h":        "⏰ 1小时",
		"3d":        "⏰ 3天",
		"7d":        "⏰ 7天",
		"30d":       "⏰ 30天",
	}

	for _, duration := range allowedDurations {
		label := durationLabels[duration]
		if label == "" {
			label = fmt.Sprintf("⏰ %s", duration)
		}

		days := -1 // 永久存储
		if duration != "permanent" {
			parsedDuration := ParseStorageDuration(duration)
			if parsedDuration > 0 {
				days = int(parsedDuration.Hours() / 24)
			}
		}

		options = append(options, StorageOption{
			Value: duration,
			Label: label,
			Days:  days,
		})
	}

	return options
}

// CreateDefaultStorageConfig 创建默认存储配置
func CreateDefaultStorageConfig() *StorageConfig {
	settings := &StorageSettings{
		GuestAllowedStorageDurations: []string{"3d", "7d"},
		GuestDefaultStorageDuration:  "3d",
		EnableGuestUpload:            true,
		GuestDefaultAccessLevel:      "private",
		UserAllowedStorageDurations:  []string{"3d", "7d", "30d"},
		UserDefaultStorageDuration:   "permanent",
	}
	return NewStorageConfig(settings)
}

func (sc *StorageConfig) GetValidationTag(isGuest bool) string {
	allowedDurations := sc.GetAllowedDurations(isGuest)
	if len(allowedDurations) == 0 {
		return ""
	}

	validationTag := "omitempty,oneof="
	for i, duration := range allowedDurations {
		if i > 0 {
			validationTag += " "
		}
		validationTag += duration
	}

	return validationTag
}

func (sc *StorageConfig) GetValidationMessage(isGuest bool) string {
	allowedDurations := sc.GetAllowedDurations(isGuest)
	if len(allowedDurations) == 0 {
		return "无有效的存储时长选项"
	}

	var allowedStr string
	for i, duration := range allowedDurations {
		if i > 0 {
			allowedStr += "、"
		}
		allowedStr += duration
	}

	if isGuest {
		return fmt.Sprintf("存储时长必须是 %s（游客仅支持限时存储）", allowedStr)
	} else {
		return fmt.Sprintf("存储时长必须是 %s", allowedStr)
	}
}
