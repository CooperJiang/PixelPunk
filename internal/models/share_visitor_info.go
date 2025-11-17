package models

import (
	"pixelpunk/pkg/common"
)

/* ShareVisitorInfo 分享访客信息模型 */
type ShareVisitorInfo struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`

	ShareID      string          `gorm:"size:32;index" json:"share_id"`  // 关联的分享ID
	ShareKey     string          `gorm:"size:32;index" json:"share_key"` // 分享密钥
	VisitorName  string          `gorm:"size:100" json:"visitor_name"`   // 访客姓名
	VisitorEmail string          `gorm:"size:100" json:"visitor_email"`  // 访客邮箱
	IPAddress    string          `gorm:"size:50" json:"ip_address"`      // 访客IP地址
	UserAgent    string          `gorm:"size:255" json:"user_agent"`     // 用户代理
	Referer      string          `gorm:"size:255" json:"referer"`        // 来源页面
	VisitCount   int             `gorm:"default:1" json:"visit_count"`   // 访问次数
	LastVisitAt  common.JSONTime `json:"last_visit_at"`                  // 最后访问时间
}

func (ShareVisitorInfo) TableName() string {
	return "share_visitor_info"
}
