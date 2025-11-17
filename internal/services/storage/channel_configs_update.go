package storage

/* Update configs helper split from storage_channel_service.go (no behavior change). */

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

func UpdateChannelConfigs(channelID string, configs map[string]interface{}) error {
	db := database.GetDB()

	var channel models.StorageChannel
	if err := db.First(&channel, "id = ?", channelID).Error; err != nil {
		return errors.New(errors.CodeNotFound, "存储渠道不存在")
	}

	templates, ok := models.StorageConfigTemplates[channel.Type]
	if !ok {
		return errors.New(errors.CodeInvalidParameter, "不支持的存储类型")
	}

	templateMap := make(map[string]models.ConfigTemplate)
	for _, tmpl := range templates {
		templateMap[tmpl.KeyName] = tmpl
	}

	keys := make([]string, 0, len(configs))
	for key := range configs {
		keys = append(keys, key)
	}

	var existingConfigs []models.StorageConfigItem
	db.Where("channel_id = ? AND key_name IN ?", channelID, keys).Find(&existingConfigs)

	existingMap := make(map[string]*models.StorageConfigItem)
	for i := range existingConfigs {
		existingMap[existingConfigs[i].KeyName] = &existingConfigs[i]
	}

	toInsert := make([]models.StorageConfigItem, 0)
	toUpdate := make([]models.StorageConfigItem, 0)

	for key, value := range configs {
		tmpl, exists := templateMap[key]
		if !exists {
			return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("无效的配置项: %s", key))
		}

		var strValue string
		switch tmpl.Type {
		case "int":
			intValue, ok := value.(float64) // JSON解析数字默认为float64
			if !ok {
				if strVal, ok := value.(string); ok {
					if val, err := strconv.Atoi(strVal); err == nil {
						intValue = float64(val)
					} else {
						return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("配置项 %s 的值必须为整数", key))
					}
				} else {
					return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("配置项 %s 的值必须为整数", key))
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
						return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("配置项 %s 的值必须为布尔类型", key))
					}
				} else {
					return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("配置项 %s 的值必须为布尔类型", key))
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

		if existing, exists := existingMap[key]; exists {
			existing.Value = strValue
			existing.UpdatedAt = common.JSONTimeNow()
			toUpdate = append(toUpdate, *existing)
		} else {
			newItem := models.StorageConfigItem{
				ID:          strings.ReplaceAll(uuid.New().String(), "-", ""),
				ChannelID:   channelID,
				Name:        tmpl.Name,
				KeyName:     key,
				Type:        tmpl.Type,
				IsSecret:    tmpl.IsSecret,
				Required:    tmpl.Required,
				Description: tmpl.Description,
				Value:       strValue,
				CreatedAt:   common.JSONTimeNow(),
				UpdatedAt:   common.JSONTimeNow(),
			}
			toInsert = append(toInsert, newItem)
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if len(toInsert) > 0 {
			if err := tx.Create(&toInsert).Error; err != nil {
				return err
			}
		}

		for _, item := range toUpdate {
			if err := tx.Model(&models.StorageConfigItem{}).
				Where("id = ?", item.ID).
				Updates(map[string]interface{}{
					"value":      item.Value,
					"updated_at": item.UpdatedAt,
				}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
