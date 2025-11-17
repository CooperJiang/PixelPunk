package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* ReviewLog 审核记录模型 */
type ReviewLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	FileID     string `gorm:"size:32;not null;index:idx_review_log_file_id" json:"file_id"`
	AuditorID  uint   `gorm:"not null;index:idx_review_log_auditor_id" json:"auditor_id"`
	UploaderID uint   `gorm:"not null;index:idx_review_log_uploader_id" json:"uploader_id"`
	Action     string `gorm:"size:20;not null" json:"action"` // approve/reject
	DeleteType string `gorm:"size:20" json:"delete_type"`     // soft/hard (仅reject时使用)
	Reason     string `gorm:"type:text" json:"reason"`        // 审核原因/备注

	NSFWScore     *float64 `json:"nsfw_score"`     // AI检测的NSFW分数
	NSFWThreshold *float64 `json:"nsfw_threshold"` // 当时使用的阈值
	IsNSFW        *bool    `json:"is_nsfw"`        // AI判断结果

	File     *File `gorm:"foreignKey:FileID;references:ID" json:"file,omitempty"`
	Auditor  *User `gorm:"foreignKey:AuditorID;references:ID" json:"auditor,omitempty"`
	Uploader *User `gorm:"foreignKey:UploaderID;references:ID" json:"uploader,omitempty"`
}

func (ReviewLog) TableName() string {
	return "review_log"
}

func (r *ReviewLog) BeforeCreate(tx *gorm.DB) error {
	if r.Action != "approve" && r.Action != "reject" {
		return gorm.ErrInvalidValue
	}

	if r.Action == "reject" && r.DeleteType == "" {
		r.DeleteType = "soft"
	}

	return nil
}
