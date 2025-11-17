package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

type VectorProcessingLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	FileID   string `gorm:"size:50;not null;index:idx_file_id" json:"file_id"`
	Action   string `gorm:"size:20;not null;index:idx_action" json:"action"`
	Type     string `gorm:"size:100;not null" json:"type"`
	Data     string `gorm:"type:text;not null" json:"data"`
	Message  string `gorm:"type:text" json:"message"`
	Model    string `gorm:"size:50" json:"model"`
	Duration int    `gorm:"default:0" json:"duration"`

	ErrorCode string `gorm:"size:500" json:"error_code,omitempty"`
	TaskID    string `gorm:"size:100" json:"task_id,omitempty"`
}

func (VectorProcessingLog) TableName() string {
	return "vector_processing_log"
}

func (vpl *VectorProcessingLog) BeforeCreate(tx *gorm.DB) error {
	now := common.JSONTime(time.Now())
	if time.Time(vpl.CreatedAt).IsZero() {
		vpl.CreatedAt = now
	}
	if time.Time(vpl.UpdatedAt).IsZero() {
		vpl.UpdatedAt = now
	}
	if vpl.Data == "" {
		vpl.Data = "{}"
	}
	return nil
}

const (
	VectorLogActionStart   = "start"
	VectorLogActionSuccess = "success"
	VectorLogActionError   = "error"
	VectorLogActionRetry   = "retry"
	VectorLogActionReset   = "reset"
	VectorLogActionSkip    = "skip"
)
