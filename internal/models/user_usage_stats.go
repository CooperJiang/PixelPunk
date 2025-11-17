package models

import (
	"pixelpunk/pkg/common"
)

/* UserUsageStats 用户使用统计模型 */
type UserUsageStats struct {
	ID             uint            `gorm:"primarykey" json:"id"`
	UserID         uint            `gorm:"not null"`
	TotalImages    int             `gorm:"not null;default:0" json:"total_files"`
	TotalSize      int64           `gorm:"not null;default:0" json:"total_size"`
	TotalBandwidth int64           `gorm:"not null;default:0" json:"total_bandwidth"`
	TotalViews     int64           `gorm:"not null;default:0" json:"total_views"`
	UpdatedAt      common.JSONTime `json:"updated_at"`
}

func (UserUsageStats) TableName() string {
	return "user_usage_stats"
}

func (s *UserUsageStats) StorageUsedPercent(storageLimit int64) float64 {
	if storageLimit <= 0 {
		return 0
	}
	return float64(s.TotalSize) / float64(storageLimit) * 100
}
