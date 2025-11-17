package models

import (
	"time"
)

/* AIJob 打标任务队列表（方案B） */
type AIJob struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FileID     string     `gorm:"size:32;uniqueIndex:idx_ai_job_file" json:"file_id"`
	Status     string     `gorm:"size:20;index" json:"status"` // queued|processing|done|failed|ignored|skipped
	Tries      int        `gorm:"default:0" json:"tries"`
	Priority   int        `gorm:"default:0;index" json:"priority"`
	LeaseUntil *time.Time `gorm:"index" json:"lease_until"`
	LeaseBy    string     `gorm:"size:64;index" json:"lease_by"`
	LastError  string     `gorm:"type:text" json:"last_error"`
}

func (AIJob) TableName() string {
	return "ai_job"
}
