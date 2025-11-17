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

type GoogleOAuthService struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ProxyConfig  *ProxyConfig
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type GoogleUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func NewGoogleOAuthService(clientID, clientSecret, redirectURI string, proxyConfig *ProxyConfig) *GoogleOAuthService {
	return &GoogleOAuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		ProxyConfig:  proxyConfig,
	}
}

func (s *GoogleOAuthService) ExchangeCode(code string) (*GoogleTokenResponse, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

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
		return nil, fmt.Errorf("google API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp GoogleTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("未获取到 access_token")
	}

	return &tokenResp, nil
}

func (s *GoogleOAuthService) GetUserInfo(accessToken string) (*GoogleUserInfo, error) {
	userURL := "https://www.googleapis.com/oauth2/v2/userinfo"

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
		return nil, fmt.Errorf("google API returned error (status %d): %s", resp.StatusCode, string(body))
	}

	var userInfo GoogleUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if userInfo.Email == "" {
		return nil, fmt.Errorf("未获取到用户邮箱")
	}

	return &userInfo, nil
}

func (s *GoogleOAuthService) FindOrCreateUser(googleUser *GoogleUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var user models.User
	err := db.Where("google_id = ?", googleUser.ID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	username := googleUser.Email
	if googleUser.Name != "" {
		username = googleUser.Name
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

	googleID := googleUser.ID
	newUser := models.User{
		Username:  username,
		Email:     googleUser.Email,
		Avatar:    googleUser.Picture,
		GoogleID:  &googleID,
		PathAlias: utils.GenerateRandomString(16),
		Status:    1,
		Role:      3,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &newUser, nil
}
