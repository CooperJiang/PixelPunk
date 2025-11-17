package user

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
)

func InitUserSettings(userID uint) error {
	var count int64
	if err := database.DB.Model(&models.UserSettings{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "检查用户设置记录失败")
	}

	if count > 0 {
		return nil
	}

	settings := models.UserSettings{
		UserID:             userID,
		StorageLimit:       models.DefaultStorageLimit,
		BandwidthLimit:     models.DefaultBandwidthLimit,
		DefaultAccessLevel: "private",
		OptimizeImages:     true,
		CreatedAt:          common.JSONTimeNow(),
		UpdatedAt:          common.JSONTimeNow(),
	}

	if err := database.DB.Create(&settings).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "创建用户设置记录失败")
	}

	return nil
}

func GetUserSettings(userID uint) (*models.UserSettings, error) {
	var settings models.UserSettings

	err := database.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := InitUserSettings(userID); err != nil {
				return nil, err
			}
			err = database.DB.Where("user_id = ?", userID).First(&settings).Error
			if err != nil {
				return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取用户设置失败")
			}
		} else {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取用户设置失败")
		}
	}

	return &settings, nil
}

func UpdateUserSettings(userID uint, storageLimit, bandwidthLimit int64, defaultAccessLevel string, optimizeImages bool) (*models.UserSettings, error) {
	settings, err := GetUserSettings(userID)
	if err != nil {
		return nil, err
	}

	if storageLimit > 0 {
		settings.StorageLimit = storageLimit
	}
	if bandwidthLimit > 0 {
		settings.BandwidthLimit = bandwidthLimit
	}
	if defaultAccessLevel != "" {
		settings.DefaultAccessLevel = defaultAccessLevel
	}
	settings.OptimizeImages = optimizeImages
	settings.UpdatedAt = common.JSONTimeNow()

	if err := database.DB.Save(settings).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新用户设置失败")
	}

	return settings, nil
}
