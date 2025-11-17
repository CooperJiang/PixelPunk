package storage

/* Repository and read helpers split from storage_channel_service.go (no behavior change). */

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

type testChannelRepository struct{}

func (r *testChannelRepository) GetChannel(channelID string) (*models.StorageChannel, error) {
	return GetChannelByID(channelID)
}

func (r *testChannelRepository) GetChannelConfig(channelID string) (map[string]interface{}, error) {
	return GetChannelConfigMap(channelID)
}

func (r *testChannelRepository) GetDefaultChannel() (*models.StorageChannel, error) {
	return GetDefaultChannel()
}

func (r *testChannelRepository) GetActiveChannels() ([]*models.StorageChannel, error) {
	channels, err := GetAllChannels()
	if err != nil {
		return nil, err
	}
	var activeChannels []*models.StorageChannel
	for i := range channels {
		if channels[i].Status == 1 {
			activeChannels = append(activeChannels, &channels[i])
		}
	}
	return activeChannels, nil
}

func GetAllChannels() ([]models.StorageChannel, error) {
	var channels []models.StorageChannel
	if err := database.GetDB().Find(&channels).Error; err != nil {
		return nil, err
	}

	for i := range channels {
		var imageCount int64
		if err := database.GetDB().Model(&models.File{}).
			Where("storage_provider_id = ?", channels[i].ID).
			Count(&imageCount).Error; err != nil {
			imageCount = 0
		}
		channels[i].FileCount = imageCount

		configMap, err := GetChannelConfigMap(channels[i].ID)
		if err == nil {
			if customDomain, exists := configMap["custom_domain"]; exists {
				if domainStr, ok := customDomain.(string); ok && domainStr != "" {
					channels[i].CustomDomain = domainStr
				}
			}

			if bucket, exists := configMap["bucket"]; exists {
				if bucketStr, ok := bucket.(string); ok && bucketStr != "" {
					channels[i].Bucket = bucketStr
				}
			}
		}

		var lastUpload models.File
		if err := database.GetDB().Model(&models.File{}).
			Where("storage_provider_id = ?", channels[i].ID).
			Order("created_at DESC").
			First(&lastUpload).Error; err == nil {
			channels[i].LastUploadAt = &lastUpload.CreatedAt
		}
	}

	return channels, nil
}

func GetChannelByID(channelID string) (*models.StorageChannel, error) {
	var channel models.StorageChannel
	if err := database.GetDB().First(&channel, "id = ?", channelID).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func GetChannelConfigs(channelID string) ([]models.StorageConfigItem, error) {
	channel, err := GetChannelByID(channelID)
	if err != nil {
		return nil, err
	}

	templates, ok := models.StorageConfigTemplates[channel.Type]
	if !ok {
		return nil, fmt.Errorf("不支持的存储类型: %s", channel.Type)
	}

	err = EnsureChannelConfigItems(channelID, channel.Type)
	if err != nil {
	}

	var dbConfigs []models.StorageConfigItem
	if err := database.GetDB().Where("channel_id = ?", channelID).Find(&dbConfigs).Error; err != nil {
		return nil, err
	}

	configMap := make(map[string]models.StorageConfigItem)
	for _, config := range dbConfigs {
		configMap[config.KeyName] = config
	}

	var result []models.StorageConfigItem
	for _, tmpl := range templates {
		if existingConfig, exists := configMap[tmpl.KeyName]; exists {
			if existingConfig.IsSecret && existingConfig.Value != "" {
				existingConfig.Value = "******" // 脱敏
			}
			result = append(result, existingConfig)
		} else {
			defaultConfig := models.StorageConfigItem{
				ID:          "", // 新配置项ID为空
				ChannelID:   channelID,
				Name:        tmpl.Name,
				KeyName:     tmpl.KeyName,
				Type:        tmpl.Type,
				IsSecret:    tmpl.IsSecret,
				Required:    tmpl.Required,
				Description: tmpl.Description,
				Value:       "", // 默认空值
			}
			result = append(result, defaultConfig)
		}
	}

	return result, nil
}

func EnsureChannelConfigItems(channelID, channelType string) error {
	db := database.GetDB()

	templates, ok := models.StorageConfigTemplates[channelType]
	if !ok {
		return fmt.Errorf("不支持的存储类型: %s", channelType)
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, tmpl := range templates {
			var count int64
			if err := tx.Model(&models.StorageConfigItem{}).
				Where("channel_id = ? AND key_name = ?", channelID, tmpl.KeyName).
				Count(&count).Error; err != nil {
				return err
			}

			if count == 0 {
				configItemUUID := uuid.New().String()
				configItemUUID = strings.ReplaceAll(configItemUUID, "-", "")

				configItem := models.StorageConfigItem{
					ID:          configItemUUID,
					ChannelID:   channelID,
					Name:        tmpl.Name,
					KeyName:     tmpl.KeyName,
					Type:        tmpl.Type,
					IsSecret:    tmpl.IsSecret,
					Required:    tmpl.Required,
					Description: tmpl.Description,
					Value:       "", // 默认空值
					CreatedAt:   common.JSONTimeNow(),
					UpdatedAt:   common.JSONTimeNow(),
				}

				if err := tx.Create(&configItem).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func GetChannelConfigsForExport(channelID string) ([]models.StorageConfigItem, error) {
	channel, err := GetChannelByID(channelID)
	if err != nil {
		return nil, err
	}

	templates, ok := models.StorageConfigTemplates[channel.Type]
	if !ok {
		return nil, fmt.Errorf("不支持的存储类型: %s", channel.Type)
	}

	err = EnsureChannelConfigItems(channelID, channel.Type)
	if err != nil {
	}

	var dbConfigs []models.StorageConfigItem
	if err := database.GetDB().Where("channel_id = ?", channelID).Find(&dbConfigs).Error; err != nil {
		return nil, err
	}

	configMap := make(map[string]models.StorageConfigItem)
	for _, config := range dbConfigs {
		configMap[config.KeyName] = config
	}

	var result []models.StorageConfigItem
	for _, tmpl := range templates {
		if existingConfig, exists := configMap[tmpl.KeyName]; exists {
			result = append(result, existingConfig)
		} else {
			defaultConfig := models.StorageConfigItem{
				ID:          "", // 新配置项ID为空
				ChannelID:   channelID,
				Name:        tmpl.Name,
				KeyName:     tmpl.KeyName,
				Type:        tmpl.Type,
				IsSecret:    tmpl.IsSecret,
				Required:    tmpl.Required,
				Description: tmpl.Description,
				Value:       "", // 默认空值
			}
			result = append(result, defaultConfig)
		}
	}

	return result, nil
}

func GetChannelConfigMap(channelID string) (map[string]interface{}, error) {
	configs, err := GetChannelConfigs(channelID)
	if err != nil {
		return nil, err
	}

	channel, err := GetChannelByID(channelID)
	if err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for _, config := range configs {
		value, err := config.GetDecryptedValue()
		if err != nil {
			return nil, err
		}

		switch config.Type {
		case "int":
			if intVal, err := strconv.Atoi(value); err == nil {
				result[config.KeyName] = intVal
			} else {
				result[config.KeyName] = value
			}
		case "bool":
			if boolVal, err := strconv.ParseBool(value); err == nil {
				result[config.KeyName] = boolVal
			} else {
				result[config.KeyName] = value
			}
		default:
			result[config.KeyName] = value
		}
	}

	result = models.GetChannelConfigWithDefaults(channelID, result, channel.Type)

	return result, nil
}

func GetDefaultChannel() (*models.StorageChannel, error) {
	db := database.GetDB()
	var channel models.StorageChannel
	if err := db.Where("is_default = ? AND status = ?", true, 1).First(&channel).Error; err != nil {
		return nil, errors.New(errors.CodeNotFound, "未找到默认存储渠道，请在管理后台设置默认渠道")
	}
	return &channel, nil
}
