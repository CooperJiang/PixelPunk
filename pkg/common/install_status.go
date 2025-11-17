package common

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// InstallStatus 系统安装状态
type InstallStatus struct {
	Installed    bool   `json:"installed"`
	DatabaseOK   bool   `json:"database_ok"`
	RedisOK      bool   `json:"redis_ok"`
	Message      string `json:"message"`
	DeployMode   string `json:"deploy_mode"`   // standalone/docker/compose
	SetupLevel   int    `json:"setup_level"`   // 1=基础配置 2=完整配置
	ConfigPreset bool   `json:"config_preset"` // 配置是否预设
}

// InstallManager 安装状态管理器
type InstallManager struct {
	mutex           sync.RWMutex
	systemInstalled bool
	installMode     bool
	installing      bool       // 正在安装标志
	installMutex    sync.Mutex // 安装过程互斥锁
}

var installManager = &InstallManager{
	systemInstalled: false,
	installMode:     false,
	installing:      false,
}

func GetInstallManager() *InstallManager {
	return installManager
}

func (im *InstallManager) SetInstallMode(mode bool) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.installMode = mode
}

// IsInstallMode 检查是否为安装模式
func (im *InstallManager) IsInstallMode() bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return im.installMode
}

func (im *InstallManager) SetSystemInstalled(installed bool) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.systemInstalled = installed
}

// IsSystemInstalled 检查系统是否已安装
func (im *InstallManager) IsSystemInstalled() bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return im.systemInstalled
}

func (im *InstallManager) GetStatus() InstallStatus {
	im.mutex.RLock()
	defer im.mutex.RUnlock()

	var message string
	if im.installing {
		message = "系统正在安装中，请等待安装完成"
	} else if im.installMode {
		message = "系统未安装，请先完成安装配置"
	} else if im.systemInstalled {
		message = "系统已安装并运行正常"
	} else {
		message = "系统状态未知"
	}

	deployMode := GetDeployMode()
	setupLevel := GetSetupLevel(deployMode, im.installMode)
	configPreset := IsConfigPreset()

	return InstallStatus{
		Installed:    im.systemInstalled,
		DatabaseOK:   !im.installMode && im.systemInstalled,
		RedisOK:      true,
		Message:      message,
		DeployMode:   deployMode,
		SetupLevel:   setupLevel,
		ConfigPreset: configPreset,
	}
}

// GetDeployMode 获取部署模式
func GetDeployMode() string {
	mode := strings.ToLower(os.Getenv("DEPLOY_MODE"))
	switch mode {
	case "docker":
		return "docker"
	case "compose":
		return "compose"
	default:
		return "standalone"
	}
}

// GetSetupLevel 获取Setup配置级别
func GetSetupLevel(deployMode string, installMode bool) int {
	if !installMode {
		return 0
	}

	switch deployMode {
	case "compose":
		return 1
	case "docker":
		return 1
	default:
		return 2
	}
}

// IsConfigPreset 检查配置是否预设
func IsConfigPreset() bool {
	preset := strings.ToLower(os.Getenv("CONFIG_PRESET"))
	return preset == "true" || preset == "1"
}

// StartInstall 开始安装过程（确保同时只有一个安装过程）
func (im *InstallManager) StartInstall() error {
	im.installMutex.Lock()
	defer im.installMutex.Unlock()

	im.mutex.Lock()
	defer im.mutex.Unlock()

	if im.installing {
		return fmt.Errorf("系统正在安装中，请稍后重试")
	}

	if !im.installMode {
		return fmt.Errorf("系统已安装，无需重复安装")
	}

	im.installing = true
	return nil
}

// FinishInstall 完成安装过程
func (im *InstallManager) FinishInstall(success bool) {
	im.mutex.Lock()
	defer im.mutex.Unlock()

	im.installing = false
	if success {
		im.installMode = false
		im.systemInstalled = true
	}
}

// IsInstalling 检查是否正在安装
func (im *InstallManager) IsInstalling() bool {
	im.mutex.RLock()
	defer im.mutex.RUnlock()
	return im.installing
}

// ResetToInstallMode 重置为安装模式
func (im *InstallManager) ResetToInstallMode(reason string) {
	im.mutex.Lock()
	defer im.mutex.Unlock()
	im.installMode = true
	im.systemInstalled = false
	im.installing = false
}
