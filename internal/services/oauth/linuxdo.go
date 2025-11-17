package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/utils"
	"strings"
)

type LinuxdoOAuthService struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ProxyConfig  *ProxyConfig
}

type LinuxdoTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type LinuxdoUserInfo struct {
	ID             int64  `json:"id"`
	Username       string `json:"username"`
	Name           string `json:"name"`
	AvatarTemplate string `json:"avatar_template"`
	Active         bool   `json:"active"`
	TrustLevel     int    `json:"trust_level"`
	Silenced       bool   `json:"silenced"`
}

func NewLinuxdoOAuthService(clientID, clientSecret, redirectURI string, proxyConfig *ProxyConfig) *LinuxdoOAuthService {
	return &LinuxdoOAuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		ProxyConfig:  proxyConfig,
	}
}

func (s *LinuxdoOAuthService) ExchangeCode(code string) (*LinuxdoTokenResponse, error) {
	tokenURL := "https://connect.linux.do/oauth2/token"

	data := url.Values{}
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", s.RedirectURI)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := getHTTPClient(s.ProxyConfig)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linuxdo API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp LinuxdoTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("未获取到 access_token")
	}

	return &tokenResp, nil
}

func (s *LinuxdoOAuthService) GetUserInfo(accessToken string) (*LinuxdoUserInfo, error) {
	userURL := "https://connect.linux.do/api/user"

	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	client := getHTTPClient(s.ProxyConfig)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linuxdo API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var userInfo LinuxdoUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if userInfo.ID == 0 {
		return nil, fmt.Errorf("未获取到用户 ID")
	}

	return &userInfo, nil
}

func (s *LinuxdoOAuthService) FindOrCreateUser(linuxdoUser *LinuxdoUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var user models.User
	err := db.Where("linuxdo_id = ?", linuxdoUser.ID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	username := linuxdoUser.Username
	if username == "" {
		username = linuxdoUser.Name
	}

	existingUser := models.User{}
	if db.Where("username = ?", username).First(&existingUser).Error == nil {
		for i := 1; i < 100; i++ {
			newUsername := fmt.Sprintf("%s%d", username, i)
			if db.Where("username = ?", newUsername).First(&existingUser).Error != nil {
				username = newUsername
				break
			}
		}
	}

	avatar := ""
	if linuxdoUser.AvatarTemplate != "" {
		avatar = strings.ReplaceAll(linuxdoUser.AvatarTemplate, "{size}", "120")
		if strings.HasPrefix(avatar, "/") {
			avatar = "https://linux.do" + avatar
		}
	}

	linuxdoID := linuxdoUser.ID
	newUser := models.User{
		Username:  username,
		Avatar:    avatar,
		LinuxdoID: &linuxdoID,
		PathAlias: utils.GenerateRandomString(16),
		Status:    1,
		Role:      3,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &newUser, nil
}
