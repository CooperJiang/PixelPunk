package oauth

import (
	"fmt"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/auth"
	oauthService "pixelpunk/internal/services/oauth"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

// OAuthHandler 通用 OAuth 处理函数
type OAuthHandler func(code string, proxyConfig *oauthService.ProxyConfig) (*models.User, error)

// handleOAuthLogin 通用 OAuth 登录处理逻辑
func handleOAuthLogin(c *gin.Context, code string, provider string, handler OAuthHandler, proxyConfig *oauthService.ProxyConfig) {
	user, err := handler(code, proxyConfig)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("%s 登录失败: %v", provider, err)))
		return
	}

	securitySettings, err := setting.GetSettingsByGroupAsMap("security")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "安全配置读取失败"))
		return
	}

	jwtSecret := ""
	if val, ok := securitySettings.Settings["jwt_secret"]; ok {
		if secretStr, ok := val.(string); ok {
			jwtSecret = secretStr
		}
	}
	if strings.TrimSpace(jwtSecret) == "" {
		errors.HandleError(c, errors.New(errors.CodeInternal, "JWT 密钥未配置"))
		return
	}

	expiresHours := 0
	if val, ok := securitySettings.Settings["login_expire_hours"]; ok {
		if hours, ok := val.(float64); ok && hours > 0 {
			expiresHours = int(hours)
		}
	}
	if expiresHours <= 0 {
		errors.HandleError(c, errors.New(errors.CodeInternal, "登录过期时间配置错误"))
		return
	}

	token, err := auth.GenerateToken(user.ID, user.Username, int(user.Role), jwtSecret, expiresHours)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "生成登录凭证失败"))
		return
	}

	userInfo := map[string]interface{}{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"avatar":   user.Avatar,
		"bio":      user.Bio,
		"website":  user.Website,
		"role":     user.Role,
		"status":   user.Status,
	}

	data := gin.H{
		"token":    token,
		"userInfo": userInfo,
		"email":    user.Email,
	}

	errors.ResponseSuccess(c, data, "登录成功")
}

func GithubLogin(c *gin.Context) {
	req, err := common.ValidateRequest[dto.GithubOAuthLoginDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	oauthConfig, err := setting.GetOAuthConfig()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取 OAuth 配置失败"))
		return
	}

	if !oauthConfig.Github.Enabled {
		errors.HandleError(c, errors.New(errors.CodeForbidden, "GitHub 登录功能未启用"))
		return
	}

	if oauthConfig.Github.ClientID == "" || oauthConfig.Github.ClientSecret == "" {
		errors.HandleError(c, errors.New(errors.CodeInternal, "GitHub OAuth 配置不完整"))
		return
	}

	var proxyConfig *oauthService.ProxyConfig
	if oauthConfig.Github.ProxyEnabled {
		proxyConfig = &oauthService.ProxyConfig{
			Enabled:  oauthConfig.Github.ProxyEnabled,
			Dynamic:  oauthConfig.Github.ProxyDynamic,
			APIURL:   oauthConfig.Github.ProxyAPIURL,
			Type:     oauthConfig.Github.ProxyType,
			Host:     oauthConfig.Github.ProxyHost,
			Port:     oauthConfig.Github.ProxyPort,
			Username: oauthConfig.Github.ProxyUsername,
			Password: oauthConfig.Github.ProxyPassword,
		}
	}

	handler := func(code string, proxyConfig *oauthService.ProxyConfig) (*models.User, error) {
		githubService := oauthService.NewGithubOAuthService(
			oauthConfig.Github.ClientID,
			oauthConfig.Github.ClientSecret,
			oauthConfig.Github.RedirectURI,
			proxyConfig,
		)

		tokenResp, err := githubService.ExchangeCode(code)
		if err != nil {
			return nil, fmt.Errorf("授权失败: %w", err)
		}

		githubUser, err := githubService.GetUserInfo(tokenResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %w", err)
		}

		return githubService.FindOrCreateUser(githubUser)
	}

	handleOAuthLogin(c, req.Code, "GitHub", handler, proxyConfig)
}

func GoogleLogin(c *gin.Context) {
	req, err := common.ValidateRequest[dto.GoogleOAuthLoginDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	oauthConfig, err := setting.GetOAuthConfig()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取 OAuth 配置失败"))
		return
	}

	if !oauthConfig.Google.Enabled {
		errors.HandleError(c, errors.New(errors.CodeForbidden, "Google 登录功能未启用"))
		return
	}

	if oauthConfig.Google.ClientID == "" || oauthConfig.Google.ClientSecret == "" {
		errors.HandleError(c, errors.New(errors.CodeInternal, "Google OAuth 配置不完整"))
		return
	}

	var proxyConfig *oauthService.ProxyConfig
	if oauthConfig.Google.ProxyEnabled {
		proxyConfig = &oauthService.ProxyConfig{
			Enabled:  oauthConfig.Google.ProxyEnabled,
			Dynamic:  oauthConfig.Google.ProxyDynamic,
			APIURL:   oauthConfig.Google.ProxyAPIURL,
			Type:     oauthConfig.Google.ProxyType,
			Host:     oauthConfig.Google.ProxyHost,
			Port:     oauthConfig.Google.ProxyPort,
			Username: oauthConfig.Google.ProxyUsername,
			Password: oauthConfig.Google.ProxyPassword,
		}
	}

	handler := func(code string, proxyConfig *oauthService.ProxyConfig) (*models.User, error) {
		googleService := oauthService.NewGoogleOAuthService(
			oauthConfig.Google.ClientID,
			oauthConfig.Google.ClientSecret,
			oauthConfig.Google.RedirectURI,
			proxyConfig,
		)

		tokenResp, err := googleService.ExchangeCode(code)
		if err != nil {
			return nil, fmt.Errorf("授权失败: %w", err)
		}

		googleUser, err := googleService.GetUserInfo(tokenResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %w", err)
		}

		return googleService.FindOrCreateUser(googleUser)
	}

	handleOAuthLogin(c, req.Code, "Google", handler, proxyConfig)
}

func LinuxdoLogin(c *gin.Context) {
	req, err := common.ValidateRequest[dto.LinuxdoOAuthLoginDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	oauthConfig, err := setting.GetOAuthConfig()
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取 OAuth 配置失败"))
		return
	}

	if !oauthConfig.Linuxdo.Enabled {
		errors.HandleError(c, errors.New(errors.CodeForbidden, "Linux DO 登录功能未启用"))
		return
	}

	if oauthConfig.Linuxdo.ClientID == "" || oauthConfig.Linuxdo.ClientSecret == "" {
		errors.HandleError(c, errors.New(errors.CodeInternal, "Linux DO OAuth 配置不完整"))
		return
	}

	handler := func(code string, proxyConfig *oauthService.ProxyConfig) (*models.User, error) {
		linuxdoService := oauthService.NewLinuxdoOAuthService(
			oauthConfig.Linuxdo.ClientID,
			oauthConfig.Linuxdo.ClientSecret,
			oauthConfig.Linuxdo.RedirectURI,
			nil,
		)

		tokenResp, err := linuxdoService.ExchangeCode(code)
		if err != nil {
			return nil, fmt.Errorf("授权失败: %w", err)
		}

		linuxdoUser, err := linuxdoService.GetUserInfo(tokenResp.AccessToken)
		if err != nil {
			return nil, fmt.Errorf("获取用户信息失败: %w", err)
		}

		return linuxdoService.FindOrCreateUser(linuxdoUser)
	}

	handleOAuthLogin(c, req.Code, "Linux DO", handler, nil)
}
