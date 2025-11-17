package dto

type InitChunkedUploadDTO struct {
	FileName        string      `json:"file_name" binding:"required,max=255"`
	FileSize        int64       `json:"file_size" binding:"required,min=1"`
	FileMD5         string      `json:"file_md5" binding:"required,len=32"`
	MimeType        string      `json:"mime_type" binding:"required"`
	ChunkSize       int64       `json:"chunk_size" binding:"required,min=1048576,max=10485760"` // 1MB-10MB
	FolderID        string      `json:"folder_id"`
	AccessLevel     string      `json:"access_level" binding:"omitempty,oneof=public private protected"`
	Optimize        bool        `json:"optimize"`
	WatermarkConfig interface{} `json:"watermark_config"`
}

func (d *InitChunkedUploadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"FileName.required":  "文件名不能为空",
		"FileName.max":       "文件名不能超过255个字符",
		"FileSize.required":  "文件大小不能为空",
		"FileSize.min":       "文件大小必须大于0",
		"FileMD5.required":   "文件MD5不能为空",
		"FileMD5.len":        "文件MD5必须为32位",
		"MimeType.required":  "文件类型不能为空",
		"ChunkSize.required": "分片大小不能为空",
		"ChunkSize.min":      "分片大小不能小于1MB",
		"ChunkSize.max":      "分片大小不能大于10MB",
		"AccessLevel.oneof":  "访问级别必须是 public、private 或 protected",
	}
}

type UploadChunkDTO struct {
	SessionID   string `form:"session_id" binding:"required,len=32"`
	ChunkNumber int    `form:"chunk_number" binding:"gte=0"`
	ChunkMD5    string `form:"chunk_md5" binding:"required,len=32"`
}

func (d *UploadChunkDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SessionID.required": "会话ID不能为空",
		"SessionID.len":      "会话ID必须为32位",
		"ChunkNumber.gte":    "分片序号不能小于0",
		"ChunkMD5.required":  "分片MD5不能为空",
		"ChunkMD5.len":       "分片MD5必须为32位",
	}
}

type CompleteChunkedUploadDTO struct {
	SessionID string `json:"session_id" binding:"required,len=32"`
}

func (d *CompleteChunkedUploadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SessionID.required": "会话ID不能为空",
		"SessionID.len":      "会话ID必须为32位",
	}
}

type ChunkedUploadStatusDTO struct {
	SessionID string `form:"session_id" binding:"required,len=32"`
}

func (d *ChunkedUploadStatusDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SessionID.required": "会话ID不能为空",
		"SessionID.len":      "会话ID必须为32位",
	}
}

type CancelChunkedUploadDTO struct {
	SessionID string `json:"session_id" binding:"required,len=32"`
}

func (d *CancelChunkedUploadDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"SessionID.required": "会话ID不能为空",
		"SessionID.len":      "会话ID必须为32位",
	}
}

type ChunkedUploadResponse struct {
	SessionID      string `json:"session_id"`
	Status         string `json:"status"`
	Progress       int    `json:"progress"`
	TotalChunks    int    `json:"total_chunks"`
	UploadedChunks int    `json:"uploaded_chunks"`
	Message        string `json:"message"`
}

type ChunkedUploadCompleteResponse struct {
	SessionID string      `json:"session_id"`
	FileID    string      `json:"file_id"`
	Message   string      `json:"message"`
	FileInfo  interface{} `json:"file_info,omitempty"`
}
