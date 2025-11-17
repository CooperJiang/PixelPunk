package models

import (
	"pixelpunk/pkg/common"
	"time"
)

type UserStatsResponse struct {
	Storage   StorageStats   `json:"storage"`
	Bandwidth BandwidthStats `json:"bandwidth"`
	Files     FilesStats     `json:"files"`
	Settings  UserSettings   `json:"settings"`
}

type StorageStats struct {
	Used       int64   `json:"used"`
	Limit      int64   `json:"limit"`
	Percentage float64 `json:"percentage"`
}

type BandwidthStats struct {
	Used       int64   `json:"used"`
	Limit      int64   `json:"limit"`
	Percentage float64 `json:"percentage"`
}

type FilesStats struct {
	Total int64 `json:"total"`
	Views int64 `json:"views"`
}

type SystemStatsResponse struct {
	Users     UsersStats     `json:"users"`
	Files     FilesGrowth    `json:"files"`
	Storage   StorageGrowth  `json:"storage"`
	Bandwidth BandwidthTotal `json:"bandwidth"`
	System    SystemInfo     `json:"system"`
}

type UsersStats struct {
	Total       int64   `json:"total"`
	GrowthRate  float64 `json:"growth_rate"`
	LastMonth   int64   `json:"last_month"`
	BannedUsers int64   `json:"banned_users"`
}

type FilesGrowth struct {
	Total       int64   `json:"total"`
	GrowthRate  float64 `json:"growth_rate"`
	LastMonth   int64   `json:"last_month"`
	Recommended int64   `json:"recommended"`
	Untagged    int64   `json:"untagged"`
	NSFW        int64   `json:"nsfw"`
}

type StorageGrowth struct {
	Total          int64   `json:"total"`
	GrowthRate     float64 `json:"growth_rate"`
	LastMonth      int64   `json:"last_month"`
	FormattedTotal string  `json:"formatted_total"`
}

type BandwidthTotal struct {
	Total          int64  `json:"total"`
	FormattedTotal string `json:"formatted_total"`
}

type SystemInfo struct {
	Load   SystemLoad   `json:"load"`
	Memory SystemMemory `json:"memory"`
	CPU    SystemCPU    `json:"cpu"`
	Disk   SystemDisk   `json:"disk"`
	Uptime SystemUptime `json:"uptime"`
	Info   BasicSysInfo `json:"info"`
}

type SystemLoad struct {
	OneMin     float64 `json:"1min"`
	FiveMin    float64 `json:"5min"`
	FifteenMin float64 `json:"15min"`
}

type SystemMemory struct {
	Total          uint64  `json:"total"`
	Used           uint64  `json:"used"`
	Free           uint64  `json:"free"`
	FormattedTotal string  `json:"formatted_total"`
	FormattedUsed  string  `json:"formatted_used"`
	FormattedFree  string  `json:"formatted_free"`
	UsagePercent   float64 `json:"usage_percent"`
}

type SystemCPU struct {
	Cores        int     `json:"cores"`
	UsagePercent float64 `json:"usage_percent"`
}

type SystemDisk struct {
	Total          uint64  `json:"total"`
	Used           uint64  `json:"used"`
	Free           uint64  `json:"free"`
	FormattedTotal string  `json:"formatted_total"`
	FormattedUsed  string  `json:"formatted_used"`
	FormattedFree  string  `json:"formatted_free"`
	UsagePercent   float64 `json:"usage_percent"`
}

type SystemUptime struct {
	Seconds int64  `json:"seconds"`
	Text    string `json:"text"`
}

type BasicSysInfo struct {
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Hostname string `json:"hostname"`
	Version  string `json:"version"`
}

type MonthlyUserGrowth struct {
	Month            string  `json:"month"`
	MonthDisplay     string  `json:"month_display"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
	TotalUsers       int64   `json:"total_users"`
	NewUsers         int64   `json:"new_users"`
	ActiveUsers      int64   `json:"active_users"`
	GrowthRate       float64 `json:"growth_rate"`
	DailyAvgNewUsers float64 `json:"daily_avg_new_users"`
}

type LatestUserInfo struct {
	ID                    uint            `json:"id"`
	Username              string          `json:"username"`
	Email                 string          `json:"email"`
	Avatar                string          `json:"avatar"`
	AvatarFullURL         string          `json:"avatar_full_url"`
	CreatedAt             common.JSONTime `json:"created_at"`
	DaysSinceRegistration int             `json:"days_since_registration"`
	LastActiveAt          common.JSONTime `json:"last_active_at"`
	LastActiveTimeAgo     string          `json:"last_active_time_ago"`
	Status                string          `json:"status"`
	Role                  string          `json:"role"`
	ImagesCount           int64           `json:"files_count"`
	StorageUsed           int64           `json:"storage_used"`
	StorageFormatted      string          `json:"storage_formatted"`
	ViewCount             int64           `json:"view_count"`
}

type LatestFileInfo struct {
	ID              string           `json:"id"` // 修改为string类型，与Image模型一致
	FileName        string           `json:"file_name"`
	DisplayName     string           `json:"display_name"`
	Size            int64            `json:"size"`
	SizeFormatted   string           `json:"size_formatted"`
	Format          string           `json:"format"` // 文件格式，如 jpg, png, webp
	URL             string           `json:"url"`
	ThumbnailURL    string           `json:"thumbnail_url"`
	FullURL         string           `json:"full_url"`
	FullThumbURL    string           `json:"full_thumb_url"`
	CreatedAt       common.JSONTime  `json:"created_at"`
	User            FileOwner        `json:"user"`
	AIInfo          interface{}      `json:"ai_info,omitempty"` // AI信息
	StorageDuration string           `json:"storage_duration"`  // 存储时长：3d/7d/30d/permanent
	ExpiresAt       *common.JSONTime `json:"expires_at"`        // 过期时间
	IsTimeLimited   bool             `json:"is_time_limited"`   // 是否为限时存储
}

type FileOwner struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type DailyFileUpload struct {
	Date           time.Time `json:"date"`
	DateFormatted  string    `json:"date_formatted"`  // 格式化日期，如 "2023-06-01"
	DayOfWeek      string    `json:"day_of_week"`     // 星期几，如 "星期一"
	Count          int64     `json:"count"`           // 当日上传数量
	CumulativeSize int64     `json:"cumulative_size"` // 当日上传大小合计 (bytes)
	SizeFormatted  string    `json:"size_formatted"`  // 格式化后的大小，如 "52 MB"
	ActiveUsers    int64     `json:"active_users"`    // 当日有上传行为的活跃用户数
}
