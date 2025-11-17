package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

type UploadSession struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `gorm:"index:idx_upload_session_created_at" json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	SessionID string `gorm:"size:32;uniqueIndex;not null" json:"session_id"`
	UserID    uint   `gorm:"not null;index:idx_upload_session_user_id" json:"user_id"`

	FileName string `gorm:"size:255;not null" json:"file_name"`
	FileSize int64  `gorm:"not null" json:"file_size"`
	FileMD5  string `gorm:"size:32" json:"file_md5"`
	MimeType string `gorm:"size:50" json:"mime_type"`

	ChunkSize   int64 `gorm:"not null" json:"chunk_size"`
	TotalChunks int   `gorm:"not null" json:"total_chunks"`

	Status   string `gorm:"size:20;not null;default:pending" json:"status"` // pending/uploading/completed/failed/expired
	Progress int    `gorm:"default:0" json:"progress"`

	FolderID    string `gorm:"size:32" json:"folder_id"`
	AccessLevel string `gorm:"size:20;default:private" json:"access_level"`
	Optimize    bool   `gorm:"default:true" json:"optimize"`

	WatermarkConfig string `gorm:"type:text" json:"watermark_config"`

	FileID string `gorm:"size:32" json:"file_id"`

	ExpiresAt time.Time `gorm:"index:idx_upload_session_expires_at" json:"expires_at"` // 24小时后过期
}

func (UploadSession) TableName() string {
	return "upload_session"
}

func (us *UploadSession) BeforeCreate(tx *gorm.DB) error {
	if us.Status == "" {
		us.Status = "pending"
	}

	if us.AccessLevel == "" {
		us.AccessLevel = "private"
	}

	if us.ExpiresAt.IsZero() {
		us.ExpiresAt = time.Now().Add(24 * time.Hour)
	}

	return nil
}

func (us *UploadSession) IsExpired() bool {
	return time.Now().After(us.ExpiresAt)
}

func (us *UploadSession) IsCompleted() bool {
	return us.Status == "completed"
}

func (us *UploadSession) IsFailed() bool {
	return us.Status == "failed"
}

func (us *UploadSession) IsActive() bool {
	return us.Status == "pending" || us.Status == "uploading"
}
