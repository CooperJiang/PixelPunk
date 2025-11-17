package models

import (
	"pixelpunk/pkg/common"
	"time"

	"gorm.io/gorm"
)

/* Share 分享模型 */
type Share struct {
	ID        string          `gorm:"primarykey;size:32" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	UserID      uint   `gorm:"not null;index" json:"user_id"`
	ShareKey    string `gorm:"size:32;not null;uniqueIndex:idx_share_share_key" json:"share_key"`
	Name        string `gorm:"size:100" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Password    string `gorm:"size:100" json:"-"` // 访问密码(可选)，不返回到前端

	ExpiredDays int              `gorm:"default:0" json:"expired_days"` // 过期天数(0表示永不过期)
	ExpiredAt   *common.JSONTime `json:"expired_at"`                    // 计算得出的过期时间

	MaxViews     int `gorm:"default:0" json:"max_views"`     // 最大访问次数(0表示不限制)
	CurrentViews int `gorm:"default:0" json:"current_views"` // 当前访问次数

	Status int `gorm:"default:1;index" json:"status"` // 状态：1正常 2已过期 3已删除 4已禁用

	CollectVisitorInfo    bool `gorm:"default:false" json:"collect_visitor_info"`   // 是否收集访客信息
	NotificationOnAccess  bool `gorm:"default:false" json:"notification_on_access"` // 是否在被访问时通知创建者
	NotificationThreshold int  `gorm:"default:100" json:"notification_threshold"`   // 访问通知阈值，默认100次
}

func (Share) TableName() string {
	return "share"
}

func (s *Share) BeforeCreate(tx *gorm.DB) error {
	if s.ExpiredDays > 0 {
		expiredAt := time.Now().AddDate(0, 0, s.ExpiredDays)
		jsonTime := common.JSONTime(expiredAt)
		s.ExpiredAt = &jsonTime
	}

	if s.Status == 0 {
		s.Status = common.ShareStatusNormal
	}

	return nil
}

/* IsExpired 判断分享是否已过期 */
func (s *Share) IsExpired() bool {
	if s.ExpiredAt == nil {
		return false
	}
	return time.Now().After(time.Time(*s.ExpiredAt))
}

/* IsNormal 判断分享是否处于正常状态 */
func (s *Share) IsNormal() bool {
	return s.Status == common.ShareStatusNormal && !s.IsExpired()
}

/* IsAccessible 判断分享是否可访问 */
func (s *Share) IsAccessible() bool {
	if !s.IsNormal() {
		return false
	}

	if s.MaxViews > 0 && s.CurrentViews >= s.MaxViews {
		return false
	}

	return true
}

/* CanAccessWithPassword 校验密码是否正确 */
func (s *Share) CanAccessWithPassword(password string) bool {
	if s.Password == "" {
		return true
	}

	return s.Password == password
}
