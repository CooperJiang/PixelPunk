package models

import (
	"time"
)

/* VectorJob 队列表（向量化队列） */
type VectorJob struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	FileID     string     `gorm:"size:32;uniqueIndex:idx_vector_job_file" json:"file_id"`
	Status     string     `gorm:"size:20;index" json:"status"` // queued|processing|done|failed
	Attempt    int        `gorm:"default:0" json:"attempt"`
	Priority   int        `gorm:"default:0;index" json:"priority"`
	LeaseUntil *time.Time `gorm:"index" json:"lease_until"`
	LeaseBy    string     `gorm:"size:64;index" json:"lease_by"`
	LastError  string     `gorm:"type:text" json:"last_error"`
}

func (VectorJob) TableName() string { return "vector_job" }
