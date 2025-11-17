package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* UserBandwidthUsage 用户带宽使用记录表 (按月统计) */
type UserBandwidthUsage struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	UserID    uint            `gorm:"not null;index:idx_user_month" json:"user_id"`
	Year      int             `gorm:"not null;index:idx_user_month" json:"year"`  // 年份
	Month     int             `gorm:"not null;index:idx_user_month" json:"month"` // 月份
	UsedBytes int64           `gorm:"not null;default:0" json:"used_bytes"`       // 当月已用流量(字节)
	UpdatedAt common.JSONTime `json:"updated_at"`
}

func (UserBandwidthUsage) TableName() string {
	return "user_bandwidth_usage"
}

func (u *UserBandwidthUsage) BeforeCreate(tx *gorm.DB) error {
	if u.UsedBytes == 0 {
		u.UsedBytes = 0
	}
	return nil
}
