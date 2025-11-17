package dto

import "pixelpunk/pkg/common"

type VectorListRequest struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	Status   string `json:"status" form:"status"`
	Model    string `json:"model" form:"model"`
	Keyword  string `json:"keyword" form:"keyword"`
}

func (r *VectorListRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码不能小于1",
		"PageSize.min": "每页数量不能小于1",
		"PageSize.max": "每页数量不能超过100",
	}
}

type VectorBatchActionRequest struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1"`
	Action  string   `json:"action" binding:"required,oneof=reset retry delete"`
}

func (r *VectorBatchActionRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要选择一个文件",
		"Action.required":  "操作类型不能为空",
		"Action.oneof":     "操作类型只能是reset、retry或delete",
	}
}

type VectorRetryRequest struct {
	FileID string `json:"file_id" binding:"required"`
}

func (r *VectorRetryRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileID.required": "文件ID不能为空",
	}
}

type VectorStatsResponse struct {
	TotalCount     int64 `json:"total_count"`
	PendingCount   int64 `json:"pending_count"`
	RunningCount   int64 `json:"running_count"`
	CompletedCount int64 `json:"completed_count"`
	FailedCount    int64 `json:"failed_count"`
	ResetCount     int64 `json:"reset_count"`
	ActiveWorkers  int   `json:"active_workers"`
	MaxWorkers     int   `json:"max_workers"`
}

type VectorItemResponse struct {
	ID                 uint             `json:"id"`
	FileID             string           `json:"file_id"`
	Status             string           `json:"status"`
	Model              string           `json:"model"`
	Description        string           `json:"description"`
	Dimension          int              `json:"dimension"`
	RetryCount         int              `json:"retry_count"`
	ErrorMessage       string           `json:"error_message,omitempty"`
	LastRetryAt        *common.JSONTime `json:"last_retry_at,omitempty"`
	ProcessingDuration int              `json:"processing_duration"`
	CreatedAt          common.JSONTime  `json:"created_at"`
	UpdatedAt          common.JSONTime  `json:"updated_at"`

	FileName     string `json:"file_name,omitempty"`
	FileURL      string `json:"file_url,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`

	VectorPreview []float32    `json:"vector_preview,omitempty"`
	VectorNorm    float64      `json:"vector_norm,omitempty"`
	VectorStats   *VectorStats `json:"vector_stats,omitempty"`
}

type VectorStats struct {
	Min    float32 `json:"min"`
	Max    float32 `json:"max"`
	Mean   float32 `json:"mean"`
	StdDev float32 `json:"stddev"`
}

type VectorListResponse struct {
	Data       []VectorItemResponse `json:"data"`
	Pagination PaginationResponse   `json:"pagination"`
}

type PaginationResponse struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type VectorLogListRequest struct {
	Page     int    `json:"page" form:"page" binding:"min=1"`
	PageSize int    `json:"page_size" form:"page_size" binding:"min=1,max=100"`
	FileID   string `json:"file_id" form:"file_id"`
	Action   string `json:"action" form:"action"`
	Model    string `json:"model" form:"model"`
}

func (r *VectorLogListRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":     "页码不能小于1",
		"PageSize.min": "每页数量不能小于1",
		"PageSize.max": "每页数量不能超过100",
	}
}

type VectorLogItemResponse struct {
	ID        uint                   `json:"id"`
	FileID    string                 `json:"file_id"`
	Action    string                 `json:"action"`
	Type      string                 `json:"type"`
	Data      map[string]interface{} `json:"data"`
	Message   string                 `json:"message,omitempty"`
	Model     string                 `json:"model"`
	Duration  int                    `json:"duration"`
	ErrorCode string                 `json:"error_code,omitempty"`
	TaskID    string                 `json:"task_id,omitempty"`
	CreatedAt common.JSONTime        `json:"created_at"`
}

type VectorLogListResponse struct {
	Data       []VectorLogItemResponse `json:"data"`
	Pagination PaginationResponse      `json:"pagination"`
}

type VectorRegenerateAllRequest struct {
	Force bool `json:"force"`
}

func (r *VectorRegenerateAllRequest) GetValidationMessages() map[string]string {
	return map[string]string{}
}
