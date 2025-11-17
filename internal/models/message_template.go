package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* MessageTemplate 消息模板模型 */
type MessageTemplate struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Type        string `gorm:"size:100;not null;uniqueIndex" json:"type"`
	Title       string `gorm:"size:200;not null" json:"title"`
	Content     string `gorm:"type:text;not null" json:"content"`
	Description string `gorm:"size:500" json:"description"`

	IsEnabled bool   `gorm:"default:true" json:"is_enabled"`
	SendEmail bool   `gorm:"default:false" json:"send_email"`
	ShowToast bool   `gorm:"default:true" json:"show_toast"`
	ToastType string `gorm:"size:20;default:info" json:"toast_type"`

	DefaultActionType  string `gorm:"size:50" json:"default_action_type"`
	DefaultActionText  string `gorm:"size:100" json:"default_action_text"`
	DefaultActionStyle string `gorm:"size:50;default:primary" json:"default_action_style"`
	ActionURLTemplate  string `gorm:"size:500" json:"action_url_template"`
}

func (MessageTemplate) TableName() string {
	return "message_template"
}

func (mt *MessageTemplate) BeforeCreate(tx *gorm.DB) error {
	if mt.ToastType == "" {
		mt.ToastType = "info"
	}
	if mt.DefaultActionStyle == "" {
		mt.DefaultActionStyle = "primary"
	}
	return nil
}

func (mt *MessageTemplate) GetPrimaryType() string {
	if mt.Type == "" {
		return "system"
	}
	for i, char := range mt.Type {
		if char == '.' {
			return mt.Type[:i]
		}
	}
	return mt.Type
}

/* IsTemplateEnabled 是否启用模板 */
func (mt *MessageTemplate) IsTemplateEnabled() bool {
	return mt.IsEnabled
}

func (mt *MessageTemplate) ShouldSendEmail() bool {
	return mt.SendEmail && mt.IsEnabled
}

func (mt *MessageTemplate) ShouldShowToast() bool {
	return mt.ShowToast && mt.IsEnabled
}
