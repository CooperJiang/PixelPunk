package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* FileDownloadLog 文件下载历史记录 */
type FileDownloadLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `gorm:"index" json:"created_at"`

	UserID    uint   `gorm:"index" json:"user_id"`                  // 下载用户ID (分享下载时为0，手动管理无外键约束)
	FileID    string `gorm:"size:32;not null;index" json:"file_id"` // 文件ID
	FileSize  int64  `gorm:"not null" json:"file_size"`             // 下载的文件大小
	FileType  string `gorm:"size:20;index" json:"file_type"`        // 文件类型
	IPAddress string `gorm:"size:45;index" json:"ip_address"`       // 下载时的IP地址
	UserAgent string `gorm:"size:500" json:"user_agent"`            // 用户代理信息
	ShareKey  string `gorm:"size:32;index" json:"share_key"`        // 分享密钥（分享下载时使用）

	DownloadSource string `gorm:"size:20;default:direct" json:"download_source"` // direct/share/api

}

func (FileDownloadLog) TableName() string {
	return "file_download_log"
}

func (log *FileDownloadLog) BeforeCreate(tx *gorm.DB) error {
	if log.DownloadSource == "" {
		log.DownloadSource = "direct"
	}
	return nil
}
