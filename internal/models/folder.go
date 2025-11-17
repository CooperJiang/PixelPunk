package models

import (
	"errors"

	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* Folder 文件夹模型 */
type Folder struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`

	UserID        uint   `gorm:"not null;index" json:"user_id"`
	ParentID      string `gorm:"size:32;index" json:"parent_id"`
	Name          string `gorm:"size:100;not null" json:"name"`
	Permission    string `gorm:"size:10;not null;default:private;index" json:"permission"` // private 或 public
	Description   string `gorm:"size:500" json:"description"`                              // 文件夹描述
	IsRecommended bool   `gorm:"default:false;index" json:"is_recommended"`                // 是否是精选资源
	SortOrder     int    `gorm:"default:0" json:"sort_order"`                              // 排序值
}

func (Folder) TableName() string {
	return "folder"
}

/* IsRoot 是否是根文件夹 */
func (f *Folder) IsRoot() bool {
	return f.ParentID == ""
}

func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.Permission == "" {
		f.Permission = "private"
	}

	if f.SortOrder == 0 {
		var maxOrder int
		tx.Model(&Folder{}).
			Where("user_id = ? AND parent_id = ?", f.UserID, f.ParentID).
			Select("COALESCE(MAX(sort_order), 0)").
			Scan(&maxOrder)
		f.SortOrder = maxOrder + 1
	}

	return nil
}

/* GetFolderByIDAndUserID 根据ID和用户ID获取文件夹，确保权限验证 */
func GetFolderByIDAndUserID(db *gorm.DB, folderID string, userID uint) (*Folder, error) {
	var folder Folder
	err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("文件夹不存在或无权访问")
		}
		return nil, err
	}
	return &folder, nil
}
