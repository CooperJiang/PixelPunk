package dto

// UploadFileDTO 上传文件请求DTO
type UploadFileDTO struct {
	FolderID        string `form:"folder_id" json:"folder_id"`
	FilePath        string `form:"file_path" json:"file_path"`
	AccessLevel     string `form:"access_level" json:"access_level"`
	Optimize        bool   `form:"optimize" json:"optimize"`
	StorageDuration string `form:"storage_duration" json:"storage_duration"`
	Watermark       string `form:"watermark" json:"watermark"`
	// 兼容旧参数（将逐步淘汰）
	WatermarkEnabled bool   `form:"watermark_enabled" json:"watermark_enabled"`
	WatermarkType    string `form:"watermark_type" json:"watermark_type" binding:"omitempty,oneof=file"`
	WatermarkConfig  string `form:"watermark_config" json:"watermark_config"`
}

func (d *UploadFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"AccessLevel.oneof":   "访问级别必须是 public、private 或 protected",
		"WatermarkType.oneof": "水印类型必须是 file",
	}
}

type UpdateFileDTO struct {
	FolderID    string `json:"folder_id"`
	AccessLevel string `json:"access_level" binding:"omitempty,oneof=public private protected"`
	Name        string `json:"name" binding:"omitempty,max=255"`
}

func (d *UpdateFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"AccessLevel.oneof": "访问级别必须是 public、private 或 protected",
		"Name.max":          "文件名称不能超过255个字符",
	}
}

// FileListQueryDTO 文件列表查询DTO
type FileListQueryDTO struct {
	Page          int    `form:"page" binding:"omitempty,min=1"`
	Size          int    `form:"size" binding:"omitempty,min=1,max=100"`
	FolderID      string `form:"folder_id"`
	Sort          string `form:"sort" binding:"omitempty,oneof=newest oldest name size width height quality nsfw_score"`
	AccessLevel   string `form:"access_level" binding:"omitempty,oneof=public private protected"`
	Keyword       string `form:"keyword" binding:"omitempty,max=100"`
	Tags          string `form:"tags"`           // 逗号分隔的标签字符串
	CategoryID    string `form:"categoryId"`     // 逗号分隔的分类ID字符串
	DominantColor string `form:"dominant_color"` // 逗号分隔的颜色字符串
	Resolution    string `form:"resolution"`
	MinWidth      int    `form:"min_width"`
	MaxWidth      int    `form:"max_width"`
	MinHeight     int    `form:"min_height"`
	MaxHeight     int    `form:"max_height"`
}

func (d *FileListQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Page.min":          "页码必须大于等于1",
		"Size.min":          "每页数量必须大于等于1",
		"Size.max":          "每页数量必须小于等于100",
		"Sort.oneof":        "排序方式必须是 newest、oldest、name、size、width、height、quality 或 nsfw_score",
		"AccessLevel.oneof": "访问级别必须是 public、private 或 protected",
		"Keyword.max":       "搜索关键字不能超过100个字符",
	}
}

// FileStatsQueryDTO 文件统计数据查询DTO
type FileStatsQueryDTO struct {
	Period string `form:"period" binding:"omitempty,oneof=day week month year"`
}

func (d *FileStatsQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Period.oneof": "时间范围必须是 day、week、month 或 year",
	}
}

// AdminRecommendFileDTO 管理员推荐文件DTO
type AdminRecommendFileDTO struct {
	FileID string `json:"file_id" binding:"required"`
}

func (d *AdminRecommendFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileID.required": "文件ID不能为空",
	}
}

// AdminBatchRecommendFileDTO 管理员批量推荐文件DTO
type AdminBatchRecommendFileDTO struct {
	FileIDs       []string `json:"file_ids" binding:"required,min=1"`
	IsRecommended bool     `json:"is_recommended"`
}

func (d *AdminBatchRecommendFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要选择一个文件",
	}
}

// AdminBatchDeleteFileDTO 管理员批量删除文件DTO
type AdminBatchDeleteFileDTO struct {
	FileIDs []string `json:"file_ids" binding:"required,min=1"`
}

func (d *AdminBatchDeleteFileDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要一个文件",
	}
}

// ReorderFilesDTO 重新排序文件请求DTO
type ReorderFilesDTO struct {
	FolderID string   `json:"folder_id"`                         // 文件夹ID，空字符串表示根目录
	FileIDs  []string `json:"file_ids" binding:"required,min=1"` // 按新顺序排列的文件ID数组
}

func (d *ReorderFilesDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要一个文件",
	}
}

// BatchUploadFailure 批量上传失败信息
type BatchUploadFailure struct {
	Filename string `json:"filename"` // 失败的文件名
	Error    string `json:"error"`    // 错误信息
	Index    int    `json:"index"`    // 文件在批量上传中的索引
}

// BatchUploadResultDTO 批量上传结果DTO
type BatchUploadResultDTO struct {
	TotalFiles   int                  `json:"total_files"`   // 总文件数
	SuccessCount int                  `json:"success_count"` // 成功上传数量
	FailureCount int                  `json:"failure_count"` // 失败上传数量
	SuccessFiles []interface{}        `json:"success_files"` // 成功上传的文件信息（使用FileDetailResponse）
	Failures     []BatchUploadFailure `json:"failures"`      // 失败的文件信息
	Message      string               `json:"message"`       // 总体处理结果消息
}

// CheckDuplicateDTO MD5预检查请求DTO
type CheckDuplicateDTO struct {
	MD5      string `json:"md5" binding:"required,len=32"`
	FileName string `json:"filename" binding:"required,max=255"`
	FileSize int64  `json:"file_size" binding:"required,min=1"`
}

func (d *CheckDuplicateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"MD5.required":      "文件MD5不能为空",
		"MD5.len":           "MD5长度必须为32位",
		"FileName.required": "文件名不能为空",
		"FileName.max":      "文件名不能超过255个字符",
		"FileSize.required": "文件大小不能为空",
		"FileSize.min":      "文件大小必须大于0",
	}
}

// InstantUploadDTO 秒传请求DTO
type InstantUploadDTO struct {
	MD5              string `json:"md5" binding:"required,len=32"`
	FileName         string `json:"filename" binding:"required,max=255"`
	FileSize         int64  `json:"file_size" binding:"required,min=1"`
	FolderID         string `json:"folder_id"`
	AccessLevel      string `json:"access_level" binding:"omitempty,oneof=public private protected"`
	Optimize         bool   `json:"optimize"`
	StorageDuration  string `json:"storage_duration"`
	Watermark        string `json:"watermark"`
	WatermarkEnabled bool   `json:"watermark_enabled"`
	WatermarkType    string `json:"watermark_type" binding:"omitempty,oneof=file"`
	WatermarkConfig  string `json:"watermark_config"`
}

func (d *InstantUploadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"MD5.required":        "文件MD5不能为空",
		"MD5.len":             "MD5长度必须为32位",
		"FileName.required":   "文件名不能为空",
		"FileName.max":        "文件名不能超过255个字符",
		"FileSize.required":   "文件大小不能为空",
		"FileSize.min":        "文件大小必须大于0",
		"AccessLevel.oneof":   "访问级别必须是 public、private 或 protected",
		"WatermarkType.oneof": "水印类型必须是 file",
	}
}

// GuestUploadDTO 游客上传文件请求DTO
type GuestUploadDTO struct {
	AccessLevel      string `form:"access_level" json:"access_level"`
	Optimize         bool   `form:"optimize" json:"optimize"`
	StorageDuration  string `form:"storage_duration" json:"storage_duration" binding:"required"`
	Fingerprint      string `form:"fingerprint" json:"fingerprint"`
	Watermark        string `form:"watermark" json:"watermark"`
	WatermarkEnabled bool   `form:"watermark_enabled" json:"watermark_enabled"`
	WatermarkType    string `form:"watermark_type" json:"watermark_type" binding:"omitempty,oneof=file"`
	WatermarkConfig  string `form:"watermark_config" json:"watermark_config"`
}

func (d *GuestUploadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"AccessLevel.oneof":        "访问级别必须是 public、private 或 protected",
		"StorageDuration.required": "游客上传必须指定存储时长",
		"WatermarkType.oneof":      "水印类型必须是 file",
	}
}

// MoveFilesRequest 批量移动文件请求DTO
type MoveFilesRequest struct {
	FileIDs        []string `json:"file_ids" binding:"required,min=1,max=50"` // 要移动的文件ID列表
	TargetFolderID string   `json:"target_folder_id"`                         // 目标文件夹ID，空字符串表示根目录
}

func (d *MoveFilesRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要选择一个文件",
		"FileIDs.max":      "单次最多只能移动50个文件",
	}
}

// ReorderFilesRequest 重新排序文件请求DTO
type ReorderFilesRequest struct {
	FolderID string   `json:"folder_id"`                         // 文件夹ID，空字符串表示根目录
	FileIDs  []string `json:"file_ids" binding:"required,min=1"` // 按新顺序排列的文件ID数组
}

func (d *ReorderFilesRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileIDs.required": "文件ID列表不能为空",
		"FileIDs.min":      "至少需要一个文件",
	}
}
