package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* FileCategory 文件分类模型 */
type FileCategory struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"size:50;not null" json:"name"` // 分类名称
	Description string `gorm:"type:text" json:"description"` // 分类描述
	UserID      uint   `gorm:"not null" json:"user_id"`      // 用户ID
	SortOrder   int    `gorm:"default:0" json:"sort_order"`  // 用户自定义排序

	Source     string `gorm:"size:20;default:user" json:"source"` // user/system_template/ai_suggestion
	TemplateID *uint  `gorm:"null" json:"template_id"`            // 来自系统模板时的模板ID

	FileCount int `gorm:"default:0" json:"file_count"` // 文件数量（冗余字段）

	Status    string          `gorm:"size:20;default:active" json:"status"` // active/archived
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
}

func (FileCategory) TableName() string {
	return "file_category"
}

func (c *FileCategory) BeforeCreate(tx *gorm.DB) error {
	var count int64
	err := tx.Model(&FileCategory{}).
		Where("user_id = ? AND name = ?", c.UserID, c.Name).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return gorm.ErrDuplicatedKey
	}

	if c.SortOrder == 0 {
		var maxOrder int
		tx.Model(&FileCategory{}).
			Where("user_id = ?", c.UserID).
			Select("COALESCE(MAX(sort_order), 0)").
			Scan(&maxOrder)
		c.SortOrder = maxOrder + 1
	}

	return nil
}

/* CategoryTemplate 系统分类模板模型 */
type CategoryTemplate struct {
	ID          uint   `gorm:"primarykey" json:"id"`
	Name        string `gorm:"size:50;not null;unique" json:"name"` // 模板分类名
	Description string `gorm:"type:text" json:"description"`        // 描述
	Icon        string `gorm:"size:50" json:"icon"`                 // 图标
	SortOrder   int    `gorm:"default:0" json:"sort_order"`         // 推荐排序
	IsPopular   bool   `gorm:"default:false" json:"is_popular"`     // 是否热门
	UsageCount  int    `gorm:"default:0" json:"usage_count"`        // 被采纳次数

	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
}

func (CategoryTemplate) TableName() string {
	return "category_template"
}

/* FileCategoryRelation 文件分类关联模型 */
type FileCategoryRelation struct {
	FileID     string          `gorm:"size:32;not null" json:"file_id"`      // 文件ID
	CategoryID uint            `gorm:"not null" json:"category_id"`          // 分类ID
	UserID     uint            `gorm:"not null" json:"user_id"`              // 用户ID（冗余字段）
	Source     string          `gorm:"size:20;default:manual" json:"source"` // manual/ai
	CreatedAt  common.JSONTime `json:"created_at"`

	Category FileCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	File     File         `gorm:"foreignKey:FileID" json:"file,omitempty"`
}

func (FileCategoryRelation) TableName() string {
	return "file_category_relation"
}

func (r *FileCategoryRelation) BeforeCreate(tx *gorm.DB) error {
	err := tx.Where("file_id = ?", r.FileID).Delete(&FileCategoryRelation{}).Error
	if err != nil {
		return err
	}

	if r.UserID == 0 {
		var file File
		err := tx.Where("id = ?", r.FileID).First(&file).Error
		if err != nil {
			return err
		}
		r.UserID = file.UserID
	}

	return nil
}
