package health

import (
	"pixelpunk/pkg/email"
)

// EmailChecker 邮件服务健康检查器
type EmailChecker struct{}

// Name 返回检查项名称
func (c *EmailChecker) Name() string {
	return "email"
}

// Check 执行健康检查
func (c *EmailChecker) Check() (Status, map[string]interface{}) {
	if !email.IsMailEnabled() {
		return StatusUp, map[string]interface{}{
			"status":  "disabled",
			"message": "邮件服务未启用，请在后台管理系统中配置邮件设置",
		}
	}

	details := map[string]interface{}{
		"status":        "enabled",
		"message":       "邮件服务已启用，配置来源：数据库",
		"config_source": "database",
	}

	return StatusUp, details
}

// Type 返回检查类型
func (c *EmailChecker) Type() CheckType {
	return CheckTypeComplete // 不是基础检查
}

// init 注册邮件健康检查器
func init() {
	RegisterChecker(&EmailChecker{})
}
