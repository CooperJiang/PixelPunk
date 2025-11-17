package stats

import (
	"encoding/json"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/sysinfo"
	"pixelpunk/pkg/utils"
	"time"
)

func GetDashboardUserStats() (*models.DashboardUserStatsResponse, error) {
	cacheKey := "dashboard:user_stats"
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardUserStatsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardUserStatsResponse

	var totalUsers int64
	if err := database.DB.Model(&models.User{}).Count(&totalUsers).Error; err != nil {
		return nil, err
	}

	var bannedUsers int64
	if err := database.DB.Model(&models.User{}).Where("status = ?", common.UserStatusDisabled).Count(&bannedUsers).Error; err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	var newUsersToday int64
	if err := database.DB.Model(&models.User{}).Where("created_at >= ?", today).Count(&newUsersToday).Error; err != nil {
		return nil, err
	}

	var activeUsersToday int64
	if err := database.DB.Model(&models.File{}).
		Where("created_at >= ?", today).
		Distinct("user_id").
		Count(&activeUsersToday).Error; err != nil {
		return nil, err
	}

	lastMonth := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	var lastMonthUsers int64
	if err := database.DB.Model(&models.User{}).Where("created_at < ?", lastMonth).Count(&lastMonthUsers).Error; err != nil {
		return nil, err
	}

	growthRate := calculateGrowthRate(lastMonthUsers, totalUsers)

	result = models.DashboardUserStatsResponse{
		TotalUsers:       totalUsers,
		ActiveUsersToday: activeUsersToday,
		NewUsersToday:    newUsersToday,
		GrowthRate:       growthRate,
		BannedUsers:      bannedUsers,
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

func GetDashboardFileStats() (*models.DashboardFileStatsResponse, error) {
	cacheKey := "dashboard:file_stats"
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardFileStatsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardFileStatsResponse

	var totalImages int64
	if err := database.DB.Model(&models.File{}).Count(&totalImages).Error; err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	var uploadedToday int64
	if err := database.DB.Model(&models.File{}).Where("created_at >= ?", today).Count(&uploadedToday).Error; err != nil {
		return nil, err
	}

	var aiTaggedCount int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status = ?", common.AITaggingStatusDone).
		Count(&aiTaggedCount).Error; err != nil {
		return nil, err
	}

	aiTaggedPercentage := float64(0)
	if totalImages > 0 {
		aiTaggedPercentage = float64(aiTaggedCount) / float64(totalImages) * 100
	}

	var pendingReview int64
	if err := database.DB.Model(&models.File{}).
		Where("status = ?", "pending_review").
		Count(&pendingReview).Error; err != nil {
		return nil, err
	}

	var nsfwDetected int64
	if err := database.DB.Model(&models.FileAIInfo{}).Where("is_nsfw = ?", true).Count(&nsfwDetected).Error; err != nil {
		return nil, err
	}

	var untaggedCount int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status IN (?)", []string{
			common.AITaggingStatusNone,
			common.AITaggingStatusFailed,
			common.AITaggingStatusSkipped,
		}).Count(&untaggedCount).Error; err != nil {
		return nil, err
	}

	result = models.DashboardFileStatsResponse{
		TotalImages:        totalImages,
		UploadedToday:      uploadedToday,
		AITaggedPercentage: aiTaggedPercentage,
		PendingReview:      pendingReview,
		NSFWDetected:       nsfwDetected,
		UntaggedCount:      untaggedCount,
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

func GetDashboardStorageStats() (*models.DashboardStorageStatsResponse, error) {
	cacheKey := "dashboard:storage_stats"
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardStorageStatsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardStorageStatsResponse

	var totalStorage *int64
	if err := database.DB.Model(&models.File{}).Select("SUM(size)").Row().Scan(&totalStorage); err != nil {
		return nil, err
	}
	finalTotalStorage := int64(0)
	if totalStorage != nil {
		finalTotalStorage = *totalStorage
	}

	var imageCount int64
	if err := database.DB.Model(&models.File{}).Count(&imageCount).Error; err != nil {
		return nil, err
	}

	averageFileSize := int64(0)
	if imageCount > 0 {
		averageFileSize = finalTotalStorage / imageCount
	}

	now := time.Now()
	currentMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	previousMonthStart := currentMonthStart.AddDate(0, -1, 0)

	var addedCurrentMonth int64
	if err := database.DB.Model(&models.File{}).
		Where("created_at >= ?", currentMonthStart).
		Select("COALESCE(SUM(size), 0)").Row().Scan(&addedCurrentMonth); err != nil {
		return nil, err
	}

	var addedPreviousMonth int64
	if err := database.DB.Model(&models.File{}).
		Where("created_at >= ? AND created_at < ?", previousMonthStart, currentMonthStart).
		Select("COALESCE(SUM(size), 0)").Row().Scan(&addedPreviousMonth); err != nil {
		return nil, err
	}

	growthRate := calculateGrowthRate(addedPreviousMonth, addedCurrentMonth)

	today := time.Now().Truncate(24 * time.Hour)
	var newStorageToday *int64
	if err := database.DB.Model(&models.File{}).
		Where("created_at >= ?", today).
		Select("SUM(size)").Row().Scan(&newStorageToday); err != nil {
		return nil, err
	}

	finalNewStorageToday := int64(0)
	if newStorageToday != nil {
		finalNewStorageToday = *newStorageToday
	}

	usagePercentage := float64(0)

	var totalBandwidth int64
	if err := database.DB.Model(&models.GlobalStats{}).
		Select("COALESCE(MAX(total_bandwidth), 0)").Row().Scan(&totalBandwidth); err != nil {
		totalBandwidth = 0
	}

	result = models.DashboardStorageStatsResponse{
		TotalStorage:        finalTotalStorage,
		FormattedStorage:    utils.FormatBytes(finalTotalStorage),
		GrowthRate:          growthRate,
		AverageFileSize:     averageFileSize,
		UsagePercentage:     usagePercentage,
		NewStorageToday:     finalNewStorageToday,
		FormattedNewStorage: utils.FormatBytes(finalNewStorageToday),
		TotalBandwidth:      totalBandwidth,
		FormattedBandwidth:  utils.FormatBytes(totalBandwidth),
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

func GetDashboardSystemResources() (*models.DashboardSystemResourcesResponse, error) {
	sysInfo := sysinfo.GetSystemInfo()

	result := models.DashboardSystemResourcesResponse{
		CPUUsage:    sysInfo.CPU.UsagePercent,
		MemoryUsage: sysInfo.Memory.UsagePercent,
		DiskUsage:   sysInfo.Disk.UsagePercent,
		LoadAverage: []float64{sysInfo.Load.OneMin, sysInfo.Load.FiveMin, sysInfo.Load.FifteenMin},
		Uptime:      sysInfo.Uptime.Text,
	}

	return &result, nil
}
