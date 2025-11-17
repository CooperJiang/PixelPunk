package models

import (
	"pixelpunk/pkg/common"
	"time"
)

/* ShareAccessToken 分享访问令牌模型 */
type ShareAccessToken struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`

	ShareID   string          `gorm:"size:32;index" json:"share_id"`  // 关联的分享ID
	ShareKey  string          `gorm:"size:32;index" json:"share_key"` // 分享密钥，冗余存储方便查询
	Token     string          `gorm:"size:64"`                        // 访问令牌
	ExpiredAt common.JSONTime `json:"expired_at"`                     // 过期时间
	ClientIP  string          `gorm:"size:50" json:"client_ip"`       // 客户端IP
	UserAgent string          `gorm:"size:255" json:"user_agent"`     // 用户代理
	IsAdmin   bool            `gorm:"default:false" json:"is_admin"`  // 是否为管理员令牌
}

func (ShareAccessToken) TableName() string {
	return "share_access_token"
}

/* IsExpired 判断令牌是否已过期 */
func (t *ShareAccessToken) IsExpired() bool {
	return time.Now().After(time.Time(t.ExpiredAt))
}
