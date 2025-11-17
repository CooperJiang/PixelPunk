package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* GuestUploadLimit 游客上传限制模型 */
type GuestUploadLimit struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Fingerprint string `gorm:"size:64;not null;index:idx_guest_fingerprint" json:"fingerprint"` // 浏览器指纹
	IP          string `gorm:"size:45;not null;index:idx_guest_ip" json:"ip"`                   // IP地址
	UserAgent   string `gorm:"size:500" json:"user_agent"`                                      // 用户代理信息

	DailyLimit   int        `gorm:"not null;default:10" json:"daily_limit"` // 每日上传限制
	HourlyLimit  int        `gorm:"not null;default:5" json:"hourly_limit"` // 每小时上传限制
	MinuteLimit  int        `gorm:"not null;default:1" json:"minute_limit"` // 每分钟上传限制
	IsBanned     bool       `gorm:"default:false" json:"is_banned"`         // 是否被封禁
	BanReason    string     `gorm:"size:255" json:"ban_reason"`             // 封禁原因
	BanExpiresAt *time.Time `gorm:"index" json:"ban_expires_at"`            // 封禁过期时间

	DailyCount    int        `gorm:"not null;default:0" json:"daily_count"`  // 今日已上传数量
	HourlyCount   int        `gorm:"not null;default:0" json:"hourly_count"` // 本小时已上传数量
	MinuteCount   int        `gorm:"not null;default:0" json:"minute_count"` // 本分钟已上传数量
	LastUploadAt  *time.Time `json:"last_upload_at"`                         // 最后上传时间
	LastResetDate *time.Time `gorm:"index" json:"last_reset_date"`           // 最后重置日期（用于日限制）
	LastResetHour *time.Time `gorm:"index" json:"last_reset_hour"`           // 最后重置小时（用于小时限制）
}

func (GuestUploadLimit) TableName() string {
	return "guest_upload_limit"
}

func (g *GuestUploadLimit) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	g.LastResetDate = &now
	g.LastResetHour = &now
	return nil
}

/* IsExpired 检查限制是否已过期 */
func (g *GuestUploadLimit) IsExpired() bool {
	if g.BanExpiresAt == nil {
		return false
	}
	return time.Now().After(*g.BanExpiresAt)
}

func (g *GuestUploadLimit) CanUpload() bool {
	g.ResetIfNeeded()

	if g.IsBanned && !g.IsExpired() {
		return false
	}

	return g.MinuteCount < g.MinuteLimit &&
		g.HourlyCount < g.HourlyLimit &&
		g.DailyCount < g.DailyLimit
}

func (g *GuestUploadLimit) ResetIfNeeded() {
	now := time.Now()

	if g.LastResetHour == nil ||
		now.Minute() != g.LastResetHour.Minute() ||
		now.Hour() != g.LastResetHour.Hour() {
		g.MinuteCount = 0
	}

	if g.LastResetHour == nil ||
		now.Hour() != g.LastResetHour.Hour() ||
		now.Day() != g.LastResetHour.Day() {
		g.HourlyCount = 0
		g.LastResetHour = &now
	}

	if g.LastResetDate == nil ||
		now.Year() != g.LastResetDate.Year() ||
		now.Month() != g.LastResetDate.Month() ||
		now.Day() != g.LastResetDate.Day() {
		g.DailyCount = 0
		g.LastResetDate = &now
	}
}

/* IncrementUploadCount 增加上传计数 */
func (g *GuestUploadLimit) IncrementUploadCount() {
	g.ResetIfNeeded()

	g.MinuteCount++
	g.HourlyCount++
	g.DailyCount++
	g.LastUploadAt = &time.Time{}
	*g.LastUploadAt = time.Now()
}

/* Ban 封禁游客 */
func (g *GuestUploadLimit) Ban(reason string, duration time.Duration) {
	g.IsBanned = true
	g.BanReason = reason
	if duration > 0 {
		expiresAt := time.Now().Add(duration)
		g.BanExpiresAt = &expiresAt
	}
}

func (g *GuestUploadLimit) Unban() {
	g.IsBanned = false
	g.BanReason = ""
	g.BanExpiresAt = nil
}
