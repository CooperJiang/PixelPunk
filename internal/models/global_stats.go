package models

import (
	"pixelpunk/pkg/common"
	"time"
)

/* GlobalStats 全局统计数据表 */
type GlobalStats struct {
	ID   uint      `gorm:"primaryKey" json:"id"`
	Date time.Time `gorm:"type:date"` // 使用日期作为唯一索引

	TotalImages    int64 `gorm:"not null;default:0" json:"total_files"`     // 总文件数
	TotalStorage   int64 `gorm:"not null;default:0" json:"total_storage"`   // 总存储空间(字节)
	TotalBandwidth int64 `gorm:"not null;default:0" json:"total_bandwidth"` // 总带宽使用量(字节)
	TotalFolders   int64 `gorm:"not null;default:0" json:"total_folders"`   // 总文件夹数
	TotalUsers     int64 `gorm:"not null;default:0" json:"total_users"`     // 总用户数
	TotalViews     int64 `gorm:"not null;default:0" json:"total_views"`     // 总访问量

	NewImages    int64 `gorm:"not null;default:0" json:"new_files"`     // 新增文件数
	NewStorage   int64 `gorm:"not null;default:0" json:"new_storage"`   // 新增存储空间(字节)
	NewBandwidth int64 `gorm:"not null;default:0" json:"new_bandwidth"` // 新增带宽使用量(字节)
	NewFolders   int64 `gorm:"not null;default:0" json:"new_folders"`   // 新增文件夹数
	NewUsers     int64 `gorm:"not null;default:0" json:"new_users"`     // 新增用户数
	NewViews     int64 `gorm:"not null;default:0" json:"new_views"`     // 新增访问量

	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
}

func (GlobalStats) TableName() string {
	return "global_stats"
}
