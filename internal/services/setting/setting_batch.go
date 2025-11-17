package setting

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/hooks"
	"strings"
	"time"

	"gorm.io/gorm"
)

/* normalizeSiteBaseURL 规范化网站基础地址 */
func normalizeSiteBaseURL(url string) string {
	url = strings.TrimSpace(url)
	// 只移除末尾的斜杠，保留协议前缀（https:// 或 http://）
	url = strings.TrimSuffix(url, "/")
	return url
}

/* normalizeSettingValue 根据 key 规范化设置值 */
func normalizeSettingValue(key string, value interface{}) interface{} {
	// 对 site_base_url 进行特殊处理
	if key == "site_base_url" {
		if strValue, ok := value.(string); ok {
			return normalizeSiteBaseURL(strValue)
		}
	}
	return value
}

/* BatchCreateSettings 批量创建设置 */
func BatchCreateSettings(createDTOs *dto.BatchSettingCreateDTO) (*dto.BatchSettingResponseDTO, error) {
	db := database.GetDB()
	result := &dto.BatchSettingResponseDTO{Success: []dto.SettingResponseDTO{}, Failed: []dto.BatchFailedItem{}}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, createDTO := range createDTOs.Settings {
			var count int64
			if err := tx.Model(&models.Setting{}).Where("`key` = ?", createDTO.Key).Count(&count).Error; err != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: createDTO.Key, Message: "检查设置键名失败: " + err.Error()})
				continue
			}
			if count > 0 {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: createDTO.Key, Message: "设置键名已存在"})
				continue
			}

			normalizedValue := normalizeSettingValue(createDTO.Key, createDTO.Value)

			var valueToStore string
			var valueJSON []byte
			var jsonErr error
			if createDTO.Type == models.SettingTypeString {
				if strValue, ok := normalizedValue.(string); ok {
					valueJSON, jsonErr = json.Marshal(strValue)
				} else {
					valueJSON, jsonErr = json.Marshal(fmt.Sprintf("%v", normalizedValue))
				}
			} else {
				valueJSON, jsonErr = json.Marshal(normalizedValue)
			}
			if jsonErr != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: createDTO.Key, Message: "设置值格式错误: " + jsonErr.Error()})
				continue
			}
			valueToStore = string(valueJSON)

			setting := models.Setting{
				Key:         createDTO.Key,
				Value:       valueToStore,
				Type:        createDTO.Type,
				Group:       createDTO.Group,
				Description: createDTO.Description,
				IsSystem:    createDTO.IsSystem,
			}
			if err := tx.Create(&setting).Error; err != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: createDTO.Key, Message: "创建设置失败: " + err.Error()})
				continue
			}

			result.Success = append(result.Success, dto.SettingResponseDTO{
				ID:          setting.ID,
				Key:         setting.Key,
				Value:       normalizedValue,
				Type:        setting.Type,
				Group:       setting.Group,
				Description: setting.Description,
				IsSystem:    setting.IsSystem,
				CreatedAt:   time.Time(setting.CreatedAt).Format(time.RFC3339),
				UpdatedAt:   time.Time(setting.UpdatedAt).Format(time.RFC3339),
			})
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBCommitFailed, "批量创建设置失败")
	}

	updatedGroups := make(map[string]bool)
	for _, item := range result.Success {
		updatedGroups[item.Group] = true
		setSettingToCache(item.Key, &item)
		deleteSettingGroupFromCache(item.Group)
		var valueStr string
		if s, ok := item.Value.(string); ok {
			valueStr = s
		} else if data, err := json.Marshal(item.Value); err == nil {
			valueStr = string(data)
		}
		notifySettingChanged(item.Group, item.Key, valueStr)
	}
	for group := range updatedGroups {
		hooks.TriggerSettingUpdate(group)
	}
	return result, nil
}

/* BatchUpsertSettings 批量更新或创建设置 */
func BatchUpsertSettings(upsertDTOs *dto.BatchUpsertSettingDTO) (*dto.BatchSettingResponseDTO, error) {
	db := database.GetDB()
	result := &dto.BatchSettingResponseDTO{Success: []dto.SettingResponseDTO{}, Failed: []dto.BatchFailedItem{}}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, settingDTO := range upsertDTOs.Settings {
			var setting models.Setting
			err := tx.Where("`key` = ?", settingDTO.Key).First(&setting).Error

			normalizedValue := normalizeSettingValue(settingDTO.Key, settingDTO.Value)

			var valueToStore string
			var valueJSON []byte
			var jsonErr error
			if settingDTO.Type == models.SettingTypeString {
				if strValue, ok := normalizedValue.(string); ok {
					valueJSON, jsonErr = json.Marshal(strValue)
				} else {
					valueJSON, jsonErr = json.Marshal(fmt.Sprintf("%v", normalizedValue))
				}
			} else {
				valueJSON, jsonErr = json.Marshal(normalizedValue)
			}
			if jsonErr != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: settingDTO.Key, Message: "设置值格式错误: " + jsonErr.Error()})
				continue
			}
			valueToStore = string(valueJSON)

			if err == nil {
				if setting.IsSystem && setting.Group != settingDTO.Group {
					result.Failed = append(result.Failed, dto.BatchFailedItem{Key: settingDTO.Key, Message: "系统设置不能修改分组"})
					continue
				}
				setting.Value = valueToStore
				setting.Type = settingDTO.Type
				setting.Group = settingDTO.Group
				setting.Description = settingDTO.Description
				if err := tx.Save(&setting).Error; err != nil {
					result.Failed = append(result.Failed, dto.BatchFailedItem{Key: settingDTO.Key, Message: "更新设置失败: " + err.Error()})
					continue
				}
				result.Success = append(result.Success, dto.SettingResponseDTO{
					ID:          setting.ID,
					Key:         setting.Key,
					Value:       normalizedValue,
					Type:        setting.Type,
					Group:       setting.Group,
					Description: setting.Description,
					IsSystem:    setting.IsSystem,
					CreatedAt:   time.Time(setting.CreatedAt).Format(time.RFC3339),
					UpdatedAt:   time.Time(setting.UpdatedAt).Format(time.RFC3339),
				})
			} else {
				newSetting := models.Setting{
					Key:         settingDTO.Key,
					Value:       valueToStore,
					Type:        settingDTO.Type,
					Group:       settingDTO.Group,
					Description: settingDTO.Description,
					IsSystem:    settingDTO.IsSystem,
				}
				if err := tx.Create(&newSetting).Error; err != nil {
					result.Failed = append(result.Failed, dto.BatchFailedItem{Key: settingDTO.Key, Message: "创建设置失败: " + err.Error()})
					continue
				}
				result.Success = append(result.Success, dto.SettingResponseDTO{
					ID:          newSetting.ID,
					Key:         newSetting.Key,
					Value:       normalizedValue,
					Type:        newSetting.Type,
					Group:       newSetting.Group,
					Description: newSetting.Description,
					IsSystem:    newSetting.IsSystem,
					CreatedAt:   time.Time(newSetting.CreatedAt).Format(time.RFC3339),
					UpdatedAt:   time.Time(newSetting.UpdatedAt).Format(time.RFC3339),
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBCommitFailed, "批量更新或创建设置失败")
	}

	updatedGroups := make(map[string]bool)
	for _, item := range result.Success {
		updatedGroups[item.Group] = true
		invalidateSettingCaches(item.Group, item.Key)
		var valueStr string
		if s, ok := item.Value.(string); ok {
			valueStr = s
		} else if data, err := json.Marshal(item.Value); err == nil {
			valueStr = string(data)
		}
		notifySettingChanged(item.Group, item.Key, valueStr)
	}
	for group := range updatedGroups {
		hooks.TriggerSettingUpdate(group)
	}
	return result, nil
}

/* BatchUpdateSettings 批量更新设置（仅更新已存在项，不创建新项） */
func BatchUpdateSettings(updateDTOs *dto.BatchSettingUpdateDTO) (*dto.BatchSettingResponseDTO, error) {
	db := database.GetDB()
	result := &dto.BatchSettingResponseDTO{Success: []dto.SettingResponseDTO{}, Failed: []dto.BatchFailedItem{}}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	err := db.Transaction(func(tx *gorm.DB) error {
		for _, updateDTO := range updateDTOs.Settings {
			var setting models.Setting
			err := tx.Where("`key` = ?", updateDTO.Key).First(&setting).Error
			if err != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: updateDTO.Key, Message: "设置不存在"})
				continue
			}
			if setting.IsSystem && setting.Group != updateDTO.Group {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: updateDTO.Key, Message: "系统设置不能修改分组"})
				continue
			}

			normalizedValue := normalizeSettingValue(updateDTO.Key, updateDTO.Value)

			var valueJSON []byte
			var jsonErr error
			if updateDTO.Type == models.SettingTypeString {
				if strValue, ok := normalizedValue.(string); ok {
					valueJSON, jsonErr = json.Marshal(strValue)
				} else {
					valueJSON, jsonErr = json.Marshal(fmt.Sprintf("%v", normalizedValue))
				}
			} else {
				valueJSON, jsonErr = json.Marshal(normalizedValue)
			}
			if jsonErr != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: updateDTO.Key, Message: "设置值格式错误: " + jsonErr.Error()})
				continue
			}
			setting.Value = string(valueJSON)
			setting.Type = updateDTO.Type
			setting.Group = updateDTO.Group
			setting.Description = updateDTO.Description
			if err := tx.Save(&setting).Error; err != nil {
				result.Failed = append(result.Failed, dto.BatchFailedItem{Key: updateDTO.Key, Message: "更新设置失败: " + err.Error()})
				continue
			}
			result.Success = append(result.Success, dto.SettingResponseDTO{
				ID: setting.ID, Key: setting.Key, Value: normalizedValue, Type: setting.Type, Group: setting.Group, Description: setting.Description,
				IsSystem: setting.IsSystem, CreatedAt: time.Time(setting.CreatedAt).Format(time.RFC3339), UpdatedAt: time.Time(setting.UpdatedAt).Format(time.RFC3339),
			})
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBCommitFailed, "批量更新设置失败")
	}

	updatedGroups := make(map[string]bool)
	for _, item := range result.Success {
		updatedGroups[item.Group] = true
		invalidateSettingCaches(item.Group, item.Key)
		var valueStr string
		if s, ok := item.Value.(string); ok {
			valueStr = s
		} else if data, err := json.Marshal(item.Value); err == nil {
			valueStr = string(data)
		}
		notifySettingChanged(item.Group, item.Key, valueStr)
	}
	for group := range updatedGroups {
		hooks.TriggerSettingUpdate(group)
	}

	return result, nil
}
