package models

import (
	"encoding/json"
	"pixelpunk/pkg/common"
)

/* ShareAccessLog 分享访问记录表 */
type ShareAccessLog struct {
	ID      string `gorm:"primarykey;size:32" json:"id"`
	ShareID string `gorm:"size:32;not null;index" json:"share_id"`

	AccessedAt common.JSONTime `json:"accessed_at"`                 // 访问时间
	IPAddress  string          `gorm:"size:50" json:"ip_address"`   // 访问者IP
	UserAgent  string          `gorm:"type:text" json:"user_agent"` // 用户代理
	Referer    string          `gorm:"size:255" json:"referer"`     // 来源页面

	VisitorName  string `gorm:"size:100" json:"visitor_name"`  // 访客姓名(可选)
	VisitorEmail string `gorm:"size:100" json:"visitor_email"` // 访客邮箱(可选)

	ViewedItems json.RawMessage `gorm:"type:json" json:"viewed_items"` // 查看了哪些项目
}

func (ShareAccessLog) TableName() string {
	return "share_access_log"
}
