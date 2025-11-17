package models

import (
	"pixelpunk/pkg/common"
	"time"
)

// PasswordResetToken 密码重置令牌模型
type PasswordResetToken struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID    uint       `gorm:"not null;index" json:"user_id"`
	Token     string     `gorm:"size:64;not null;uniqueIndex" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	IsUsed    bool       `gorm:"default:false" json:"is_used"`
	UsedAt    *time.Time `json:"used_at"`
	IPAddress string     `gorm:"size:45" json:"ip_address"` // 请求IP

	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (PasswordResetToken) TableName() string {
	return "password_reset_token"
}

// IsExpired 检查token是否已过期
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// IsValid 检查token是否有效（未使用且未过期）
func (p *PasswordResetToken) IsValid() bool {
	return !p.IsUsed && !p.IsExpired()
}
