package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* GuestUploadLog 游客上传日志模型 */
type GuestUploadLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	GuestUploadLimitID uint   `gorm:"not null;index" json:"guest_upload_limit_id"` // 关联的游客限制记录ID
	FileID             string `gorm:"size:32;index" json:"file_id"`                // 关联的文件ID

	Fingerprint     string     `gorm:"size:64;not null;index" json:"fingerprint"` // 浏览器指纹
	IP              string     `gorm:"size:45;not null;index" json:"ip"`          // IP地址
	UserAgent       string     `gorm:"size:500" json:"user_agent"`                // 用户代理信息
	StorageDuration string     `gorm:"size:20" json:"storage_duration"`           // 选择的存储时长
	ExpiresAt       *time.Time `gorm:"index" json:"expires_at"`                   // 预计过期时间

	Status       string `gorm:"size:20;not null;default:'success'" json:"status"` // 状态：success/failed/blocked
	Reason       string `gorm:"size:255" json:"reason"`                           // 失败或封禁原因
	FileSize     int64  `json:"file_size"`                                        // 文件大小
	FileName     string `gorm:"size:255" json:"file_name"`                        // 文件名
	OriginalName string `gorm:"size:255" json:"original_name"`                    // 原始文件名

	GuestUploadLimit *GuestUploadLimit `gorm:"foreignKey:GuestUploadLimitID;references:ID" json:"guest_upload_limit"`
	File             *File             `gorm:"foreignKey:FileID;references:ID" json:"file"`
}

func (GuestUploadLog) TableName() string {
	return "guest_upload_log"
}

func (l *GuestUploadLog) BeforeCreate(tx *gorm.DB) error {
	if l.Status == "" {
		l.Status = "success"
	}
	return nil
}

/* IsSuccess 检查是否上传成功 */
func (l *GuestUploadLog) IsSuccess() bool {
	return l.Status == "success"
}

/* IsFailed 检查是否上传失败 */
func (l *GuestUploadLog) IsFailed() bool {
	return l.Status == "failed"
}

/* IsBlocked 检查是否被阻止 */
func (l *GuestUploadLog) IsBlocked() bool {
	return l.Status == "blocked"
}

/* MarkAsSuccess 标记为成功 */
func (l *GuestUploadLog) MarkAsSuccess() {
	l.Status = "success"
	l.Reason = ""
}

/* MarkAsFailed 标记为失败 */
func (l *GuestUploadLog) MarkAsFailed(reason string) {
	l.Status = "failed"
	l.Reason = reason
}

/* MarkAsBlocked 标记为被阻止 */
func (l *GuestUploadLog) MarkAsBlocked(reason string) {
	l.Status = "blocked"
	l.Reason = reason
}
