package email

import (
	"crypto/tls"
	"errors"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/hooks"
	"strings"
	"sync"

	"gopkg.in/gomail.v2"
)

var (
	mailer *gomail.Dialer
	mutex  sync.RWMutex
)

// Init 初始化邮件服务
func Init() {
	// 注册设置更新钩子，当邮件设置更新时自动刷新配置
	hooks.RegisterSettingUpdateHook("mail", func(group string) error {
		return RefreshMailer()
	})

	_ = RefreshMailer()
}

// GetMailerFromSettings 从设置中获取邮件发送器
func GetMailerFromSettings() (*gomail.Dialer, map[string]interface{}, error) {
	mailSettings, err := setting.GetSettingsByGroupAsMap("mail")
	if err != nil {
		return nil, nil, err
	}

	host, hostOk := mailSettings.Settings["smtp_host"].(string)
	username, usernameOk := mailSettings.Settings["smtp_username"].(string)
	password, passwordOk := mailSettings.Settings["smtp_password"].(string)

	// 获取端口，并转换为整数
	port := 25 // 默认端口
	portValue, portOk := mailSettings.Settings["smtp_port"].(float64)
	if portOk {
		port = int(portValue)
	}

	if !hostOk || !usernameOk || !passwordOk || host == "" {
		return nil, mailSettings.Settings, errors.New("邮件设置不完整")
	}

	dialer := gomail.NewDialer(host, port, username, password)

	encryption, _ := mailSettings.Settings["smtp_encryption"].(string)
	if encryption == "ssl" || encryption == "tls" {
		dialer.SSL = true
	} else if encryption == "starttls" {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return dialer, mailSettings.Settings, nil
}

// RefreshMailer 刷新邮件发送器
func RefreshMailer() error {
	newMailer, _, err := GetMailerFromSettings()
	if err != nil {
		return err
	}

	mutex.Lock()
	mailer = newMailer
	mutex.Unlock()

	return nil
}

// IsMailEnabled 检查邮件服务是否可用
func IsMailEnabled() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return mailer != nil
}

// SendMail 发送邮件
func SendMail(to, subject, body string) error {
	if !IsMailEnabled() {
		return errors.New("邮件服务未启用，请在后台管理系统中配置邮件设置")
	}

	// 从数据库获取发件人信息
	mailSettings, err := setting.GetSettingsByGroupAsMap("mail")
	if err != nil {
		return errors.New("无法获取邮件配置：" + err.Error())
	}

	fromAddress, _ := mailSettings.Settings["smtp_from_address"].(string)
	fromName, _ := mailSettings.Settings["smtp_from_name"].(string)

	if fromAddress == "" {
		return errors.New("邮件配置不完整：请在后台管理系统中配置发件人地址")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromAddress, fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	replyTo, _ := mailSettings.Settings["smtp_reply_to"].(string)
	if replyTo != "" {
		m.SetHeader("Reply-To", replyTo)
	}

	mutex.RLock()
	currentMailer := mailer
	mutex.RUnlock()

	if err := currentMailer.DialAndSend(m); err != nil {
		// 将常见的英文错误翻译成中文
		return translateMailError(err)
	}

	return nil
}

// translateMailError 将邮件发送错误翻译成中文
func translateMailError(err error) error {
	errMsg := err.Error()

	// SMTP服务器地址错误
	if strings.Contains(errMsg, "no such host") || strings.Contains(errMsg, "lookup") {
		return errors.New("SMTP服务器地址无效,请检查邮件配置")
	}

	if strings.Contains(errMsg, "connection refused") {
		return errors.New("SMTP服务器拒绝连接,请检查端口配置")
	}

	if strings.Contains(errMsg, "authentication failed") || strings.Contains(errMsg, "Invalid credentials") || strings.Contains(errMsg, "535") {
		return errors.New("SMTP认证失败,请检查用户名和密码")
	}

	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "i/o timeout") {
		return errors.New("SMTP服务器连接超时")
	}

	// TLS/SSL错误
	if strings.Contains(errMsg, "tls") || strings.Contains(errMsg, "certificate") {
		return errors.New("TLS/SSL连接失败,请检查加密方式")
	}

	// 其他错误,添加中文前缀
	return errors.New("邮件发送失败: " + errMsg)
}

// TestMailSettings 测试邮件设置是否正确
func TestMailSettings(host string, port int, username, password, encryption, fromAddress, fromName, testEmail string) error {
	dialer := gomail.NewDialer(host, port, username, password)

	if encryption == "ssl" || encryption == "tls" {
		dialer.SSL = true
	} else if encryption == "starttls" {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromAddress, fromName))
	m.SetHeader("To", testEmail)
	m.SetHeader("Subject", "测试邮件")
	m.SetBody("text/html", "<h1>这是一封测试邮件</h1><p>如果您收到这封邮件，说明邮件服务配置正确。</p>")

	if err := dialer.DialAndSend(m); err != nil {
		return translateMailError(err)
	}

	return nil
}
