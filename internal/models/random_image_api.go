package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* RandomImageAPI 随机图片API模型 */
type RandomImageAPI struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID     uint    `gorm:"not null;index" json:"user_id"`
	Name       string  `gorm:"size:100;not null" json:"name"`                                  // API名称
	APIKey     string  `gorm:"size:32;not null;uniqueIndex:idx_random_api_key" json:"api_key"` // 随机密钥（rnd_xxxxxxxxxxxx）
	FolderID   *string `gorm:"size:32;index" json:"folder_id"`                                 // 绑定文件夹ID，NULL表示全部公开图片
	Status     int     `gorm:"default:1;index" json:"status"`                                  // 1:正常 2:禁用
	ReturnType string  `gorm:"size:20;not null;default:'redirect'" json:"return_type"`         // 返回类型：redirect(302重定向) 或 direct(直接返回图片)

	CallCount    int64            `gorm:"default:0" json:"call_count"` // 调用次数统计
	LastCalledAt *common.JSONTime `json:"last_called_at"`              // 最后调用时间

	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Folder *Folder `gorm:"foreignKey:FolderID" json:"folder,omitempty"`
}

/* RandomImageAPIStatus 随机图片API状态常量 */
const (
	RandomImageAPIStatusActive   = 1 // 正常状态
	RandomImageAPIStatusDisabled = 2 // 禁用状态
)

/* RandomImageAPIReturnType 随机图片API返回类型常量 */
const (
	RandomImageAPIReturnTypeRedirect = "redirect" // 302重定向
	RandomImageAPIReturnTypeDirect   = "direct"   // 直接返回图片
)

func (RandomImageAPI) TableName() string {
	return "random_image_apis"
}

func (r *RandomImageAPI) BeforeCreate(tx *gorm.DB) error {
	if r.Status == 0 {
		r.Status = RandomImageAPIStatusActive
	}
	if r.ReturnType == "" {
		r.ReturnType = RandomImageAPIReturnTypeRedirect
	}
	return nil
}

func (r *RandomImageAPI) IsActive() bool {
	return r.Status == RandomImageAPIStatusActive
}

func (r *RandomImageAPI) IsDisabled() bool {
	return r.Status == RandomImageAPIStatusDisabled
}
