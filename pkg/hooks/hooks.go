package hooks

import (
	"sync"
)

// SettingUpdateHook 定义设置更新时的钩子函数类型
type SettingUpdateHook func(group string) error

var (
	settingUpdateHooks = make(map[string][]SettingUpdateHook)
	mutex              sync.RWMutex
)

// RegisterSettingUpdateHook 注册设置更新钩子
// group: 设置组名，为空表示所有组
// hook: 当指定组的设置更新时调用的函数
func RegisterSettingUpdateHook(group string, hook SettingUpdateHook) {
	mutex.Lock()
	defer mutex.Unlock()

	settingUpdateHooks[group] = append(settingUpdateHooks[group], hook)
}

// TriggerSettingUpdate 触发设置更新钩子
// group: 已更新的设置组名
func TriggerSettingUpdate(group string) {
	mutex.RLock()
	defer mutex.RUnlock()

	for _, hook := range settingUpdateHooks[group] {
		err := hook(group)
		if err != nil {
		}
	}

	// 调用通用钩子（空组名）
	for _, hook := range settingUpdateHooks[""] {
		err := hook(group)
		if err != nil {
		}
	}
}
