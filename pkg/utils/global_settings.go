package utils

import (
	"sync"
)

// 全局设置缓存结构体，统一管理各种系统设置
type globalSettingsCache struct {
	baseUrl           string
	hideRemoteUrl     bool
	aiAnalysisEnabled bool   // AI分析功能是否启用
	adminEmail        string // 管理员邮箱
	mutex             sync.RWMutex
}

// 全局单例
var globalSettings = &globalSettingsCache{
	hideRemoteUrl:     true,
	aiAnalysisEnabled: true,
	adminEmail:        "", // 默认管理员邮箱
}

func SetBaseUrl(baseUrl string) {
	globalSettings.mutex.Lock()
	defer globalSettings.mutex.Unlock()

	globalSettings.baseUrl = baseUrl
}

func GetBaseUrl() string {
	globalSettings.mutex.RLock()
	defer globalSettings.mutex.RUnlock()

	return globalSettings.baseUrl
}

func SetHideRemoteUrl(hide bool) {
	globalSettings.mutex.Lock()
	defer globalSettings.mutex.Unlock()

	globalSettings.hideRemoteUrl = hide
}

func GetHideRemoteUrl() bool {
	globalSettings.mutex.RLock()
	defer globalSettings.mutex.RUnlock()

	return globalSettings.hideRemoteUrl
}

func SetAiAnalysisEnabled(enabled bool) {
	globalSettings.mutex.Lock()
	defer globalSettings.mutex.Unlock()

	globalSettings.aiAnalysisEnabled = enabled
}

func GetAiAnalysisEnabled() bool {
	globalSettings.mutex.RLock()
	defer globalSettings.mutex.RUnlock()

	return globalSettings.aiAnalysisEnabled
}

func SetAdminEmail(email string) {
	globalSettings.mutex.Lock()
	defer globalSettings.mutex.Unlock()

	globalSettings.adminEmail = email
}

func GetAdminEmail() string {
	globalSettings.mutex.RLock()
	defer globalSettings.mutex.RUnlock()

	return globalSettings.adminEmail
}
