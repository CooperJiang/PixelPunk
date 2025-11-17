package stats

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/utils"

	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func GetUserStats(userID uint) (*models.UserStatsResponse, error) {
	var stats models.UserUsageStats
	err := database.DB.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		stats = models.UserUsageStats{
			UserID: userID,
		}
	}

	var settings models.UserSettings
	err = database.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		settings = models.UserSettings{
			UserID:             userID,
			StorageLimit:       models.DefaultStorageLimit,
			BandwidthLimit:     models.DefaultBandwidthLimit,
			DefaultAccessLevel: "private",
		}
	}

	response := &models.UserStatsResponse{
		Storage: models.StorageStats{
			Used:       stats.TotalSize,
			Limit:      settings.StorageLimit,
			Percentage: stats.StorageUsedPercent(settings.StorageLimit),
		},
		Bandwidth: models.BandwidthStats{
			Used:       stats.TotalBandwidth,
			Limit:      settings.BandwidthLimit,
			Percentage: float64(stats.TotalBandwidth) / float64(settings.BandwidthLimit) * 100,
		},
		Files: models.FilesStats{
			Total: int64(stats.TotalImages),
			Views: stats.TotalViews,
		},
		Settings: settings,
	}

	return response, nil
}

func GetLatestFiles(limit int) ([]models.LatestFileInfo, error) {
	cacheKey := fmt.Sprintf("dashboard:latest_files:%d", limit)

	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result []models.LatestFileInfo
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return result, nil
		}
	}

	var images []models.File

	err := database.DB.Model(&models.File{}).
		Where("status <> ?", "pending_deletion").
		Order("created_at DESC").
		Limit(limit).
		Find(&images).Error

	if err != nil {
		return nil, err
	}

	result := make([]models.LatestFileInfo, 0, len(images))
	for _, file := range images {
		var user models.User
		database.DB.Model(&models.User{}).First(&user, file.UserID)

		fullURL, fullThumbURL, _ := storage.GetFullURLs(file)

		aiInfo, _ := getFileAIInfo(file.ID)

		result = append(result, models.LatestFileInfo{
			ID:            file.ID, // 直接使用字符串ID，不需要转换
			FileName:      file.FileName,
			DisplayName:   file.DisplayName,
			Size:          file.Size,
			SizeFormatted: utils.FormatBytes(file.Size),
			Format:        file.Format,
			URL:           file.URL,
			ThumbnailURL:  file.ThumbURL,
			FullURL:       fullURL,
			FullThumbURL:  fullThumbURL,
			CreatedAt:     file.CreatedAt,
			User: models.FileOwner{
				ID:       user.ID,
				Username: user.Username,
			},
			AIInfo:          aiInfo,
			StorageDuration: file.StorageDuration,
			ExpiresAt:       (*common.JSONTime)(file.ExpiresAt),
			IsTimeLimited:   file.IsTimeLimitedStorage(),
		})
	}

	signer := utils.GetURLSigner()
	for i := range result {
		if result[i].ID != "" {
			result[i].FullURL = signer.SignFileURL(result[i].ID, utils.SIGNATURE_DURATION)
			result[i].FullThumbURL = signer.SignThumbURL(result[i].ID, utils.SIGNATURE_DURATION)
		}
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return result, nil
}

func GetDailyActiveUploadUsers(startDate, endDate time.Time) (map[string]int64, error) {
	result := make(map[string]int64)

	days := int(endDate.Sub(startDate).Hours()/24) + 1

	for i := 0; i < days; i++ {
		currDate := startDate.AddDate(0, 0, i)
		dateStr := currDate.Format("2006-01-02")
		result[dateStr] = 0
	}

	for i := 0; i < days; i++ {
		currDate := startDate.AddDate(0, 0, i)
		nextDate := currDate.AddDate(0, 0, 1)
		dateStr := currDate.Format("2006-01-02")

		var activeUsers int64
		err := database.DB.Model(&models.File{}).
			Where("created_at >= ? AND created_at < ?", currDate, nextDate).
			Distinct("user_id").
			Count(&activeUsers).Error

		if err != nil {
			return result, err
		}

		result[dateStr] = activeUsers
	}

	return result, nil
}

func calculateGrowthRate(oldValue, newValue int64) float64 {
	if oldValue == 0 {
		return 100.0 // 如果旧值为0，增长率为100%
	}
	rate := float64(newValue-oldValue) / float64(oldValue) * 100
	return float64(int(rate*10+0.5)) / 10
}

func getFileAIInfo(fileID string) (interface{}, error) {
	var aiInfo models.FileAIInfo

	if err := database.DB.Where("file_id = ?", fileID).First(&aiInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var tagNames []string
	if len(aiInfo.Tags) > 0 {
		if err := json.Unmarshal(aiInfo.Tags, &tagNames); err != nil {
			logger.Warn("解析文件标签JSON失败: %v\n", err)
		}
	}

	if len(tagNames) == 0 {
		var tags []models.GlobalTag
		database.DB.Model(&models.GlobalTag{}).
			Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
			Where("file_global_tag_relation.file_id = ?", fileID).
			Find(&tags)

		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}

	response := map[string]interface{}{
		"description":     aiInfo.Description,
		"tags":            tagNames,
		"dominant_color":  aiInfo.DominantColor,
		"resolution":      aiInfo.Resolution,
		"is_nsfw":         aiInfo.IsNSFW,
		"nsfw_score":      aiInfo.NSFWScore,
		"nsfw_evaluation": aiInfo.NSFWEvaluation,
	}

	return response, nil
}

func CheckUserStorageAvailable(userID uint, fileSize int64) (bool, error) {
	var stats models.UserUsageStats
	err := database.DB.Where("user_id = ?", userID).First(&stats).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("查询用户存储统计失败: %v", err)
	}

	var settings models.UserSettings
	err = database.DB.Where("user_id = ?", userID).First(&settings).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, fmt.Errorf("查询用户设置失败: %v", err)
	}

	if err == gorm.ErrRecordNotFound {
		settings = models.UserSettings{
			UserID:             userID,
			StorageLimit:       models.DefaultStorageLimit,
			BandwidthLimit:     models.DefaultBandwidthLimit,
			DefaultAccessLevel: "private",
		}
	}

	totalSizeAfterUpload := stats.TotalSize + fileSize

	if totalSizeAfterUpload > settings.StorageLimit {
		return false, nil
	}

	return true, nil
}
