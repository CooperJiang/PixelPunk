package models

import (
	"encoding/json"
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* VectorVerificationTask 向量验证任务模型 */
type VectorVerificationTask struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	TaskID    string `gorm:"size:32;uniqueIndex;not null" json:"task_id"`
	Status    string `gorm:"type:varchar(20);default:pending" json:"status"`
	TaskType  string `gorm:"type:varchar(20);default:manual" json:"task_type"`
	CreatorID *uint  `gorm:"index" json:"creator_id"`

	TotalCount     int `gorm:"default:0" json:"total_count"`
	ProcessedCount int `gorm:"default:0" json:"processed_count"`
	VerifiedCount  int `gorm:"default:0" json:"verified_count"`
	MissingCount   int `gorm:"default:0" json:"missing_count"`
	ErrorCount     int `gorm:"default:0" json:"error_count"`

	StartedAt         *time.Time `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	EstimatedDuration *int       `json:"estimated_duration"` // 预估耗时(秒)

	BatchSize        int    `gorm:"default:100" json:"batch_size"`
	FilterConditions string `gorm:"type:json" json:"filter_conditions"` // JSON字符串
	ErrorDetails     string `gorm:"type:text" json:"error_details"`

	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Creator *User `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`
}

/* VectorVerificationTaskStatus 任务状态常量 */
const (
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
)

/* VectorVerificationTaskType 任务类型常量 */
const (
	TaskTypeManual    = "manual"
	TaskTypeScheduled = "scheduled"
	TaskTypePartial   = "partial"
)

func (VectorVerificationTask) TableName() string {
	return "vector_verification_task"
}

/* BeforeCreate GORM钩子：创建前 */
func (t *VectorVerificationTask) BeforeCreate(tx *gorm.DB) error {
	if t.TaskID == "" {
		t.TaskID = generateTaskID()
	}
	if t.BatchSize == 0 {
		t.BatchSize = 100
	}
	return nil
}

func (t *VectorVerificationTask) SetFilterConditions(filters map[string]interface{}) error {
	if filters == nil {
		t.FilterConditions = ""
		return nil
	}

	data, err := json.Marshal(filters)
	if err != nil {
		return err
	}
	t.FilterConditions = string(data)
	return nil
}

/* GetFilterConditions 获取筛选条件 */
func (t *VectorVerificationTask) GetFilterConditions() (map[string]interface{}, error) {
	if t.FilterConditions == "" {
		return make(map[string]interface{}), nil
	}

	var filters map[string]interface{}
	err := json.Unmarshal([]byte(t.FilterConditions), &filters)
	return filters, err
}

/* GetProgress 获取任务进度百分比 */
func (t *VectorVerificationTask) GetProgress() float64 {
	if t.TotalCount == 0 {
		return 0.0
	}
	return float64(t.ProcessedCount) / float64(t.TotalCount) * 100
}

/* GetDuration 获取任务已执行时长 */
func (t *VectorVerificationTask) GetDuration() time.Duration {
	if t.StartedAt == nil {
		return 0
	}

	endTime := time.Now()
	if t.CompletedAt != nil {
		endTime = *t.CompletedAt
	}

	return endTime.Sub(*t.StartedAt)
}

/* IsCompleted 检查任务是否已完成 */
func (t *VectorVerificationTask) IsCompleted() bool {
	return t.Status == TaskStatusCompleted || t.Status == TaskStatusFailed
}

/* IsRunning 检查任务是否正在运行 */
func (t *VectorVerificationTask) IsRunning() bool {
	return t.Status == TaskStatusRunning
}

/* MarkAsStarted 标记任务开始 */
func (t *VectorVerificationTask) MarkAsStarted() {
	now := time.Now()
	t.Status = TaskStatusRunning
	t.StartedAt = &now
}

/* MarkAsCompleted 标记任务完成 */
func (t *VectorVerificationTask) MarkAsCompleted() {
	now := time.Now()
	t.Status = TaskStatusCompleted
	t.CompletedAt = &now
}

/* MarkAsFailed 标记任务失败 */
func (t *VectorVerificationTask) MarkAsFailed(errorMsg string) {
	now := time.Now()
	t.Status = TaskStatusFailed
	t.CompletedAt = &now
	t.ErrorDetails = errorMsg
}

/* TaskSummary 任务摘要结构 */
type TaskSummary struct {
	TaskID            string     `json:"task_id"`
	Status            string     `json:"status"`
	TaskType          string     `json:"task_type"`
	Progress          float64    `json:"progress"`
	TotalCount        int        `json:"total_count"`
	ProcessedCount    int        `json:"processed_count"`
	VerifiedCount     int        `json:"verified_count"`
	MissingCount      int        `json:"missing_count"`
	ErrorCount        int        `json:"error_count"`
	Duration          string     `json:"duration"`
	EstimatedDuration *int       `json:"estimated_duration"`
	StartedAt         *time.Time `json:"started_at"`
	CompletedAt       *time.Time `json:"completed_at"`
	CreatorID         *uint      `json:"creator_id"`
}

func (t *VectorVerificationTask) ToSummary() *TaskSummary {
	return &TaskSummary{
		TaskID:            t.TaskID,
		Status:            t.Status,
		TaskType:          t.TaskType,
		Progress:          t.GetProgress(),
		TotalCount:        t.TotalCount,
		ProcessedCount:    t.ProcessedCount,
		VerifiedCount:     t.VerifiedCount,
		MissingCount:      t.MissingCount,
		ErrorCount:        t.ErrorCount,
		Duration:          t.GetDuration().String(),
		EstimatedDuration: t.EstimatedDuration,
		StartedAt:         t.StartedAt,
		CompletedAt:       t.CompletedAt,
		CreatorID:         t.CreatorID,
	}
}

func generateTaskID() string {
	return common.GenerateUniqueString()[:32] // 截取前32位
}
