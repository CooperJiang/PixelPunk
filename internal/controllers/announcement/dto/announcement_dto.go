package dto

import "pixelpunk/pkg/common"

type AnnouncementCreateDTO struct {
	Title    string `json:"title" binding:"required,max=255"`                         // 公告标题
	Content  string `json:"content" binding:"required"`                               // 内容（Markdown格式）
	Summary  string `json:"summary" binding:"max=500"`                                // 摘要
	IsPinned bool   `json:"is_pinned"`                                                // 是否置顶
	Status   string `json:"status" binding:"required,oneof=draft published archived"` // 状态
}

func (d *AnnouncementCreateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Title.required":   "公告标题不能为空",
		"Title.max":        "公告标题长度不能超过255个字符",
		"Content.required": "公告内容不能为空",
		"Summary.max":      "摘要长度不能超过500个字符",
		"Status.required":  "状态不能为空",
		"Status.oneof":     "状态必须是draft、published或archived",
	}
}

type AnnouncementUpdateDTO struct {
	Title    *string `json:"title" binding:"omitempty,max=255"`
	Content  *string `json:"content"`
	Summary  *string `json:"summary" binding:"omitempty,max=500"`
	IsPinned *bool   `json:"is_pinned"`
	Status   *string `json:"status" binding:"omitempty,oneof=draft published archived"`
}

func (d *AnnouncementUpdateDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Title.max":    "公告标题长度不能超过255个字符",
		"Summary.max":  "摘要长度不能超过500个字符",
		"Status.oneof": "状态必须是draft、published或archived",
	}
}

type AnnouncementQueryDTO struct {
	Status   string `form:"status" json:"status"`       // 状态筛选
	IsPinned *bool  `form:"is_pinned" json:"is_pinned"` // 是否置顶
	Page     int    `form:"page" json:"page"`           // 页码
	PageSize int    `form:"page_size" json:"page_size"` // 每页数量
	Keyword  string `form:"keyword" json:"keyword"`     // 关键词搜索
}

type AnnouncementResponseDTO struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Summary   string          `json:"summary"`
	IsPinned  bool            `json:"is_pinned"`
	Status    string          `json:"status"`
	ViewCount int             `json:"view_count"`
	CreatedBy uint            `json:"created_by"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
}

type AnnouncementListResponseDTO struct {
	Data       []AnnouncementResponseDTO `json:"data"`
	Pagination PaginationResponseDTO     `json:"pagination"`
}

type PaginationResponseDTO struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	Total     int64 `json:"total"`
	TotalPage int   `json:"total_page"`
}

type AnnouncementSimpleDTO struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Summary   string          `json:"summary"`
	IsPinned  bool            `json:"is_pinned"`
	ViewCount int             `json:"view_count"`
	CreatedAt common.JSONTime `json:"created_at"`
}

type AnnouncementDetailDTO struct {
	ID        uint            `json:"id"`
	Title     string          `json:"title"`
	Content   string          `json:"content"`
	Summary   string          `json:"summary"`
	IsPinned  bool            `json:"is_pinned"`
	ViewCount int             `json:"view_count"`
	CreatedAt common.JSONTime `json:"created_at"`
}

type PublicAnnouncementListDTO struct {
	Announcements []AnnouncementSimpleDTO `json:"announcements"`
	Total         int                     `json:"total"`
	Config        map[string]interface{}  `json:"config"` // 公告系统配置
}
