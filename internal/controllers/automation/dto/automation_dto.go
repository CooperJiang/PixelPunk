package dto


// TaggingTaskStats 打标任务统计信息
type TaggingTaskStats struct {
	NoneCount       int `json:"none_count"`
	PendingCount    int `json:"pending_count"`
	ProcessingCount int `json:"processing_count"`
	DoneCount       int `json:"done_count"`
	FailedCount     int `json:"failed_count"`
	IgnoredCount    int `json:"ignored_count"`
	TotalCount      int `json:"total_count"`
	QueuePosition   int `json:"queue_position,omitempty"` // 可选：队列位置
}

// VectorTaskStats 向量化任务统计信息
type VectorTaskStats struct {
	PendingCount    int `json:"pending_count"`
	ProcessingCount int `json:"processing_count"`
	CompletedCount  int `json:"completed_count"`
	FailedCount     int `json:"failed_count"`
	ResetCount      int `json:"reset_count"`
	TotalCount      int `json:"total_count"`
	QueuePosition   int `json:"queue_position,omitempty"` // 可选：队列位置
}

// SystemStatus 系统队列状态
type SystemStatus struct {
	TaggingEnabled     bool `json:"tagging_enabled"`
	VectorEnabled      bool `json:"vector_enabled"`
	TaggingConcurrency int  `json:"tagging_concurrency"`
	VectorConcurrency  int  `json:"vector_concurrency"`
	TotalQueueLength   int  `json:"total_queue_length"`
}

// AutomationOverviewResponse 自动任务总览响应
type AutomationOverviewResponse struct {
	Tagging      TaggingTaskStats `json:"tagging"`
	Vector       VectorTaskStats  `json:"vector"`
	SystemStatus SystemStatus     `json:"system_status"`
}

// TaggingTaskItem 打标任务详情项
type TaggingTaskItem struct {
	ID            string `json:"id"`
	FileID        string `json:"file_id"`
	FileName      string `json:"file_name"`
	ThumbnailURL  string `json:"thumbnail_url,omitempty"`
	Status        string `json:"status"`
	Tries         int    `json:"tries"`
	ErrorMessage  string `json:"error_message,omitempty"`
	Format        string `json:"format"`
	SizeFormatted string `json:"size_formatted"`
	Resolution    string `json:"resolution"`
	UpdatedAt     string `json:"updated_at"`
	CreatedAt     string `json:"created_at"`
}

// VectorTaskItem 向量任务详情项
type VectorTaskItem struct {
	ID           int    `json:"id"`
	FileID       string `json:"file_id"`
	FileName     string `json:"file_name"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	Status       string `json:"status"`
	Model        string `json:"model"`
	RetryCount   int    `json:"retry_count"`
	ErrorMessage string `json:"error_message,omitempty"`
	UpdatedAt    string `json:"updated_at"`
	CreatedAt    string `json:"created_at"`
}

// TaskListQuery 任务列表查询参数
type TaskListQuery struct {
	Status string `form:"status"`
	Page   int    `form:"page" binding:"min=1"`
	Limit  int    `form:"limit" binding:"min=1,max=100"`
}

// TaskListResponse 任务列表响应
type TaskListResponse struct {
	Data  interface{} `json:"data"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
}
