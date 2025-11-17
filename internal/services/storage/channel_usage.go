package storage

import (
	"fmt"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
)

func GetTotalStorageUsageByChannels() (int64, string, error) {
	db := database.GetDB()

	var channels []models.StorageChannel
	if err := db.Where("status = ?", 1).Find(&channels).Error; err != nil {
		return 0, "", fmt.Errorf("获取存储渠道列表失败: %v", err)
	}

	var totalStorage int64

	for _, channel := range channels {
		var channelStorage *int64
		if err := db.Model(&models.File{}).
			Where("storage_provider_id = ? AND status <> ?", channel.ID, "pending_deletion").
			Select("SUM(size)").Row().Scan(&channelStorage); err != nil {
			logger.Error("获取渠道 %s 存储使用量失败: %v", channel.Name, err)
			continue
		}

		if channelStorage != nil {
			totalStorage += *channelStorage
		}
	}

	formattedStorage := utils.FormatBytes(totalStorage)

	return totalStorage, formattedStorage, nil
}

func RefreshChannelCache(channelID string) error {
	mgr, err := createStorageManager()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "创建存储管理器失败")
	}

	if err := mgr.RefreshAdapter(channelID); err != nil {
		return errors.Wrap(err, errors.CodeInternal, "刷新渠道缓存失败")
	}

	return nil
}

func ClearAllChannelCache() error {
	mgr, err := createStorageManager()
	if err != nil {
		return errors.Wrap(err, errors.CodeInternal, "创建存储管理器失败")
	}

	mgr.ClearCache()

	return nil
}
