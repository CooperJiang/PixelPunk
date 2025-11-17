package stats

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/sysinfo"
	"time"
)

func GetDashboardUploadTrends(days int) (*models.DashboardUploadTrendsResponse, error) {
	if days <= 0 {
		days = 7 // 默认7天
	}

	cacheKey := fmt.Sprintf("dashboard:upload_trends:%d", days)
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardUploadTrendsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardUploadTrendsResponse
	result.Timeline = make([]string, days)
	result.UploadCounts = make([]int64, days)
	result.StorageGrowth = make([]int64, days)

	now := time.Now()

	for i := 0; i < days; i++ {
		date := now.AddDate(0, 0, -i).Truncate(24 * time.Hour)
		nextDate := date.Add(24 * time.Hour)

		result.Timeline[days-1-i] = date.Format("01-02")

		var uploadCount int64
		if err := database.DB.Model(&models.File{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Count(&uploadCount).Error; err != nil {
			return nil, err
		}
		result.UploadCounts[days-1-i] = uploadCount

		var storageGrowth *int64
		if err := database.DB.Model(&models.File{}).
			Where("created_at >= ? AND created_at < ?", date, nextDate).
			Select("SUM(size)").Row().Scan(&storageGrowth); err != nil {
			return nil, err
		}
		finalStorageGrowth := int64(0)
		if storageGrowth != nil {
			finalStorageGrowth = *storageGrowth
		}
		result.StorageGrowth[days-1-i] = finalStorageGrowth
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

/* GetDashboardAIServices 获取AI服务状态数据（实时，不缓存） */
func GetDashboardAIServices() (*models.DashboardAIServicesResponse, error) {
	var result models.DashboardAIServicesResponse

	var taggingQueueSize int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status = ?", common.AITaggingStatusPending).
		Count(&taggingQueueSize).Error; err != nil {
		return nil, err
	}

	var taggingProcessing int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status = ?", common.AITaggingStatusPending).
		Count(&taggingProcessing).Error; err != nil {
		return nil, err
	}

	today := time.Now().Truncate(24 * time.Hour)
	var taggingCompletedToday int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status = ? AND updated_at >= ?", common.AITaggingStatusDone, today).
		Count(&taggingCompletedToday).Error; err != nil {
		return nil, err
	}

	result.Tagging = models.AITaggingStats{
		QueueSize:      taggingQueueSize,
		Processing:     taggingProcessing,
		CompletedToday: taggingCompletedToday,
	}

	var vectorQueueSize int64
	var vectorProcessing int64
	var vectorCompletionRate float64

	if database.DB.Migrator().HasTable(&models.VectorProcessingLog{}) {
		var totalVectorTasks int64
		var completedVectorTasks int64

		if err := database.DB.Model(&models.VectorProcessingLog{}).
			Where("created_at >= ?", today).
			Count(&totalVectorTasks).Error; err == nil {

			if err := database.DB.Model(&models.VectorProcessingLog{}).
				Where("created_at >= ? AND action = ?", today, "success").
				Count(&completedVectorTasks).Error; err == nil && totalVectorTasks > 0 {
				vectorCompletionRate = float64(completedVectorTasks) / float64(totalVectorTasks) * 100
			}
		}
	}

	result.Vectors = models.AIVectorStats{
		QueueSize:      vectorQueueSize,
		Processing:     vectorProcessing,
		CompletionRate: vectorCompletionRate,
	}

	var contentReviewPending int64
	var contentReviewApprovedToday int64
	var contentReviewRejectedToday int64

	if database.DB.Migrator().HasTable("review_log") {
	}

	result.ContentReview = models.AIContentReviewStats{
		Pending:       contentReviewPending,
		ApprovedToday: contentReviewApprovedToday,
		RejectedToday: contentReviewRejectedToday,
	}

	taggingHealth := "healthy"
	if taggingQueueSize > 1000 {
		taggingHealth = "degraded"
	}

	vectorHealth := "healthy"
	if vectorCompletionRate < 80 && vectorQueueSize > 0 {
		vectorHealth = "degraded"
	}

	result.HealthStatus = models.AIHealthStatus{
		TaggingService: taggingHealth,
		VectorService:  vectorHealth,
	}

	return &result, nil
}

func GetDashboardShareStats() (*models.DashboardShareStatsResponse, error) {
	cacheKey := "dashboard:share_stats"
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardShareStatsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardShareStatsResponse

	var totalShares int64
	if err := database.DB.Model(&models.Share{}).Count(&totalShares).Error; err != nil {
		return nil, err
	}

	var activeShares int64
	if err := database.DB.Model(&models.Share{}).
		Where("status = ?", 1).
		Count(&activeShares).Error; err != nil {
		return nil, err
	}

	var inactiveShares int64
	if err := database.DB.Model(&models.Share{}).
		Where("status != ?", 1).
		Count(&inactiveShares).Error; err != nil {
		return nil, err
	}

	var totalVisits *int64
	if err := database.DB.Model(&models.Share{}).
		Select("SUM(current_views)").Row().Scan(&totalVisits); err != nil {
		return nil, err
	}
	finalTotalVisits := int64(0)
	if totalVisits != nil {
		finalTotalVisits = *totalVisits
	}

	var totalDownloads int64 = 0
	if database.DB.Migrator().HasTable(&models.FileDownloadLog{}) {
		if err := database.DB.Model(&models.FileDownloadLog{}).
			Where("share_key IS NOT NULL AND share_key != ''").
			Count(&totalDownloads).Error; err != nil {
			totalDownloads = 0
		}
	}

	lastMonth := time.Now().AddDate(0, -1, 0).Truncate(24 * time.Hour)
	var lastMonthShares int64
	if err := database.DB.Model(&models.Share{}).
		Where("created_at < ?", lastMonth).
		Count(&lastMonthShares).Error; err != nil {
		return nil, err
	}

	growthRate := calculateGrowthRate(lastMonthShares, totalShares)

	today := time.Now().Truncate(24 * time.Hour)
	var newSharesToday int64
	if err := database.DB.Model(&models.Share{}).
		Where("created_at >= ?", today).
		Count(&newSharesToday).Error; err != nil {
		return nil, err
	}

	result = models.DashboardShareStatsResponse{
		TotalShares:    totalShares,
		PublicShares:   activeShares,   // 用正常状态的分享代替公开分享
		PrivateShares:  inactiveShares, // 用非正常状态的分享代替私密分享
		TotalVisits:    finalTotalVisits,
		TotalDownloads: totalDownloads,
		GrowthRate:     growthRate,
		NewSharesToday: newSharesToday,
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

/* GetDashboardTagStats 获取标签统计数据 - 重新设计更合理的统计维度 */
func GetDashboardTagStats() (*models.DashboardTagStatsResponse, error) {
	cacheKey := "dashboard:tag_stats"
	if cached, err := cache.Get(cacheKey); err == nil && cached != "" {
		var result models.DashboardTagStatsResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &result, nil
		}
	}

	var result models.DashboardTagStatsResponse

	var totalTagsCount int64
	if err := database.DB.Model(&models.GlobalTag{}).Count(&totalTagsCount).Error; err != nil {
		return nil, err
	}

	var taggedImagesCount int64
	if err := database.DB.Model(&models.FileGlobalTagRelation{}).
		Distinct("file_id").
		Count(&taggedImagesCount).Error; err != nil {
		return nil, err
	}

	var untaggedImages int64
	if err := database.DB.Model(&models.File{}).
		Where("ai_tagging_status IN (?)", []string{
			common.AITaggingStatusNone,
			common.AITaggingStatusFailed,
			common.AITaggingStatusSkipped,
		}).Count(&untaggedImages).Error; err != nil {
		return nil, err
	}

	var popularTags []struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Count int64  `json:"count"`
	}
	if err := database.DB.Table("global_tag as t").
		Select("t.id, t.name, COUNT(r.id) as count").
		Joins("LEFT JOIN file_global_tag_relation as r ON t.id = r.tag_id").
		Group("t.id, t.name").
		Order("count DESC").
		Limit(12).
		Scan(&popularTags).Error; err != nil {
		return nil, err
	}

	popularTagItems := make([]models.DashboardTagItem, len(popularTags))
	for i, tag := range popularTags {
		popularTagItems[i] = models.DashboardTagItem{
			Name:  tag.Name,
			Count: int64(tag.Count),
		}
	}

	result = models.DashboardTagStatsResponse{
		AITagsCount:     totalTagsCount,    // AI标签种类数量
		TaggedImages:    taggedImagesCount, // 已标记的文件数量
		UntaggedImages:  untaggedImages,    // 未标记文件数量
		PopularTags:     popularTagItems,   // 热门标签列表
		ManualTagsCount: 0,                 // 手动标签数量为0（保持兼容性）
	}

	if data, err := json.Marshal(result); err == nil {
		cache.Set(cacheKey, string(data), 5*time.Second)
	}

	return &result, nil
}

/* GetDashboardSystemInfo 获取系统信息（实时，不缓存） */
func GetDashboardSystemInfo() (*models.DashboardSystemInfoResponse, error) {
	sysInfo := sysinfo.GetSystemInfo()

	result := models.DashboardSystemInfoResponse{
		Version: sysInfo.Info.Version, // 假设版本信息在这里，如果没有可以从其他地方获取
		Uptime:  sysInfo.Uptime.Text,
		Status:  "online", // 简单的在线状态
	}

	return &result, nil
}
