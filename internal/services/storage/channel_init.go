package storage

import (
	"fmt"
	"strings"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateDefaultLocalChannel() (*models.StorageChannel, error) {
	db := database.GetDB()

	uuidStr := uuid.New().String()
	uuidStr = strings.ReplaceAll(uuidStr, "-", "") // 移除连字符，生成32位ID

	newChannel := &models.StorageChannel{
		ID:        uuidStr,
		Name:      "本地存储",
		Type:      "local",
		Status:    1,
		IsDefault: true,
		IsLocal:   true, // 设置为本地存储渠道
		Remark:    "系统自动创建的默认存储渠道",
		CreatedAt: common.JSONTimeNow(),
		UpdatedAt: common.JSONTimeNow(),
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(newChannel).Error; err != nil {
			return err
		}

		configItems := []models.StorageConfigItem{
			{
				ID:          strings.ReplaceAll(uuid.New().String(), "-", ""),
				ChannelID:   newChannel.ID,
				Name:        "存储路径",
				KeyName:     "base_path",
				Type:        "string",
				Value:       "uploads/files",
				Required:    true,
				Description: "本地存储的基础路径",
				CreatedAt:   common.JSONTimeNow(),
				UpdatedAt:   common.JSONTimeNow(),
			},
			{
				ID:          strings.ReplaceAll(uuid.New().String(), "-", ""),
				ChannelID:   newChannel.ID,
				Name:        "缩略图路径",
				KeyName:     "thumbnail_path",
				Type:        "string",
				Value:       "uploads/thumbnails",
				Required:    true,
				Description: "文件缩略图存储路径",
				CreatedAt:   common.JSONTimeNow(),
				UpdatedAt:   common.JSONTimeNow(),
			},
		}

		for _, item := range configItems {
			if err := tx.Create(&item).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("创建默认存储渠道失败: %v", err)
	}

	return newChannel, nil
}


func MigrateLocalChannels() error {
	return database.GetDB().Model(&models.StorageChannel{}).
		Where("type = ? AND is_local = ?", "local", false).
		Update("is_local", true).Error
}

func CheckAndInitDefaultChannel() error {
	db := database.GetDB()

	if err := MigrateLocalChannels(); err != nil {
		logger.Error(fmt.Sprintf("迁移本地存储渠道失败: %v", err))
	}

	var count int64
	if err := db.Model(&models.StorageChannel{}).Where("status = ?", 1).Count(&count).Error; err != nil {
		return fmt.Errorf("检查存储渠道失败: %v", err)
	}

	if count == 0 {
		_, err := CreateDefaultLocalChannel()
		if err != nil {
			return fmt.Errorf("创建默认存储渠道失败: %v", err)
		}
	}

	return nil
}

func InitStorageService() {
	if err := MigrateLocalChannels(); err != nil {
		logger.Error(fmt.Sprintf("迁移本地存储渠道失败: %v", err))
	}

	db := database.GetDB()
	var count int64

	if err := db.Model(&models.StorageChannel{}).Where("status = ?", 1).Count(&count).Error; err != nil {
		return
	}

	if count == 0 {
		_, err := CreateDefaultLocalChannel()
		if err != nil {
			return
		}
	}
}
