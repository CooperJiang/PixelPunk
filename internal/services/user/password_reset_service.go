package user

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/email"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"
	"time"

	"gorm.io/gorm"
)

// 生成随机token
func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// SendPasswordResetEmail 发送密码重置邮件
func SendPasswordResetEmail(emailAddr string, ipAddress string) error {
	db := database.GetDB()

	user, err := FindUserByEmail(emailAddr)
	if err != nil {
		return errors.New(errors.CodeUserNotFound, "该邮箱尚未注册")
	}

	token, err := generateResetToken()
	if err != nil {
		return errors.New(errors.CodeInternal, "生成重置令牌失败")
	}

	expiresAt := time.Now().Add(1 * time.Hour)

	// 将旧的未使用token标记为已使用
	db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND is_used = ?", user.ID, false).
		Update("is_used", true)

	resetToken := models.PasswordResetToken{
		UserID:    user.ID,
		Token:     token,
		ExpiresAt: expiresAt,
		IsUsed:    false,
		IPAddress: ipAddress,
	}

	if err := db.Create(&resetToken).Error; err != nil {
		return errors.New(errors.CodeDBCreateFailed, "创建重置令牌失败")
	}

	websiteSettings, err := setting.GetSettingsByGroupAsMap("website")
	if err != nil {
		return errors.New(errors.CodeInternal, "获取网站设置失败")
	}

	baseURL, ok := websiteSettings.Settings["site_base_url"].(string)
	if !ok || baseURL == "" {
		return errors.New(errors.CodeInternal, "网站基础地址未配置，请联系管理员")
	}

	resetURL := fmt.Sprintf("%s/reset-password?token=%s", baseURL, token)

	subject := "密码重置"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
				<h2 style="color: #05d9e8;">密码重置请求</h2>
				<p>您好,</p>
				<p>我们收到了您的密码重置请求。请点击下面的按钮重置您的密码:</p>
				<div style="margin: 30px 0; text-align: center;">
					<a href="%s" style="display: inline-block; padding: 12px 30px; background-color: #05d9e8; color: #ffffff; text-decoration: none; border-radius: 5px; font-weight: bold;">重置密码</a>
				</div>
				<p>或者复制以下链接到浏览器中打开:</p>
				<p style="word-break: break-all; color: #666;"><a href="%s">%s</a></p>
				<p style="color: #999; font-size: 14px;">此链接将在1小时后失效。</p>
				<p style="color: #999; font-size: 14px;">如果您没有请求重置密码，请忽略此邮件。</p>
				<hr style="margin: 30px 0; border: none; border-top: 1px solid #eee;">
				<p style="color: #999; font-size: 12px;">此邮件由系统自动发送，请勿回复。</p>
			</div>
		</body>
		</html>
	`, resetURL, resetURL, resetURL)

	if err := email.SendMail(emailAddr, subject, body); err != nil {
		return errors.Wrap(err, errors.CodeEmailSendFailed, err.Error())
	}

	return nil
}

// VerifyResetToken 验证重置token
func VerifyResetToken(token string) (*models.User, error) {
	db := database.GetDB()

	var resetToken models.PasswordResetToken
	if err := db.Where("token = ?", token).First(&resetToken).Error; err != nil {
		return nil, errors.New(errors.CodeInvalidToken, "无效的重置令牌")
	}

	if !resetToken.IsValid() {
		if resetToken.IsUsed {
			return nil, errors.New(errors.CodeTokenUsed, "该重置令牌已被使用")
		}
		if resetToken.IsExpired() {
			return nil, errors.New(errors.CodeTokenExpired, "重置令牌已过期")
		}
	}

	var user models.User
	if err := db.First(&user, resetToken.UserID).Error; err != nil {
		return nil, errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	return &user, nil
}

// ResetPasswordWithToken 使用token重置密码
func ResetPasswordWithToken(token, newPassword string) error {
	db := database.GetDB()

	user, err := VerifyResetToken(token)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return errors.New(errors.CodeInternal, "密码加密失败")
	}

	// 使用 GORM Transaction 方法替代手动事务管理，确保 SQLite 兼容性
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&models.User{}).Where("id = ?", user.ID).Update("password", hashedPassword).Error; err != nil {
			return errors.New(errors.CodeDBUpdateFailed, "更新密码失败")
		}

		// 标记token为已使用
		now := time.Now()
		if err := tx.Model(&models.PasswordResetToken{}).
			Where("token = ?", token).
			Updates(map[string]interface{}{
				"is_used": true,
				"used_at": &now,
			}).Error; err != nil {
			return errors.New(errors.CodeDBUpdateFailed, "更新令牌状态失败")
		}

		return nil
	})
}

// CleanupExpiredTokens 清理过期的重置token
func CleanupExpiredTokens() error {
	db := database.GetDB()

	result := db.Where("expires_at < ?", time.Now()).Delete(&models.PasswordResetToken{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected > 0 {
	}

	return nil
}
