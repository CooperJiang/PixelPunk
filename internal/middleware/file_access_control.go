package middleware

import (
	"fmt"
	"net/http"
	"net/url"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/access_control"
	"pixelpunk/internal/services/auth"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/internal/services/setting"
	"pixelpunk/internal/services/share"
	"pixelpunk/internal/services/stats"
	"pixelpunk/internal/services/user"
	"pixelpunk/pkg/assets"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func FileAccessControlMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fileObj, exists := c.Get("file_info")
		if !exists {
			logger.Error("[ACCESS_CONTROL] 未找到文件信息，请确保 FileInfoExtractorMiddleware 在前")
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}

		file, ok := fileObj.(models.File)
		if !ok {
			logger.Error("[ACCESS_CONTROL] 文件信息类型转换失败")
			assets.ServeDefaultFile(c, assets.FileTypeNotFound)
			return
		}

		isThumbObj, _ := c.Get("isThumb")
		isThumb, _ := isThumbObj.(bool)

		if isSpecialAccessScenario(c) {
			if !isThumb {
				go updateFileStats(file.ID, file.UserID, file.Size)
			}
			c.Next()
			return
		}

		isInternalRequest := isFromConfiguredBaseUrl(c)

		if !isThumb {
			go updateFileStats(file.ID, file.UserID, file.Size)
		}

		if isInternalRequest && (file.AccessLevel == "public" || file.AccessLevel == "private") {
			if file.Status == "pending_review" {
				assets.ServeDefaultFile(c, assets.FileTypeReview)
				return
			}
			c.Next()
			return
		}

		shareKey := c.Query("share")
		accessToken := c.Query("access_token")
		if shareKey != "" {
			if accessToken != "" {
				valid, _ := share.ValidateAccessToken(shareKey, accessToken)
				if valid {
					c.Next()
					return
				}
				assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
				return
			}
			if verifyShareAccess(c, shareKey, file.ID) {
				c.Next()
				return
			}
			assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
			return
		}

		if !handleUserAccessControl(c, file) {
			return
		}

		handleFileAccessLevel(c, file, isInternalRequest)
	}
}

func verifyShareAccess(c *gin.Context, shareKey string, fileID string) bool {
	var share models.Share
	if err := database.DB.Where("share_key = ? AND status = ?", shareKey, common.ShareStatusNormal).First(&share).Error; err != nil {
		return false
	}

	if share.ExpiredAt != nil && time.Now().After(time.Time(*share.ExpiredAt)) {
		database.DB.Model(&share).Update("status", common.ShareStatusExpired)
		return false
	}

	if share.MaxViews > 0 && share.CurrentViews >= share.MaxViews {
		return false
	}

	if share.Password != "" {
		return false
	}

	var count int64
	database.DB.Model(&models.ShareItem{}).
		Where("share_id = ? AND item_type = ? AND item_id = ?", share.ID, common.ShareItemTypeFile, fileID).
		Count(&count)
	if count > 0 {
		return true
	}

	var targetImage models.File
	if err := database.DB.Where("id = ?", fileID).First(&targetImage).Error; err != nil {
		return false
	}

	var folderShares []models.ShareItem
	database.DB.Where("share_id = ? AND item_type = ?", share.ID, common.ShareItemTypeFolder).Find(&folderShares)

	for _, folderShare := range folderShares {
		if targetImage.FolderID == folderShare.ItemID {
			return true
		}

		if targetImage.FolderID != "" && targetImage.FolderID != "0" {
			currentFolderID := targetImage.FolderID
			for currentFolderID != "" && currentFolderID != "0" {
				var folder models.Folder
				if err := database.DB.Where("id = ?", currentFolderID).First(&folder).Error; err != nil {
					break // 找不到文件夹，退出循环
				}

				if folder.ID == folderShare.ItemID {
					return true
				}

				currentFolderID = folder.ParentID
			}
		}
	}

	return false
}

type FileAccessConfig struct {
	PublicCacheMaxAge  int // 公开文件缓存时间（秒）
	PrivateCacheMaxAge int // 私有文件缓存时间（秒）
}

var DefaultFileAccessConfig = FileAccessConfig{
	PublicCacheMaxAge:  60 * 60 * 24 * 60, // 60天 (5,184,000秒)
	PrivateCacheMaxAge: 60 * 60 * 24,      // 1天 (86,400秒)
}

var CurrentFileAccessConfig = DefaultFileAccessConfig

func SetFileAccessConfig(config FileAccessConfig) {
	CurrentFileAccessConfig = config
}

func updateFileStats(fileID string, userID uint, size int64) {
	filesvc.UpdateViews(fileID)

	filesvc.UpdateBandwidth(fileID, size)

	user.UpdateViewsUsage(userID, 1)

	user.UpdateBandwidthUsage(userID, size)

	stats.GetStatsAdapter().RecordFileViewed(size)
}

func handleFileAccessLevel(c *gin.Context, file models.File, isInternalRequest bool) bool {
	switch file.AccessLevel {
	case "public":
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", CurrentFileAccessConfig.PublicCacheMaxAge))

		if file.Status == "pending_review" {
			assets.ServeDefaultFile(c, assets.FileTypeReview)
			return true
		}

		c.Next()
		return true

	case "protected":
		return handleProtectedAccess(c, file, isInternalRequest)

	case "private":
		return handlePrivateAccess(c, file)

	default:
		assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
		return false
	}
}

func handleProtectedAccess(c *gin.Context, file models.File, isInternalRequest bool) bool {
	if !isInternalRequest {
		assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
		return false
	}

	if CanUserAccessProtectedFile(c, file.UserID) {
		c.Header("Cache-Control", "private, max-age=3600")

		if file.Status == "pending_review" {
			assets.ServeDefaultFile(c, assets.FileTypeReview)
			return true
		}

		c.Next()
		return true
	}

	assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
	return false
}

func handlePrivateAccess(c *gin.Context, file models.File) bool {
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", CurrentFileAccessConfig.PublicCacheMaxAge))

	if file.Status == "pending_review" {
		assets.ServeDefaultFile(c, assets.FileTypeReview)
		return true
	}

	c.Next()
	return true
}

func handleUserAccessControl(c *gin.Context, file models.File) bool {
	config, err := access_control.GetUserAccessControl(file.UserID)
	if err != nil {
		return true // 获取配置失败，默认允许访问
	}

	if !config.EnableRefererCheck && !config.EnableIPCheck {
		return true
	}

	ip := c.ClientIP()
	referer := c.GetHeader("Referer")
	domain := extractDomainFromReferer(referer)

	if config.EnableIPCheck {
		isIPAllowed := access_control.CheckUserIP(config, ip, isIPInList)
		if !isIPAllowed {
			handleBlockAction(c, config, ip, domain)
			return false
		}
	}

	if config.EnableRefererCheck {
		if referer == "" && !config.AllowEmptyReferer {
			handleBlockAction(c, config, ip, domain)
			return false
		}

		if referer != "" {
			isDomainAllowed := access_control.CheckUserDomain(config, domain, referer, isDomainInList)
			if !isDomainAllowed {
				handleBlockAction(c, config, ip, domain)
				return false
			}
		}
	}

	return true
}

func handleBlockAction(c *gin.Context, config *models.UserAccessControl, ip, domain string) {
	switch config.BlockAction {
	case models.BlockActionRedirect:
		if config.RedirectURL != "" {
			c.Redirect(http.StatusTemporaryRedirect, config.RedirectURL)
			c.Abort()
			return
		}
		fallthrough
	case models.BlockActionBlock:
		assets.ServeDefaultFile(c, assets.FileTypeUnauthorized)
	case models.BlockActionWatermark:
		c.Set("watermark", true)
		c.Next()
	case models.BlockActionThumbnail:
		c.Set("forceThumbnail", true)
		c.Next()
	}
}

func isSpecialAccessScenario(c *gin.Context) bool {
	return isAdminRequest(c)
}

func isAdminRequest(c *gin.Context) bool {
	timeParam := c.Query("t")
	signatureParam := c.Query("s")

	if timeParam != "" && signatureParam != "" {
		signer := utils.GetURLSigner()
		if signer.VerifyFileURL(c.Request.URL.Path, timeParam, signatureParam) {
			return true
		}
	}

	authHeader := c.GetHeader("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		jwtSecret := getJWTSecret()
		claims, err := auth.ParseToken(tokenString, jwtSecret)
		if err == nil && (claims.Role == common.UserRoleAdmin || claims.Role == common.UserRoleSuperAdmin) {
			return true
		}
	}

	return false
}

func getJWTSecret() string {
	jwtSecret := ""
	securitySettings, err := setting.GetSettingsByGroupAsMap("security")
	if err == nil {
		if val, ok := securitySettings.Settings["jwt_secret"]; ok {
			if secretStr, ok := val.(string); ok {
				jwtSecret = secretStr
			}
		}
	}
	return jwtSecret
}

func isFromConfiguredBaseUrl(c *gin.Context) bool {
	baseUrl := utils.GetBaseUrl()
	referer := c.GetHeader("Referer")

	if baseUrl == "" || referer == "" {
		return false
	}

	if isLocalDevelopment(baseUrl, referer) {
		return true
	}

	baseUrlObj, err1 := url.Parse(baseUrl)
	refererObj, err2 := url.Parse(referer)
	if err1 != nil || err2 != nil {
		return false
	}

	baseHost := baseUrlObj.Hostname()
	refererHost := refererObj.Hostname()
	basePort := baseUrlObj.Port()
	refererPort := refererObj.Port()

	if baseHost == refererHost {
		if (basePort == "" && refererPort == "") || basePort == refererPort {
			return true
		}

		if basePort == "" {
			if baseUrlObj.Scheme == "https" && refererPort == "443" {
				return true
			}
			if baseUrlObj.Scheme == "http" && refererPort == "80" {
				return true
			}
		}

		if refererPort == "" {
			if refererObj.Scheme == "https" && basePort == "443" {
				return true
			}
			if refererObj.Scheme == "http" && basePort == "80" {
				return true
			}
		}
	}

	return false
}

func isLocalDevelopment(baseUrl, referer string) bool {
	baseUrlObj, err1 := url.Parse(baseUrl)
	refererObj, err2 := url.Parse(referer)

	if err1 != nil || err2 != nil {
		return false
	}

	baseHost := baseUrlObj.Hostname()
	refererHost := refererObj.Hostname()

	isLocalBase := baseHost == "localhost" || baseHost == "127.0.0.1"
	isLocalReferer := refererHost == "localhost" || refererHost == "127.0.0.1"

	return isLocalBase && isLocalReferer
}
