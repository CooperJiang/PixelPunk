package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

type FileTaggingLog struct {
	ID         uint            `gorm:"primarykey" json:"id"`
	FileID     string          `gorm:"size:32;index" json:"file_id"`
	FileType   string          `gorm:"size:20;index" json:"file_type"`
	Status     string          `gorm:"size:20;not null" json:"status"`
	Action     string          `gorm:"size:20;not null" json:"action"`
	Type       string          `gorm:"size:100;not null" json:"type"`
	Data       string          `gorm:"type:text;not null" json:"data"`
	Message    string          `gorm:"type:text" json:"message"`
	OperatorID uint            `gorm:"index" json:"operator_id"`
	Duration   int64           `gorm:"default:0" json:"duration"`
	CreatedAt  common.JSONTime `json:"created_at"`
	UpdatedAt  common.JSONTime `json:"updated_at"`
}

func (FileTaggingLog) TableName() string {
	return "file_tagging_log"
}

func (log *FileTaggingLog) BeforeCreate(tx *gorm.DB) error {
	if log.Action == "" {
		log.Action = common.TaggingActionAuto
	}
	if log.Data == "" {
		log.Data = "{}"
	}
	return nil
}
