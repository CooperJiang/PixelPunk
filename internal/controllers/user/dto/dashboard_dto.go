package dto

import (
	"time"
)

type DashboardQueryDTO struct {
	RecentFilesLimit int `json:"recent_files_limit" form:"recent_files_limit" binding:"omitempty,min=1,max=50"`
}

type FolderQueryDTO struct {
	UserID     uint   `json:"user_id" form:"user_id"`         // 用户ID，可选，不提供则查询当前用户
	PublicOnly bool   `json:"public_only" form:"public_only"` // 是否只查询公开文件夹
	ParentID   string `json:"parent_id" form:"parent_id"`     // 父文件夹ID
	WithStats  bool   `json:"with_stats" form:"with_stats"`   // 是否包含统计信息
}

type FileBriefDTO struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	URL           string    `json:"url"`
	ThumbnailURL  string    `json:"thumb_url"`
	FullURL       string    `json:"full_url"`       // 完整的文件URL（包含域名）
	FullThumbURL  string    `json:"full_thumb_url"` // 完整的缩略图URL（包含域名）
	Size          int64     `json:"size"`
	SizeFormatted string    `json:"size_formatted"` // 格式化后的大小，如"52 KB"
	Views         int64     `json:"views"`
	FolderID      *string   `json:"folder_id"`
	FolderName    string    `json:"folder_name,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	StorageType   string    `json:"storage_type"`
}

type FolderStatDTO struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	FileCount          int    `json:"file_count"`
	TotalSize          int64  `json:"total_size"`
	TotalSizeFormatted string `json:"total_size_formatted"` // 格式化后的大小
	Description        string `json:"description,omitempty"`
	Permission         string `json:"permission"`
	CreatedAt          string `json:"created_at"`
}

type FolderListResponseDTO struct {
	Folders []FolderStatDTO `json:"folders"`
}

type DashboardResponseDTO struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`

	TotalFiles            int     `json:"total_files"`
	TotalViews            int64   `json:"total_views"`
	UsedStorage           int64   `json:"used_storage"`            // 字节
	UsedStorageFormatted  string  `json:"used_storage_formatted"`  // 格式化后的大小
	TotalStorage          int64   `json:"total_storage"`           // 字节
	TotalStorageFormatted string  `json:"total_storage_formatted"` // 格式化后的大小
	StorageUsagePercent   float64 `json:"storage_usage_percent"`   // 0-100
	TotalFolders          int     `json:"total_folders"`

	FilesUploadedToday     int   `json:"files_uploaded_today"`
	FilesUploadedThisWeek  int   `json:"files_uploaded_this_week"`
	FilesUploadedThisMonth int   `json:"files_uploaded_this_month"`
	ViewsToday             int64 `json:"views_today"`
	ViewsThisWeek          int64 `json:"views_this_week"`
	ViewsThisMonth         int64 `json:"views_this_month"`

	RecentFiles []FileBriefDTO  `json:"recent_files"`
	Folders     []FolderStatDTO `json:"folders"` // 保留向后兼容，但建议使用单独的文件夹接口
}

func (req *DashboardQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"RecentFilesLimit.min": "最近文件数量不能小于1",
		"RecentFilesLimit.max": "最近文件数量不能大于50",
	}
}

func (req *FolderQueryDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserID.min": "用户ID必须大于0",
	}
}
