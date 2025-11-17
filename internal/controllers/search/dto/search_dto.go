package dto

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/storage"
	"time"
)

type VectorSearchRequest struct {
	Query     string  `json:"query" binding:"required,min=1,max=500"`    // 搜索查询文本
	Limit     int     `json:"limit" binding:"omitempty,min=1,max=100"`   // 返回结果数量限制
	Threshold float32 `json:"threshold" binding:"omitempty,min=0,max=1"` // 相似度阈值
}

func (r *VectorSearchRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"Query.required": "搜索查询不能为空",
		"Query.min":      "搜索查询至少需要1个字符",
		"Query.max":      "搜索查询不能超过500个字符",
		"Limit.min":      "结果数量限制不能小于1",
		"Limit.max":      "结果数量限制不能超过100",
		"Threshold.min":  "相似度阈值不能小于0",
		"Threshold.max":  "相似度阈值不能大于1",
	}
}

type VectorSearchResult struct {
	FileID      string    `json:"file_id"`             // 文件ID
	Similarity  float32   `json:"similarity"`          // 相似度分数
	Description string    `json:"description"`         // 文件描述
	FileInfo    *FileInfo `json:"file_info,omitempty"` // 文件基本信息
}

type FileInfo struct {
	ID            string `json:"id"`             // 文件数据库ID
	OriginalName  string `json:"original_name"`  // 原始文件名
	DisplayName   string `json:"display_name"`   // 显示名称
	Size          int64  `json:"size"`           // 文件大小
	Width         int    `json:"width"`          // 宽度（若为文件）
	Height        int    `json:"height"`         // 高度（若为文件）
	Format        string `json:"format"`         // 文件格式
	URL           string `json:"url"`            // 文件访问URL
	ThumbURL      string `json:"thumb_url"`      // 缩略图访问URL
	FullURL       string `json:"full_url"`       // 完整URL
	FullThumbURL  string `json:"full_thumb_url"` // 完整缩略图URL
	CreatedAt     string `json:"created_at"`     // 创建时间
	SizeFormatted string `json:"size_formatted"` // 格式化的文件大小
	Resolution    string `json:"resolution"`     // 分辨率
}

type VectorSearchResponse struct {
	Query       string               `json:"query"`        // 原始查询
	Total       int                  `json:"total"`        // 结果总数
	Results     []VectorSearchResult `json:"results"`      // 搜索结果列表
	ProcessTime string               `json:"process_time"` // 处理耗时
	UsedCache   bool                 `json:"used_cache"`   // 是否使用了缓存
}

type VectorStatsRequest struct {
	UserID uint `json:"user_id" binding:"omitempty"` // 用户ID（可选，管理员可查看所有用户）
}

type VectorStatsResponse struct {
	TotalVectors     int64 `json:"total_vectors"`     // 总向量数
	CompletedVectors int64 `json:"completed_vectors"` // 已完成向量数
	PendingVectors   int64 `json:"pending_vectors"`   // 待处理向量数
	FailedVectors    int64 `json:"failed_vectors"`    // 失败向量数
	UserVectors      int64 `json:"user_vectors"`      // 用户向量数（如果指定了用户ID）
	StorageSize      int64 `json:"storage_size"`      // 存储大小（字节）
	AvgDimension     int   `json:"avg_dimension"`     // 平均向量维度
}

type VectorHealthResponse struct {
	Status      string `json:"status"`       // 状态：ok/error
	Enabled     bool   `json:"enabled"`      // 是否启用
	Provider    string `json:"provider"`     // 向量化提供者
	Model       string `json:"model"`        // 使用的模型
	DatabaseOK  bool   `json:"database_ok"`  // 数据库连接状态
	EmbeddingOK bool   `json:"embedding_ok"` // 向量化服务状态
	Message     string `json:"message"`      // 状态消息
}

type UserVectorSearchRequest struct {
	Query string `json:"query" binding:"required,min=1,max=500"` // 搜索查询文本
	Page  int    `json:"page" binding:"omitempty,min=1"`         // 页码
}

func (r *UserVectorSearchRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"Query.required": "搜索查询不能为空",
		"Query.min":      "搜索查询至少需要1个字符",
		"Query.max":      "搜索查询不能超过500个字符",
		"Page.min":       "页码必须大于等于1",
	}
}

type GalleryVectorSearchRequest struct {
	Query string `json:"query" binding:"required,min=1,max=500"` // 搜索查询文本
	Page  int    `json:"page" binding:"omitempty,min=1"`         // 页码
}

func (r *GalleryVectorSearchRequest) GetValidationMessages() map[string]string {
	return map[string]string{
		"Query.required": "搜索查询不能为空",
		"Query.min":      "搜索查询至少需要1个字符",
		"Query.max":      "搜索查询不能超过500个字符",
		"Page.min":       "页码必须大于等于1",
	}
}

type VectorSearchFileResult struct {
	models.File         // 嵌入完整的文件模型
	Similarity  float32 `json:"similarity"` // 添加相似度字段
}

func ConvertFileToInfo(file *models.File) *FileInfo {
	if file == nil {
		return nil
	}

	fullURL, fullThumbURL, _ := storage.GetFullURLs(*file)

	return &FileInfo{
		ID:            file.ID,
		OriginalName:  file.OriginalName,
		DisplayName:   file.DisplayName,
		Size:          file.Size,
		Width:         file.Width,
		Height:        file.Height,
		Format:        file.Format,
		URL:           file.URL,
		ThumbURL:      file.ThumbURL,
		FullURL:       fullURL,
		FullThumbURL:  fullThumbURL,
		CreatedAt:     time.Time(file.CreatedAt).Format("2006-01-02 15:04:05"),
		SizeFormatted: file.SizeFormatted,
		Resolution:    file.Resolution,
	}
}
