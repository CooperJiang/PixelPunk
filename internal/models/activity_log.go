package models

import (
	"encoding/json"
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* ActivityLog 活动日志模型 */
type ActivityLog struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`

	UserID *uint  `json:"user_id"`                      // 用户ID，可空（系统操作）
	Type   string `gorm:"size:50;not null" json:"type"` // 活动类型

	Data json.RawMessage `gorm:"type:json;not null" json:"data"` // 结构化参数数据（必需）

	Module     string `gorm:"size:30" json:"module"`      // 涉及的模块: image, user, share, admin, system
	EntityType string `gorm:"size:30" json:"entity_type"` // 操作对象类型: image, folder, user, config
	EntityID   string `gorm:"size:100" json:"entity_id"`  // 具体对象ID

	IsVisible bool `gorm:"default:true" json:"is_visible"` // 是否对用户可见

	Tags string `gorm:"size:255" json:"tags"` // 标签，便于分类

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (ActivityLog) TableName() string {
	return "activity_log"
}

func (a *ActivityLog) BeforeCreate(tx *gorm.DB) error {
	if a.Type == "" {
		return gorm.ErrInvalidValue
	}
	if len(a.Data) == 0 {
		a.Data = json.RawMessage("{}")
	}
	return nil
}
