package dto

import "pixelpunk/pkg/common"

type AdminUserListQueryDTO struct {
	Page    int    `form:"page" json:"page"`       // 页码
	Size    int    `form:"size" json:"size"`       // 每页数量
	Keyword string `form:"keyword" json:"keyword"` // 关键字 (用户名或邮箱)
	Status  int    `form:"status" json:"status"`   // 用户状态, 0:所有状态, 1:正常, 2:禁用, 3:删除
	Role    int    `form:"role" json:"role"`       // 用户角色, 0:所有角色, 1:超级管理员, 2:管理员, 3:普通用户
}

type AdminUpdateUserDTO struct {
	ID       uint   `json:"id" binding:"required"`                    // 用户ID
	Username string `json:"username" binding:"required,min=2,max=50"` // 用户名
	Status   int    `json:"status" binding:"required,oneof=1 2 3"`    // 用户状态, 1:正常, 2:禁用, 3:删除
	Role     int    `json:"role" binding:"required,oneof=1 2 3"`      // 用户角色, 1:超级管理员, 2:管理员, 3:普通用户
}

func (d *AdminUpdateUserDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"ID.required":       "用户ID不能为空",
		"Username.required": "用户名不能为空",
		"Username.min":      "用户名长度不能小于2个字符",
		"Username.max":      "用户名长度不能超过50个字符",
		"Status.required":   "状态不能为空",
		"Status.oneof":      "状态值无效",
		"Role.required":     "角色不能为空",
		"Role.oneof":        "角色值无效",
	}
}

type AdminUserResponseDTO struct {
	ID             uint             `json:"id"`
	Username       string           `json:"username"`
	Email          string           `json:"email"`
	Avatar         string           `json:"avatar"`
	AvatarFullPath string           `json:"avatar_full_path"`
	Status         int              `json:"status"`
	Role           int              `json:"role"`
	StorageLimit   int64            `json:"storage_limit"`    // 存储空间限制（字节）
	BandwidthLimit int64            `json:"bandwidth_limit"`  // 带宽限制（字节）
	UsedStorage    int64            `json:"used_storage"`     // 已使用存储（字节）
	UsedBandwidth  int64            `json:"used_bandwidth"`   // 已使用带宽（字节）
	TotalFiles     int              `json:"total_files"`      // 总文件数量
	TotalViews     int64            `json:"total_views"`      // 总浏览数
	LastActivityAt *common.JSONTime `json:"last_activity_at"` // 最后操作时间
	LastActivityIP string           `json:"last_activity_ip"` // 最后操作IP
	CreatedAt      common.JSONTime  `json:"created_at"`
	UpdatedAt      common.JSONTime  `json:"updated_at"`
}

type AdminUserListResponseDTO struct {
	Total int64                  `json:"total"`
	List  []AdminUserResponseDTO `json:"list"`
}

type AdminCreateUserDTO struct {
	Username       string `json:"username" binding:"required,min=2,max=50"` // 用户名
	Email          string `json:"email" binding:"required,email,max=100"`   // 邮箱
	Password       string `json:"password" binding:"required,min=6,max=50"` // 密码
	Role           int    `json:"role" binding:"required,oneof=1 2 3"`      // 用户角色, 1:超级管理员, 2:管理员, 3:普通用户
	StorageLimit   int64  `json:"storage_limit,omitempty"`                  // 存储空间限制（字节），可选
	BandwidthLimit int64  `json:"bandwidth_limit,omitempty"`                // 带宽限制（字节），可选
}

func (d *AdminCreateUserDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"Username.required": "用户名不能为空",
		"Username.min":      "用户名长度不能小于2个字符",
		"Username.max":      "用户名长度不能超过50个字符",
		"Email.required":    "邮箱不能为空",
		"Email.email":       "邮箱格式不正确",
		"Email.max":         "邮箱长度不能超过100个字符",
		"Password.required": "密码不能为空",
		"Password.min":      "密码长度不能小于6个字符",
		"Password.max":      "密码长度不能超过50个字符",
		"Role.required":     "角色不能为空",
		"Role.oneof":        "角色值无效",
	}
}

type AdminUpdateUserStorageDTO struct {
	UserID         uint  `json:"user_id" binding:"required"`         // 用户ID
	StorageLimit   int64 `json:"storage_limit" binding:"required"`   // 存储空间限制（字节）
	BandwidthLimit int64 `json:"bandwidth_limit" binding:"required"` // 带宽限制（字节）
	CurrentUserID  uint  `json:"-"`                                  // 当前操作用户ID，不从请求体获取
}

func (d *AdminUpdateUserStorageDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserID.required":         "用户ID不能为空",
		"StorageLimit.required":   "存储空间限制不能为空",
		"BandwidthLimit.required": "带宽限制不能为空",
	}
}

type AdminResetUserPasswordResponseDTO struct {
	NewPassword string `json:"new_password"` // 新密码
}

type AdminSendUserEmailDTO struct {
	UserID  uint   `json:"user_id" binding:"required"` // 用户ID
	Subject string `json:"subject" binding:"required"` // 邮件主题
	Content string `json:"content" binding:"required"` // 邮件内容
}

func (d *AdminSendUserEmailDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserID.required":  "用户ID不能为空",
		"Subject.required": "邮件主题不能为空",
		"Content.required": "邮件内容不能为空",
	}
}

type AdminToggleUserStatusDTO struct {
	UserID uint `json:"user_id" binding:"required"`            // 用户ID
	Status int  `json:"status" binding:"required,oneof=1 2 3"` // 用户状态, 1:正常, 2:禁用, 3:删除
}

func (d *AdminToggleUserStatusDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserID.required": "用户ID不能为空",
		"Status.required": "状态不能为空",
		"Status.oneof":    "状态值无效",
	}
}

type AdminBatchOperateUsersDTO struct {
	UserIDs   []uint `json:"user_ids" binding:"required"`                                       // 用户ID列表
	Operation string `json:"operation" binding:"required,oneof=enable disable delete set_role"` // 操作类型
	Role      int    `json:"role,omitempty"`                                                    // 角色（仅在set_role时需要）
}

func (d *AdminBatchOperateUsersDTO) GetValidationMessages() map[string]string {
	return map[string]string{
		"UserIDs.required":   "用户ID列表不能为空",
		"Operation.required": "操作类型不能为空",
		"Operation.oneof":    "操作类型无效",
	}
}
