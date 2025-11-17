package user

import (
	"fmt"
	"pixelpunk/internal/controllers/user/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/bandwidth"
	"pixelpunk/pkg/database"
)

func formatBytes(bytes int64) string {
	const (
		B  = 1
		KB = 1024 * B
		MB = 1024 * KB
		GB = 1024 * MB
		TB = 1024 * GB
	)

	if bytes == 0 {
		return "0 B"
	}

	var size float64
	var unit string

	switch {
	case bytes >= TB:
		size = float64(bytes) / TB
		unit = "TB"
	case bytes >= GB:
		size = float64(bytes) / GB
		unit = "GB"
	case bytes >= MB:
		size = float64(bytes) / MB
		unit = "MB"
	case bytes >= KB:
		size = float64(bytes) / KB
		unit = "KB"
	default:
		size = float64(bytes)
		unit = "B"
	}

	if size >= 100 {
		return fmt.Sprintf("%.0f %s", size, unit)
	} else if size >= 10 {
		return fmt.Sprintf("%.1f %s", size, unit)
	} else {
		return fmt.Sprintf("%.2f %s", size, unit)
	}
}

func GetWorkspaceStats(userID uint) (*dto.WorkspaceStatsDTO, error) {
	db := database.GetDB()
	response := &dto.WorkspaceStatsDTO{}

	userSettings, err := GetUserSettings(userID)
	if err != nil {
		return nil, fmt.Errorf("获取用户设置失败: %w", err)
	}

	response.TotalStorage = userSettings.StorageLimit
	response.TotalStorageFormatted = formatBytes(userSettings.StorageLimit)

	response.TotalBandwidth = userSettings.BandwidthLimit
	response.TotalBandwidthFormatted = formatBytes(userSettings.BandwidthLimit)

	var totalImages int64
	var totalSize int64

	err = db.Model(&models.File{}).Where("user_id = ?", userID).Count(&totalImages).Error
	if err != nil {
		return nil, fmt.Errorf("统计文件数量失败: %w", err)
	}

	err = db.Model(&models.File{}).Where("user_id = ?", userID).
		Select("COALESCE(SUM(size), 0)").Scan(&totalSize).Error
	if err != nil {
		return nil, fmt.Errorf("统计存储大小失败: %w", err)
	}

	response.TotalImages = int(totalImages)
	response.UsedStorage = totalSize
	response.UsedStorageFormatted = formatBytes(totalSize)

	var totalViews int64
	if totalImages > 0 {
		var imageIDs []string
		err = db.Model(&models.File{}).Where("user_id = ?", userID).Pluck("id", &imageIDs).Error
		if err != nil {
			return nil, fmt.Errorf("获取文件ID列表失败: %w", err)
		}

		if len(imageIDs) > 0 {
			err = db.Model(&models.FileStats{}).
				Where("file_id IN ?", imageIDs).
				Select("COALESCE(SUM(views), 0)").Scan(&totalViews).Error
			if err != nil {
				totalViews = 0
			}
		}
	}
	response.TotalViews = totalViews

	var totalShares int64
	err = db.Model(&models.Share{}).Where("user_id = ?", userID).Count(&totalShares).Error
	if err != nil {
		totalShares = 0
	}
	response.TotalShares = int(totalShares)

	bandwidthStats, err := bandwidth.GetUserBandwidthStats(userID)
	if err != nil {
		return nil, fmt.Errorf("获取带宽统计失败: %w", err)
	}

	var usedBandwidth int64

	if bandwidthStats.UsedBytes > 0 {
		usedBandwidth = bandwidthStats.UsedBytes
	} else {
		if totalViews > 0 && totalSize > 0 {
			avgImageSize := totalSize / totalImages
			if avgImageSize > 0 {
				usedBandwidth = avgImageSize * totalViews
			}
		}
	}

	response.UsedBandwidth = usedBandwidth
	response.UsedBandwidthFormatted = formatBytes(usedBandwidth)

	return response, nil
}
