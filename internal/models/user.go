package models

import (
	"pixelpunk/pkg/common"

	"gorm.io/gorm"
)

/* User 用户模型 */
type User struct {
	ID        uint            `gorm:"primarykey" json:"id"`
	CreatedAt common.JSONTime `json:"created_at"`
	UpdatedAt common.JSONTime `json:"updated_at"`

	Username string `gorm:"size:50;not null;uniqueIndex:idx_user_username,sort:asc" json:"username"`
	Password string `gorm:"size:100" json:"-"`
	Email    string `gorm:"size:100;uniqueIndex:idx_user_email,sort:asc" json:"email"`
	Avatar   string `gorm:"size:255" json:"avatar"`
	Bio      string `gorm:"size:500" json:"bio"`
	Website  string `gorm:"size:255" json:"website"` // 个人网站
	Status   int    `gorm:"default:1" json:"status"` // 1:正常 2:禁用 3:删除
	Role     int    `gorm:"default:3" json:"role"`   // 1:超级管理员 2:管理员 3:普通用户

	PathAlias string  `gorm:"size:32;uniqueIndex:idx_user_path_alias,sort:asc" json:"path_alias"`
	GithubID  *int64  `gorm:"uniqueIndex:idx_user_github_id,sort:asc" json:"github_id"`
	GoogleID  *string `gorm:"size:50;uniqueIndex:idx_user_google_id,sort:asc" json:"google_id"`
	LinuxdoID *int64  `gorm:"uniqueIndex:idx_user_linuxdo_id,sort:asc" json:"linuxdo_id"`

	LastActivityAt *common.JSONTime `gorm:"column:last_activity_at" json:"last_activity_at"`
	LastActivityIP string           `gorm:"size:45;column:last_activity_ip" json:"last_activity_ip"` // 支持IPv6
}

func (User) TableName() string {
	return "user"
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Status == 0 {
		u.Status = common.UserStatusNormal
	}
	if u.Role == 0 {
		u.Role = common.UserRoleUser
	}
	return nil
}

/* IsAdmin 是否是管理员 */
func (u *User) IsAdmin() bool {
	return u.Role == common.UserRoleAdmin || u.Role == common.UserRoleSuperAdmin
}

/* IsSuperAdmin 是否是超级管理员 */
func (u *User) IsSuperAdmin() bool {
	return u.Role == common.UserRoleSuperAdmin
}

/* IsNormal 是否是正常状态 */
func (u *User) IsNormal() bool {
	return u.Status == common.UserStatusNormal
}

func (u *User) IsDisabled() bool {
	return u.Status == common.UserStatusDisabled
}

/* IsDeleted 是否被删除 */
func (u *User) IsDeleted() bool {
	return u.Status == common.UserStatusDeleted
}
