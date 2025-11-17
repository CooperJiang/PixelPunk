package file

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
)

/* TemporaryLinkResult 生成临时链接的结果 */
type TemporaryLinkResult struct {
	FileID    string
	URL       string
	ExpiresIn int       // 分钟
	ExpiresAt time.Time // 绝对时间
}

/* GenerateTemporaryLink 为用户文件生成带签名的临时访问链接 */
func GenerateTemporaryLink(userID uint, fileID string, expireMinutes int) (*TemporaryLinkResult, error) {
	if fileID == "" {
		return nil, errors.New(errors.CodeInvalidParameter, "文件ID不能为空")
	}
	if expireMinutes <= 0 {
		expireMinutes = 5
	}
	if expireMinutes > 60 {
		expireMinutes = 60
	}

	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).
		Where("status <> ?", "pending_deletion").
		First(&file).Error; err != nil {
		return nil, errors.New(errors.CodeFileNotFound, "文件不存在或无权访问")
	}

	expireAt := time.Now().Add(time.Minute * time.Duration(expireMinutes))
	claims := jwt.MapClaims{
		"sub": fileID,
		"iss": "cyber_image_server",
		"iat": time.Now().Unix(),
		"exp": expireAt.Unix(),
		"uid": userID,
		"fn":  file.FileName,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtSecret := "pixel_punk_jwt_secret_key"
	if securitySettings, err := setting.GetSettingsByGroupAsMap("security"); err == nil {
		if val, ok := securitySettings.Settings["jwt_secret"]; ok {
			if secretStr, ok := val.(string); ok && secretStr != "" {
				jwtSecret = secretStr
			}
		}
	}
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, errors.Wrap(err, errors.CodeInternal, "生成访问令牌失败")
	}

	baseURL := file.URL
	if baseURL == "" {
		if file.ThumbURL != "" {
			baseURL = file.ThumbURL
		} else {
			baseURL = fmt.Sprintf("/f/%s", file.ID)
		}
	}
	url := fmt.Sprintf("%s?token=%s", baseURL, tokenString)

	return &TemporaryLinkResult{FileID: fileID, URL: url, ExpiresIn: expireMinutes, ExpiresAt: expireAt}, nil
}
