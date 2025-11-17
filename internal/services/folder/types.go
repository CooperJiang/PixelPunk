package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
)

type FolderResponse struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ParentID    string          `json:"parent_id,omitempty"`
	Permission  string          `json:"permission"`
	Description string          `json:"description"`
	FileCount   int64           `json:"file_count"`
	HasChildren bool            `json:"has_children"`
	SortOrder   int             `json:"sort_order"`
	Level       int             `json:"level"`
	CreatedAt   common.JSONTime `json:"created_at"`
	UpdatedAt   common.JSONTime `json:"updated_at"`
}

/* PaginationInfo 分页信息（仍保留以兼容调用方） */
type PaginationInfo struct {
	Total       int64 `json:"total"`
	Size        int   `json:"size"`
	CurrentPage int   `json:"current_page"`
	LastPage    int   `json:"last_page"`
	From        int   `json:"from"`
	To          int   `json:"to"`
}

/* FolderFileAIInfo 简化AI信息（新语义） */
type FolderFileAIInfo struct {
	Description    string   `json:"description"`
	Tags           []string `json:"tags"`
	DominantColor  string   `json:"dominant_color"`
	Resolution     string   `json:"resolution"`
	IsNSFW         bool     `json:"is_nsfw"`
	NSFWScore      float64  `json:"nsfw_score"`
	NSFWEvaluation string   `json:"nsfw_evaluation"`
}

/* FileResponse 简化版文件响应结构体（新语义） */
type FileResponse struct {
	ID              string            `json:"id"`
	URL             string            `json:"url"`
	ThumbnailURL    string            `json:"thumb_url"`
	FullURL         string            `json:"full_url"`
	FullThumbURL    string            `json:"full_thumb_url"`
	ShortURL        string            `json:"short_url"`
	OriginalName    string            `json:"original_name"`
	DisplayName     string            `json:"display_name"`
	Size            int64             `json:"size"`
	Width           int               `json:"width"`
	Height          int               `json:"height"`
	Format          string            `json:"format"`
	AccessLevel     string            `json:"access_level"`
	FolderID        string            `json:"folder_id,omitempty"`
	CreatedAt       common.JSONTime   `json:"created_at"`
	UpdatedAt       common.JSONTime   `json:"updated_at"`
	IsDuplicate     bool              `json:"is_duplicate"`
	MD5Hash         string            `json:"md5_hash,omitempty"`
	IsRecommended   bool              `json:"is_recommended"`
	AIInfo          *FolderFileAIInfo `json:"ai_info,omitempty"`
	EXIFInfo        *models.FileEXIF  `json:"exif_info,omitempty"`
	StorageDuration string            `json:"storage_duration"`
	ExpiresAt       *common.JSONTime  `json:"expires_at"`
	IsTimeLimited   bool              `json:"is_time_limited"`
}

type FolderContentResponse struct {
	Folders    []*FolderResponse `json:"folders"`
	Files      []*FileResponse   `json:"files"`
	Pagination *PaginationInfo   `json:"pagination,omitempty"`
}

type TreeNodeResponse struct {
	ID       string              `json:"id"`
	Label    string              `json:"label"`
	Icon     string              `json:"icon"`
	Count    int64               `json:"count,omitempty"`
	Level    int                 `json:"level"`
	Children []*TreeNodeResponse `json:"children,omitempty"`
	Data     *FolderResponse     `json:"data,omitempty"`
}
