package models

import (
	"gorm.io/gorm"
)

/* UserAccessControl 用户访问控制配置 */
type UserAccessControl struct {
	gorm.Model
	UserID uint `json:"user_id" gorm:"index;not null"` // 用户ID

	EnableRefererCheck bool   `json:"enable_referer_check"` // 是否启用防盗链检查
	AllowEmptyReferer  bool   `json:"allow_empty_referer"`  // 是否允许空Referer
	DomainWhitelist    string `json:"domain_whitelist"`     // 域名白名单,逗号分隔
	DomainBlacklist    string `json:"domain_blacklist"`     // 域名黑名单,逗号分隔

	EnableIPCheck bool   `json:"enable_ip_check"` // 是否启用IP检查
	IPWhitelist   string `json:"ip_whitelist"`    // IP白名单,逗号分隔
	IPBlacklist   string `json:"ip_blacklist"`    // IP黑名单,逗号分隔

	ControlMode string `json:"control_mode"` // 控制模式: strict(严格), moderate(适中), loose(宽松)

	BlockAction        string `json:"block_action"`         // 阻止行为: block(直接阻止), redirect(重定向), watermark(添加水印)
	RedirectURL        string `json:"redirect_url"`         // 重定向URL
	CustomErrorMessage string `json:"custom_error_message"` // 自定义错误信息
}

func (UserAccessControl) TableName() string {
	return "user_access_controls"
}

const (
	ControlModeStrict   = "strict"   // 严格模式
	ControlModeModerate = "moderate" // 适中模式
	ControlModeLoose    = "loose"    // 宽松模式
)

const (
	BlockActionBlock     = "block"     // 直接阻止
	BlockActionRedirect  = "redirect"  // 重定向
	BlockActionWatermark = "watermark" // 添加水印
	BlockActionThumbnail = "thumbnail" // 返回缩略图
)

/* ApplyPreset 应用预设配置 */
func (u *UserAccessControl) ApplyPreset(preset string) {
	switch preset {
	case ControlModeStrict:
		u.EnableRefererCheck = true
		u.AllowEmptyReferer = false
		u.EnableIPCheck = true
		u.BlockAction = BlockActionBlock
		u.ControlMode = ControlModeStrict
	case ControlModeModerate:
		u.EnableRefererCheck = true
		u.AllowEmptyReferer = true
		u.EnableIPCheck = false
		u.BlockAction = BlockActionRedirect
		u.ControlMode = ControlModeModerate
	case ControlModeLoose:
		u.EnableRefererCheck = false
		u.AllowEmptyReferer = true
		u.EnableIPCheck = false
		u.BlockAction = BlockActionWatermark
		u.ControlMode = ControlModeLoose
	}
}

/* CreateDefaultConfig 创建默认配置 */
func CreateDefaultConfig(userID uint) *UserAccessControl {
	config := &UserAccessControl{
		UserID:             userID,
		EnableRefererCheck: false,
		AllowEmptyReferer:  true,
		EnableIPCheck:      false,
		ControlMode:        ControlModeLoose,
		BlockAction:        BlockActionBlock,
	}
	return config
}
