package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* FileStats 文件统计信息模型 */
type FileStats struct {
	ID             uint             `gorm:"primarykey" json:"id"`
	FileID         string           `gorm:"size:32;not null;uniqueIndex:idx_file_stats_file_id" json:"file_id"`
	Views          int64            `gorm:"not null;default:0" json:"views"`
	Downloads      int64            `gorm:"not null;default:0" json:"downloads"` // 新增下载次数统计
	Bandwidth      int64            `gorm:"not null;default:0" json:"bandwidth"`
	LastViewAt     *common.JSONTime `json:"last_view_at"`
	LastDownloadAt *common.JSONTime `json:"last_download_at"` // 新增最后下载时间
	UpdatedAt      common.JSONTime  `json:"updated_at"`
}

func (FileStats) TableName() string {
	return "file_stats"
}

func (s *FileStats) BeforeCreate(tx *gorm.DB) error {
	if s.Views == 0 {
		s.Views = 0
	}
	if s.Downloads == 0 {
		s.Downloads = 0
	}
	if s.Bandwidth == 0 {
		s.Bandwidth = 0
	}
	return nil
}

func (s *FileStats) IncrementViews(tx *gorm.DB) error {
	now := common.JSONTimeNow()
	s.Views++
	s.LastViewAt = &now
	return tx.Save(s).Error
}

func (s *FileStats) IncrementDownloads(tx *gorm.DB) error {
	now := common.JSONTimeNow()
	s.Downloads++
	s.LastDownloadAt = &now
	return tx.Save(s).Error
}

func (s *FileStats) AddBandwidth(tx *gorm.DB, size int64) error {
	s.Bandwidth += size
	return tx.Save(s).Error
}
