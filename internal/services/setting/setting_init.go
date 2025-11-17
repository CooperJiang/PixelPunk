package setting

import (
	"encoding/json"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/hooks"
	"pixelpunk/pkg/utils"
)

/* RegisterSettingChangeHandler 注册设置变更处理器 */
func RegisterSettingChangeHandler(group, key string, handler func(value string)) {
	if settingChangeHandlers == nil {
		settingChangeHandlers = make(map[string]func(value string))
	}
	handlerKey := group + ":" + key
	settingChangeHandlers[handlerKey] = handler
}

func notifySettingChanged(group, key, value string) {
	if settingChangeHandlers == nil {
		return
	}

	handlerKey := group + ":" + key
	if handler, exists := settingChangeHandlers[handlerKey]; exists {
		handler(value)
	}
}

/* InitSettingService 初始化设置服务 */
func InitSettingService() {
	settingService = &SettingService{}

	// 只在首次初始化时创建 map，避免清空已注册的钩子
	if settingChangeHandlers == nil {
		settingChangeHandlers = make(map[string]func(value string))
	}

	if database.GetDB() != nil {
		syncGlobalSettings()
		preloadSettingsToCache()
	}

	hooks.RegisterSettingUpdateHook("website", func(group string) error {
		syncGlobalSettings()
		return nil
	})

	hooks.RegisterSettingUpdateHook("security", func(group string) error {
		syncGlobalSettings()
		return nil
	})

	hooks.RegisterSettingUpdateHook("upload", func(group string) error {
		syncGlobalSettings()
		return nil
	})

	hooks.RegisterSettingUpdateHook("vector", func(group string) error {
		syncGlobalSettings()
		return nil
	})

	hooks.RegisterSettingUpdateHook("ai", func(group string) error {
		syncGlobalSettings()
		return nil
	})
}

func preloadSettingsToCache() {
	commonGroups := []string{
		"website",
		"security",
		"upload",
		"registration",
		"email",
		"vector", // 添加向量设置到预加载列表
		"ai",     // 添加AI设置到预加载列表，确保ai_analysis_enabled等配置正确加载
	}

	go func() {
		for _, group := range commonGroups {
			_, _ = GetSettingsByGroupAsMap(group)
		}
	}()
}

/* GetSettingService 获取设置服务实例 */
func GetSettingService() *SettingService {
	return settingService
}

func syncGlobalSettings() {
	db := database.GetDB()
	if db == nil {
		return
	}

	var baseUrlValue string
	result := db.Table("setting").
		Select("value").
		Where("`key` = ? AND `group` = ?", "site_base_url", "website").
		Limit(1).
		Scan(&baseUrlValue)

	if result.Error == nil && result.RowsAffected > 0 {
		var baseUrl string
		if err := json.Unmarshal([]byte(baseUrlValue), &baseUrl); err == nil {
			utils.SetBaseUrl(baseUrl)
		}
	} else {
		utils.SetBaseUrl("")
	}

	var hideRemoteUrlValue string
	result = db.Table("setting").
		Select("value").
		Where("`key` = ? AND `group` = ?", "hide_remote_url", "security").
		Limit(1).
		Scan(&hideRemoteUrlValue)

	if result.Error == nil && result.RowsAffected > 0 {
		var hideRemoteUrl bool
		if err := json.Unmarshal([]byte(hideRemoteUrlValue), &hideRemoteUrl); err == nil {
			utils.SetHideRemoteUrl(hideRemoteUrl)
		} else {
			utils.SetHideRemoteUrl(true)
		}
	} else {
		utils.SetHideRemoteUrl(true)
	}

	var aiAnalysisEnabledValue string
	result = db.Table("setting").
		Select("value").
		Where("`key` = ? AND `group` = ?", "ai_analysis_enabled", "upload").
		Limit(1).
		Scan(&aiAnalysisEnabledValue)

	if result.Error == nil && result.RowsAffected > 0 {
		var aiAnalysisEnabled bool
		if err := json.Unmarshal([]byte(aiAnalysisEnabledValue), &aiAnalysisEnabled); err == nil {
			utils.SetAiAnalysisEnabled(aiAnalysisEnabled)
		} else {
			utils.SetAiAnalysisEnabled(true)
		}
	} else {
		utils.SetAiAnalysisEnabled(true)
	}

	var adminEmailValue string
	result = db.Table("setting").
		Select("value").
		Where("`key` = ? AND `group` = ?", "admin_email", "website").
		Limit(1).
		Scan(&adminEmailValue)

	if result.Error == nil && result.RowsAffected > 0 {
		var adminEmail string
		if err := json.Unmarshal([]byte(adminEmailValue), &adminEmail); err == nil {
			utils.SetAdminEmail(adminEmail)
		} else {
			utils.SetAdminEmail("")
		}
	} else {
		utils.SetAdminEmail("")
	}
}
