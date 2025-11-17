package models

/* DashboardUserStatsResponse 用户统计响应 */
type DashboardUserStatsResponse struct {
	TotalUsers       int64   `json:"total_users"`
	ActiveUsersToday int64   `json:"active_users_today"`
	NewUsersToday    int64   `json:"new_users_today"`
	GrowthRate       float64 `json:"growth_rate"`
	BannedUsers      int64   `json:"banned_users"`
}

/* DashboardFileStatsResponse 文件统计响应 */
type DashboardFileStatsResponse struct {
	TotalImages        int64   `json:"total_files"`
	UploadedToday      int64   `json:"uploaded_today"`
	AITaggedPercentage float64 `json:"ai_tagged_percentage"`
	PendingReview      int64   `json:"pending_review"`
	NSFWDetected       int64   `json:"nsfw_detected"`
	UntaggedCount      int64   `json:"untagged_count"`
}

/* DashboardStorageStatsResponse 存储统计响应 */
type DashboardStorageStatsResponse struct {
	TotalStorage        int64   `json:"total_storage"`
	FormattedStorage    string  `json:"formatted_storage"`
	GrowthRate          float64 `json:"growth_rate"`
	AverageFileSize     int64   `json:"average_file_size"`
	UsagePercentage     float64 `json:"usage_percentage"`
	NewStorageToday     int64   `json:"new_storage_today"`     // 今日新增存储量(字节)
	FormattedNewStorage string  `json:"formatted_new_storage"` // 今日新增存储量(格式化)
	TotalBandwidth      int64   `json:"total_bandwidth"`       // 总带宽使用(字节)
	FormattedBandwidth  string  `json:"formatted_bandwidth"`   // 格式化带宽
}

/* DashboardSystemResourcesResponse 系统资源响应 */
type DashboardSystemResourcesResponse struct {
	CPUUsage    float64   `json:"cpu_usage"`
	MemoryUsage float64   `json:"memory_usage"`
	DiskUsage   float64   `json:"disk_usage"`
	LoadAverage []float64 `json:"load_average"`
	Uptime      string    `json:"uptime"`
}

/* DashboardUploadTrendsResponse 上传趋势响应 */
type DashboardUploadTrendsResponse struct {
	Timeline      []string `json:"timeline"`
	UploadCounts  []int64  `json:"upload_counts"`
	StorageGrowth []int64  `json:"storage_growth"`
}

/* DashboardAIServicesResponse AI服务状态响应 */
type DashboardAIServicesResponse struct {
	Tagging       AITaggingStats       `json:"tagging"`
	Vectors       AIVectorStats        `json:"vectors"`
	ContentReview AIContentReviewStats `json:"content_review"`
	HealthStatus  AIHealthStatus       `json:"health_status"`
}

/* AITaggingStats AI标注统计 */
type AITaggingStats struct {
	QueueSize      int64 `json:"queue_size"`
	Processing     int64 `json:"processing"`
	CompletedToday int64 `json:"completed_today"`
}

/* AIVectorStats AI向量统计 */
type AIVectorStats struct {
	QueueSize      int64   `json:"queue_size"`
	Processing     int64   `json:"processing"`
	CompletionRate float64 `json:"completion_rate"`
}

/* AIContentReviewStats AI内容审核统计 */
type AIContentReviewStats struct {
	Pending       int64 `json:"pending"`
	ApprovedToday int64 `json:"approved_today"`
	RejectedToday int64 `json:"rejected_today"`
}

/* AIHealthStatus AI服务健康状态 */
type AIHealthStatus struct {
	TaggingService string `json:"tagging_service"` // healthy, degraded, down
	VectorService  string `json:"vector_service"`  // healthy, degraded, down
}

/* DashboardShareStatsResponse 分享统计响应 */
type DashboardShareStatsResponse struct {
	TotalShares    int64   `json:"total_shares"`
	PublicShares   int64   `json:"public_shares"`
	PrivateShares  int64   `json:"private_shares"`
	TotalVisits    int64   `json:"total_visits"`
	TotalDownloads int64   `json:"total_downloads"`
	GrowthRate     float64 `json:"growth_rate"`
	NewSharesToday int64   `json:"new_shares_today"` // 今日新增分享数
}

/* DashboardTagStatsResponse 标签统计响应 */
type DashboardTagStatsResponse struct {
	AITagsCount     int64              `json:"ai_tags_count"`     // AI标签种类数量
	TaggedImages    int64              `json:"tagged_files"`      // 已标记的文件数量
	UntaggedImages  int64              `json:"untagged_files"`    // 未标记文件数量
	PopularTags     []DashboardTagItem `json:"popular_tags"`      // 热门标签列表
	ManualTagsCount int64              `json:"manual_tags_count"` // 手动标签数量（保持兼容性，设为0）
}

/* DashboardTagItem 标签项目 */
type DashboardTagItem struct {
	Name  string `json:"name"`
	Count int64  `json:"count"`
}

/* DashboardSystemInfoResponse 系统信息响应 */
type DashboardSystemInfoResponse struct {
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
	Status  string `json:"status"`
}
