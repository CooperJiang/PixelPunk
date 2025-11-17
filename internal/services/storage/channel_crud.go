package storage

/* CRUD operations split from storage_channel_service.go (no behavior change). */

import (
	"fmt"
	"strconv"
	"strings"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateChannel(channel *models.StorageChannel, configs map[string]interface{}) error {
	db := database.GetDB()

	var count int64
	if err := db.Model(&models.StorageChannel{}).Where("type = ?", channel.Type).Count(&count).Error; err != nil {
		return fmt.Errorf("检查存储渠道失败: %v", err)
	}

	if channel.Type == "local" && count > 0 {
		return errors.New(errors.CodeValidationFailed, "本地存储渠道已存在，不能重复创建")
	}

	if channel.Type == "local" {
		channel.IsLocal = true
	}

	uuidStr := uuid.New().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "") // 移除连字符，生成32位ID
	channel.ID = uuidStr

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(channel).Error; err != nil {
			return err
		}

		templates, ok := models.StorageConfigTemplates[channel.Type]
		if !ok {
			return fmt.Errorf("不支持的存储类型: %s", channel.Type)
		}

		for _, tmpl := range templates {
			configItemUUID := uuid.New().String()
			configItemUUID = strings.ReplaceAll(configItemUUID, "-", "")

			configItem := models.StorageConfigItem{
				ID:          configItemUUID,
				ChannelID:   channel.ID,
				Name:        tmpl.Name,
				KeyName:     tmpl.KeyName,
				Type:        tmpl.Type,
				IsSecret:    tmpl.IsSecret,
				Required:    tmpl.Required,
				Description: tmpl.Description,
				CreatedAt:   common.JSONTimeNow(),
				UpdatedAt:   common.JSONTimeNow(),
			}

			if configs != nil {
				if value, exists := configs[tmpl.KeyName]; exists {
					var strValue string
					switch tmpl.Type {
					case "int":
						intValue, ok := value.(float64) // JSON解析数字默认为float64
						if !ok {
							if strVal, ok := value.(string); ok {
								if val, err := strconv.Atoi(strVal); err == nil {
									intValue = float64(val)
								} else {
									return fmt.Errorf("无效的整数值: %v", value)
								}
							} else {
								return fmt.Errorf("无效的整数值: %v", value)
							}
						}
						strValue = fmt.Sprintf("%d", int(intValue))

					case "bool":
						boolValue, ok := value.(bool)
						if !ok {
							if strVal, ok := value.(string); ok {
								if val, err := strconv.ParseBool(strVal); err == nil {
									boolValue = val
								} else {
									return fmt.Errorf("无效的布尔值: %v", value)
								}
							} else {
								return fmt.Errorf("无效的布尔值: %v", value)
							}
						}
						strValue = fmt.Sprintf("%t", boolValue)

					default: // string, password等
						if str, ok := value.(string); ok {
							strValue = str
						} else {
							strValue = fmt.Sprintf("%v", value)
						}
					}

					configItem.Value = strValue
				}
			}

			if err := tx.Create(&configItem).Error; err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func UpdateChannel(channel *models.StorageChannel) error {
	db := database.GetDB()

	var existingChannel models.StorageChannel
	if err := db.First(&existingChannel, "id = ?", channel.ID).Error; err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if existingChannel.IsDefault && !channel.IsDefault {
			var otherDefaultCount int64
			if err := tx.Model(&models.StorageChannel{}).
				Where("id <> ? AND is_default = ?", channel.ID, true).
				Count(&otherDefaultCount).Error; err != nil {
				return err
			}
			if otherDefaultCount == 0 {
				return errors.New(errors.CodeValidationFailed, "至少需要一个默认渠道，不能取消当前默认渠道")
			}
		}

		if existingChannel.IsLocal {
			channel.Status = 1
			channel.IsDefault = true
			channel.IsLocal = true
		}

		if err := tx.Save(channel).Error; err != nil {
			return err
		}

		if channel.IsDefault {
			if err := tx.Model(&models.StorageChannel{}).
				Where("id <> ? AND is_default = ?", channel.ID, true).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func DeleteChannel(channelID string) error {
	db := database.GetDB()

	var channel models.StorageChannel
	if err := db.First(&channel, "id = ?", channelID).Error; err != nil {
		return err
	}

	if channel.IsLocal {
		return errors.New(errors.CodeValidationFailed, "本地存储渠道不允许删除，只能修改")
	}

	var count int64
	if err := db.Model(&models.File{}).Where("storage_provider_id = ?", channelID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New(errors.CodeValidationFailed, fmt.Sprintf("该存储渠道仍有%d个文件使用，请先删除关联文件后移除渠道。", count))
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("channel_id = ?", channelID).Delete(&models.StorageConfigItem{}).Error; err != nil {
			return err
		}

		if err := tx.Delete(&models.StorageChannel{}, "id = ?", channelID).Error; err != nil {
			return err
		}

		return nil
	})
}

func EnableChannel(channelID string) error {
	return database.GetDB().Model(&models.StorageChannel{}).Where("id = ?", channelID).Updates(map[string]interface{}{
		"status":     1,
		"updated_at": common.JSONTimeNow(),
	}).Error
}

func DisableChannel(channelID string) error {
	db := database.GetDB()

	var channel models.StorageChannel
	if err := db.First(&channel, "id = ?", channelID).Error; err != nil {
		return err
	}

	if channel.IsLocal {
		return errors.New(errors.CodeValidationFailed, "本地存储渠道不能被禁用")
	}

	if channel.IsDefault {
		return errors.New(errors.CodeValidationFailed, "默认渠道不能被禁用")
	}

	return db.Model(&models.StorageChannel{}).Where("id = ?", channelID).Updates(map[string]interface{}{
		"status":     0,
		"updated_at": common.JSONTimeNow(),
	}).Error
}

func SetDefaultChannel(channelID string) error {
	db := database.GetDB()

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.StorageChannel{}).
			Where("is_default = ?", true).
			Updates(map[string]interface{}{
				"is_default": false,
				"updated_at": common.JSONTimeNow(),
			}).Error; err != nil {
			return err
		}

		return tx.Model(&models.StorageChannel{}).Where("id = ?", channelID).Updates(map[string]interface{}{
			"is_default": true,
			"status":     1, // 确保默认渠道是启用的
			"updated_at": common.JSONTimeNow(),
		}).Error
	})
}
