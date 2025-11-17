package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* UserSettings 用户设置模型 */
type UserSettings struct {
	ID                 uint            `gorm:"primarykey" json:"id"`
	UserID             uint            `gorm:"not null"`
	StorageLimit       int64           `gorm:"not null;default:5368709120" json:"storage_limit"`     // 默认500M
	BandwidthLimit     int64           `gorm:"not null;default:107374182400" json:"bandwidth_limit"` // 默认1GB
	DefaultAccessLevel string          `gorm:"size:20;not null;default:private" json:"default_access_level"`
	OptimizeImages     bool            `gorm:"not null;default:false" json:"optimize_files"`
	CreatedAt          common.JSONTime `json:"created_at"`
	UpdatedAt          common.JSONTime `json:"updated_at"`
}

func (UserSettings) TableName() string {
	return "user_settings"
}

/* DefaultStorageLimit 默认存储空间限制 500MB */
const DefaultStorageLimit int64 = 500 * 1024 * 1024

/* DefaultBandwidthLimit 默认带宽限制 1GB */
const DefaultBandwidthLimit int64 = 1 * 1024 * 1024 * 1024

func (s *UserSettings) BeforeCreate(tx *gorm.DB) error {
	if s.StorageLimit == 0 {
		s.StorageLimit = DefaultStorageLimit
	}
	if s.BandwidthLimit == 0 {
		s.BandwidthLimit = DefaultBandwidthLimit
	}
	if s.DefaultAccessLevel == "" {
		s.DefaultAccessLevel = "private"
	}
	return nil
}
