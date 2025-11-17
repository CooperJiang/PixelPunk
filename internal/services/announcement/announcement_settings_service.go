package announcement

import (
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/services/setting"
)

/* GetAnnouncementSettings 获取公告系统配置 */
func GetAnnouncementSettings() (map[string]interface{}, error) {
	settingsDTO, err := setting.GetSettingsByGroupAsMap("announcement")
	if err != nil {
		return nil, fmt.Errorf("获取公告配置失败: %v", err)
	}

	return settingsDTO.Settings, nil
}

/* UpdateAnnouncementSettings 更新公告系统配置 */
func UpdateAnnouncementSettings(settings map[string]interface{}) error {
	var settingsToUpdate []dto.SettingCreateDTO

	allowedKeys := map[string]string{
		"announcement_enabled":         "boolean",
		"announcement_drawer_position": "string",
		"announcement_drawer_width":    "number",
		"announcement_display_limit":   "number",
		"announcement_auto_show_delay": "number",
	}

	for key, valueType := range allowedKeys {
		if value, exists := settings[key]; exists {
			settingsToUpdate = append(settingsToUpdate, dto.SettingCreateDTO{
				Key:         key,
				Value:       value,
				Type:        valueType,
				Group:       "announcement",
				Description: getSettingDescription(key),
				IsSystem:    true,
			})
		}
	}

	if len(settingsToUpdate) == 0 {
		return fmt.Errorf("没有有效的配置项需要更新")
	}

	batchDTO := &dto.BatchUpsertSettingDTO{
		Settings: settingsToUpdate,
	}

	_, err := setting.BatchUpsertSettings(batchDTO)
	if err != nil {
		return fmt.Errorf("更新公告配置失败: %v", err)
	}

	return nil
}

/* InitializeAnnouncementSettings 初始化公告系统配置 */
func InitializeAnnouncementSettings() error {
	existingSettings, err := GetAnnouncementSettings()
	if err == nil && len(existingSettings) > 0 {
		// 配置已存在，不需要初始化
		return nil
	}

	defaultSettings := []dto.SettingCreateDTO{
		{
			Key:         "announcement_enabled",
			Value:       true,
			Type:        "boolean",
			Group:       "announcement",
			Description: "公告系统总开关",
			IsSystem:    true,
		},
		{
			Key:         "announcement_drawer_position",
			Value:       "right",
			Type:        "string",
			Group:       "announcement",
			Description: "抽屉位置(left/right)",
			IsSystem:    true,
		},
		{
			Key:         "announcement_drawer_width",
			Value:       600,
			Type:        "number",
			Group:       "announcement",
			Description: "抽屉宽度(px)",
			IsSystem:    true,
		},
		{
			Key:         "announcement_display_limit",
			Value:       10,
			Type:        "number",
			Group:       "announcement",
			Description: "显示最近N条公告",
			IsSystem:    true,
		},
		{
			Key:         "announcement_auto_show_delay",
			Value:       3,
			Type:        "number",
			Group:       "announcement",
			Description: "置顶公告自动弹窗延迟时间（秒）",
			IsSystem:    true,
		},
	}

	batchDTO := &dto.BatchUpsertSettingDTO{
		Settings: defaultSettings,
	}

	_, err = setting.BatchUpsertSettings(batchDTO)
	if err != nil {
		return fmt.Errorf("初始化公告配置失败: %v", err)
	}

	return nil
}

/* getSettingDescription 获取配置项描述 */
func getSettingDescription(key string) string {
	descriptions := map[string]string{
		"announcement_enabled":         "公告系统总开关",
		"announcement_drawer_position": "抽屉位置(left/right)",
		"announcement_drawer_width":    "抽屉宽度(px)",
		"announcement_display_limit":   "显示最近N条公告",
		"announcement_auto_show_delay": "置顶公告自动弹窗延迟时间（秒）",
	}

	if desc, exists := descriptions[key]; exists {
		return desc
	}
	return ""
}
