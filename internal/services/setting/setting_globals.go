package setting

var settingService *SettingService

var settingChangeHandlers map[string]func(value string)

const (
	SettingCachePrefix      = "setting:"
	SettingGroupPrefix      = "setting:group:"
	SettingCacheExpire      = 0
	SettingGroupCacheExpire = 0
)

type SettingService struct{}
