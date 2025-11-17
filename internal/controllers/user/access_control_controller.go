package user

import (
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/access_control"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

type AccessControlDTO struct {
	EnableRefererCheck bool   `json:"enable_referer_check"` // 是否启用防盗链检查
	AllowEmptyReferer  bool   `json:"allow_empty_referer"`  // 是否允许空Referer
	DomainWhitelist    string `json:"domain_whitelist"`     // 域名白名单,逗号分隔
	DomainBlacklist    string `json:"domain_blacklist"`     // 域名黑名单,逗号分隔
	EnableIPCheck      bool   `json:"enable_ip_check"`      // 是否启用IP检查
	IPWhitelist        string `json:"ip_whitelist"`         // IP白名单,逗号分隔
	IPBlacklist        string `json:"ip_blacklist"`         // IP黑名单,逗号分隔
	ControlMode        string `json:"control_mode"`         // 控制模式
	BlockAction        string `json:"block_action"`         // 阻止行为
	RedirectURL        string `json:"redirect_url"`         // 重定向URL
	CustomErrorMessage string `json:"custom_error_message"` // 自定义错误信息
}

type NewAccessControlDTO struct {
	Enabled            bool   `json:"enabled"`              // 是否启用访问控制
	IPMode             string `json:"ip_mode"`              // IP模式：whitelist(白名单)或blacklist(黑名单)
	DomainMode         string `json:"domain_mode"`          // 域名模式：whitelist(白名单)或blacklist(黑名单)
	RestrictionMode    string `json:"restriction_mode"`     // 限制模式：strict(严格)、moderate(适中)、loose(宽松)
	BlockAction        string `json:"block_action"`         // 阻止行为：block(直接阻止)、redirect(重定向)、watermark(添加水印)
	RedirectURL        string `json:"redirect_url"`         // 重定向URL
	CustomErrorMessage string `json:"custom_error_message"` // 自定义错误信息
	IPList             string `json:"ip_list"`              // IP列表，逗号分隔
	DomainList         string `json:"domain_list"`          // 域名列表，逗号分隔
}

func GetUserAccessControl(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	config, err := access_control.GetUserAccessControl(userID)
	if err != nil {
		errCode := errors.CodeInternal
		if e, ok := err.(*errors.Error); ok {
			errCode = e.Code
		}
		errObj := &errors.Error{
			Code:    errCode,
			Message: "获取访问控制配置失败",
			Detail:  err.Error(),
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	dto := NewAccessControlDTO{
		Enabled:            config.EnableRefererCheck || config.EnableIPCheck,
		RestrictionMode:    config.ControlMode,
		BlockAction:        config.BlockAction,
		RedirectURL:        config.RedirectURL,
		CustomErrorMessage: config.CustomErrorMessage,
	}

	if config.IPWhitelist != "" {
		dto.IPMode = "whitelist"
		dto.IPList = config.IPWhitelist
	} else {
		dto.IPMode = "blacklist"
		dto.IPList = config.IPBlacklist
	}

	if config.DomainWhitelist != "" {
		dto.DomainMode = "whitelist"
		dto.DomainList = config.DomainWhitelist
	} else {
		dto.DomainMode = "blacklist"
		dto.DomainList = config.DomainBlacklist
	}

	errors.ResponseSuccess(c, dto, "获取访问控制配置成功")
}

func ApplyAccessControlPreset(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	preset := c.Param("preset")
	if preset != models.ControlModeStrict &&
		preset != models.ControlModeModerate &&
		preset != models.ControlModeLoose {
		errObj := &errors.Error{
			Code:    errors.CodeInvalidParameter,
			Message: "参数验证失败",
			Detail:  "无效的预设模式",
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	if err := access_control.ApplyControlPreset(userID, preset); err != nil {
		errCode := errors.CodeInternal
		if e, ok := err.(*errors.Error); ok {
			errCode = e.Code
		}
		errObj := &errors.Error{
			Code:    errCode,
			Message: "应用预设访问控制配置失败",
			Detail:  err.Error(),
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	errors.ResponseSuccess(c, nil, "应用预设访问控制配置成功")
}

func ResetUserAccessControl(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	if err := access_control.DeleteUserAccessControl(userID); err != nil {
		errCode := errors.CodeInternal
		if e, ok := err.(*errors.Error); ok {
			errCode = e.Code
		}
		errObj := &errors.Error{
			Code:    errCode,
			Message: "删除访问控制配置失败",
			Detail:  err.Error(),
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	config := models.CreateDefaultConfig(userID)
	config.EnableRefererCheck = false
	config.EnableIPCheck = false

	if err := access_control.UpdateUserAccessControl(config); err != nil {
		errCode := errors.CodeInternal
		if e, ok := err.(*errors.Error); ok {
			errCode = e.Code
		}
		errObj := &errors.Error{
			Code:    errCode,
			Message: "创建默认访问控制配置失败",
			Detail:  err.Error(),
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	errors.ResponseSuccess(c, nil, "重置访问控制配置成功")
}

func CreateOrUpdateUserAccessControl(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var dto NewAccessControlDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		errors.ResponseError(c, errors.ErrInvalidParams, err.Error())
		return
	}

	validRestrictionModes := map[string]bool{
		models.ControlModeStrict:   true,
		models.ControlModeModerate: true,
		models.ControlModeLoose:    true,
		"":                         true, // 允许空值
	}

	if !validRestrictionModes[dto.RestrictionMode] {
		errObj := &errors.Error{
			Code:    errors.CodeInvalidParameter,
			Message: "参数验证失败",
			Detail:  "无效的限制模式",
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	validBlockActions := map[string]bool{
		models.BlockActionBlock:     true,
		models.BlockActionRedirect:  true,
		models.BlockActionWatermark: true,
		models.BlockActionThumbnail: true,
		"":                          true, // 允许空值
	}

	if !validBlockActions[dto.BlockAction] {
		errObj := &errors.Error{
			Code:    errors.CodeInvalidParameter,
			Message: "参数验证失败",
			Detail:  "无效的阻止行为",
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	validModes := map[string]bool{
		"whitelist": true,
		"blacklist": true,
		"":          true, // 允许空值
	}

	if !validModes[dto.IPMode] {
		errObj := &errors.Error{
			Code:    errors.CodeInvalidParameter,
			Message: "参数验证失败",
			Detail:  "无效的IP模式，应为whitelist或blacklist",
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	if !validModes[dto.DomainMode] {
		errObj := &errors.Error{
			Code:    errors.CodeInvalidParameter,
			Message: "参数验证失败",
			Detail:  "无效的域名模式，应为whitelist或blacklist",
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	config := &models.UserAccessControl{
		UserID:             userID,
		EnableRefererCheck: dto.Enabled && dto.DomainList != "",
		AllowEmptyReferer:  dto.RestrictionMode != models.ControlModeStrict, // 非严格模式允许空Referer
		EnableIPCheck:      dto.Enabled && dto.IPList != "",
		ControlMode:        dto.RestrictionMode,
		BlockAction:        dto.BlockAction,
		RedirectURL:        dto.RedirectURL,
		CustomErrorMessage: dto.CustomErrorMessage,
	}

	if dto.IPMode == "whitelist" {
		config.IPWhitelist = dto.IPList
		config.IPBlacklist = ""
	} else {
		config.IPBlacklist = dto.IPList
		config.IPWhitelist = ""
	}

	if dto.DomainMode == "whitelist" {
		config.DomainWhitelist = dto.DomainList
		config.DomainBlacklist = ""
	} else {
		config.DomainBlacklist = dto.DomainList
		config.DomainWhitelist = ""
	}

	if err := access_control.UpdateUserAccessControl(config); err != nil {
		errCode := errors.CodeInternal
		if e, ok := err.(*errors.Error); ok {
			errCode = e.Code
		}
		errObj := &errors.Error{
			Code:    errCode,
			Message: "创建或更新访问控制配置失败",
			Detail:  err.Error(),
		}
		errors.ResponseError(c, errObj, "")
		return
	}

	errors.ResponseSuccess(c, nil, "访问控制配置已创建或更新")
}
