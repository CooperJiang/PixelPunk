package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* Announcement 公告模型 */
type Announcement struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"deleted_at,omitempty"`

	Title     string `gorm:"size:255;not null" json:"title"`            // 标题
	Content   string `gorm:"type:text;not null" json:"content"`         // 内容（Markdown格式）
	Summary   string `gorm:"size:500" json:"summary"`                   // 摘要
	IsPinned  bool   `gorm:"default:false;index" json:"is_pinned"`      // 是否置顶
	Status    string `gorm:"size:20;default:draft;index" json:"status"` // 状态: draft, published, archived
	ViewCount int    `gorm:"default:0" json:"view_count"`               // 浏览次数
	CreatedBy uint   `gorm:"index" json:"created_by"`                   // 创建者ID
}

func (Announcement) TableName() string {
	return "announcements"
}

/* BeforeCreate GORM钩子 - 创建前设置默认值 */
func (a *Announcement) BeforeCreate(tx *gorm.DB) error {
	if a.Status == "" {
		a.Status = "draft"
	}
	return nil
}


func (a *Announcement) IsDraft() bool {
	return a.Status == "draft"
}

func (a *Announcement) IsPublished() bool {
	return a.Status == "published"
}

func (a *Announcement) IsArchived() bool {
	return a.Status == "archived"
}
