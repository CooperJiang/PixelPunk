package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* Message 用户消息模型 */
type Message struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID uint   `gorm:"not null" json:"user_id"`
	Type   string `gorm:"size:100;not null" json:"type"`
	Data   string `gorm:"type:text;not null" json:"data"` // 结构化消息数据，前端根据 type 动态生成标题和内容

	Status    int              `gorm:"default:1" json:"status"`   // 1:未读 2:已读 3:已删除
	Priority  int              `gorm:"default:2" json:"priority"` // 1:高 2:中 3:低
	ReadAt    *common.JSONTime `gorm:"column:read_at" json:"read_at"`
	ExpiresAt *common.JSONTime `gorm:"column:expires_at" json:"expires_at"`

	RelatedType string `gorm:"size:50" json:"related_type"`
	RelatedID   string `gorm:"size:100" json:"related_id"`
	RelatedData string `gorm:"type:text" json:"related_data"`

	IsActionable       bool   `gorm:"default:false" json:"is_actionable"`
	ActionType         string `gorm:"size:50" json:"action_type"`
	ActionURL          string `gorm:"size:500" json:"action_url"`
	ActionText         string `gorm:"size:100" json:"action_text"`
	ActionStyle        string `gorm:"size:50;default:primary" json:"action_style"`
	RequiresPermission string `gorm:"size:100" json:"requires_permission"`

	MetaData string `gorm:"type:text" json:"meta_data"`
}

func (Message) TableName() string {
	return "message"
}

func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.Status == 0 {
		m.Status = common.MessageStatusUnread
	}
	if m.Priority == 0 {
		m.Priority = common.MessagePriorityNormal
	}
	if m.ActionStyle == "" {
		m.ActionStyle = "primary"
	}

	// 确保必需的 JSON 字段有默认值
	if m.Data == "" {
		m.Data = "{}"
	}
	if m.RelatedData == "" {
		m.RelatedData = "{}"
	}
	if m.MetaData == "" {
		m.MetaData = "{}"
	}

	return nil
}

/* IsUnread 是否未读 */
func (m *Message) IsUnread() bool {
	return m.Status == common.MessageStatusUnread
}

/* IsRead 是否已读 */
func (m *Message) IsRead() bool {
	return m.Status == common.MessageStatusRead
}

/* IsDeleted 是否已删除 */
func (m *Message) IsDeleted() bool {
	return m.Status == common.MessageStatusDeleted
}

/* IsHighPriority 是否高优先级 */
func (m *Message) IsHighPriority() bool {
	return m.Priority == common.MessagePriorityHigh
}

func (m *Message) GetPrimaryType() string {
	if m.Type == "" {
		return "system"
	}
	for i, char := range m.Type {
		if char == '.' {
			return m.Type[:i]
		}
	}
	return m.Type
}
