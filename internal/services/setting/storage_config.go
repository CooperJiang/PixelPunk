package setting

import (
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
)

/* CreateStorageConfig 创建存储配置管理器 */
func CreateStorageConfig() (*common.StorageConfig, error) {
	guestSettings, err := GetSettingsByGroupAsMap("guest")
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取游客设置失败")
	}
	uploadSettings, err := GetSettingsByGroupAsMap("upload")
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取上传设置失败")
	}

	var guestAllowedDurations []string
	if durations, ok := guestSettings.Settings["guest_allowed_storage_durations"].([]interface{}); ok {
		for _, duration := range durations {
			if durationStr, ok := duration.(string); ok {
				guestAllowedDurations = append(guestAllowedDurations, durationStr)
			}
		}
	}
	var userAllowedDurations []string
	if durations, ok := uploadSettings.Settings["user_allowed_storage_durations"].([]interface{}); ok {
		for _, duration := range durations {
			if durationStr, ok := duration.(string); ok {
				userAllowedDurations = append(userAllowedDurations, durationStr)
			}
		}
	}

	settings := &common.StorageSettings{
		GuestAllowedStorageDurations: guestAllowedDurations,
		GuestDefaultStorageDuration:  getString(guestSettings.Settings, "guest_default_storage_duration"),
		EnableGuestUpload:            getBool(guestSettings.Settings, "enable_guest_upload"),
		GuestDefaultAccessLevel:      getString(guestSettings.Settings, "guest_default_access_level"),
		UserAllowedStorageDurations:  userAllowedDurations,
		UserDefaultStorageDuration:   getString(uploadSettings.Settings, "user_default_storage_duration"),
	}
	return common.NewStorageConfig(settings), nil
}

func getString(settings map[string]interface{}, key string) string {
	if value, ok := settings[key].(string); ok {
		return value
	}
	return ""
}
func getBool(settings map[string]interface{}, key string) bool {
	if value, ok := settings[key].(bool); ok {
		return value
	}
	return false
}
