package oauth

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"pixelpunk/internal/models"
	settingService "pixelpunk/internal/services/setting"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/utils"
	"time"
)

type ProxyConfig struct {
	Enabled  bool
	Dynamic  bool
	APIURL   string
	Type     string
	Host     string
	Port     string
	Username string
	Password string
}

func getHTTPClient(proxyConfig *ProxyConfig) *http.Client {
	transport := &http.Transport{}

	if proxyConfig != nil && proxyConfig.Enabled && proxyConfig.Dynamic && proxyConfig.APIURL != "" {
		dynamicProxy, err := settingService.FetchDynamicProxy(proxyConfig.APIURL, proxyConfig.Username, proxyConfig.Password)
		if err == nil {
			var proxyURL string
			if dynamicProxy.Username != "" && dynamicProxy.Password != "" {
				proxyURL = fmt.Sprintf("%s://%s:%s@%s:%d",
					dynamicProxy.Type,
					url.QueryEscape(dynamicProxy.Username),
					url.QueryEscape(dynamicProxy.Password),
					dynamicProxy.Host,
					dynamicProxy.Port,
				)
			} else {
				proxyURL = fmt.Sprintf("%s://%s:%d",
					dynamicProxy.Type,
					dynamicProxy.Host,
					dynamicProxy.Port,
				)
			}

			if proxy, err := url.Parse(proxyURL); err == nil {
				transport.Proxy = http.ProxyURL(proxy)
			}
		}
	} else if proxyConfig != nil && proxyConfig.Enabled && proxyConfig.Host != "" && proxyConfig.Port != "" {
		var proxyURL string
		if proxyConfig.Username != "" && proxyConfig.Password != "" {
			proxyURL = fmt.Sprintf("%s://%s:%s@%s:%s",
				proxyConfig.Type,
				url.QueryEscape(proxyConfig.Username),
				url.QueryEscape(proxyConfig.Password),
				proxyConfig.Host,
				proxyConfig.Port,
			)
		} else {
			proxyURL = fmt.Sprintf("%s://%s:%s",
				proxyConfig.Type,
				proxyConfig.Host,
				proxyConfig.Port,
			)
		}

		if proxy, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(proxy)
		}
	} else {
		if proxyURL := os.Getenv("HTTP_PROXY"); proxyURL != "" {
			if proxy, err := url.Parse(proxyURL); err == nil {
				transport.Proxy = http.ProxyURL(proxy)
			}
		} else if proxyURL := os.Getenv("HTTPS_PROXY"); proxyURL != "" {
			if proxy, err := url.Parse(proxyURL); err == nil {
				transport.Proxy = http.ProxyURL(proxy)
			}
		}
	}

	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

type GithubOAuthService struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	ProxyConfig  *ProxyConfig
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUserInfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	Bio       string `json:"bio"`
	Blog      string `json:"blog"`
}

func NewGithubOAuthService(clientID, clientSecret, redirectURI string, proxyConfig *ProxyConfig) *GithubOAuthService {
	return &GithubOAuthService{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		ProxyConfig:  proxyConfig,
	}
}

func (s *GithubOAuthService) ExchangeCode(code string) (*GithubTokenResponse, error) {
	tokenURL := "https://github.com/login/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", s.ClientID)
	data.Set("client_secret", s.ClientSecret)
	data.Set("code", code)

	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.URL.RawQuery = data.Encode()
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
		return nil, fmt.Errorf("GitHub 返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var tokenResp GithubTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("未获取到 access_token")
	}

	return &tokenResp, nil
}

func (s *GithubOAuthService) GetUserInfo(accessToken string) (*GithubUserInfo, error) {
	userURL := "https://api.github.com/user"

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
		return nil, fmt.Errorf("GitHub 返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var userInfo GithubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if userInfo.Email == "" {
		emails, err := s.GetUserEmails(accessToken)
		if err == nil && len(emails) > 0 {
			for _, email := range emails {
				if email.Primary && email.Verified {
					userInfo.Email = email.Email
					break
				}
			}
			if userInfo.Email == "" && len(emails) > 0 {
				userInfo.Email = emails[0].Email
			}
		}
	}

	return &userInfo, nil
}

type GithubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (s *GithubOAuthService) GetUserEmails(accessToken string) ([]GithubEmail, error) {
	emailURL := "https://api.github.com/user/emails"

	req, err := http.NewRequest("GET", emailURL, nil)
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
		return nil, fmt.Errorf("GitHub 返回错误: %s, 状态码: %d", string(body), resp.StatusCode)
	}

	var emails []GithubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return emails, nil
}

func (s *GithubOAuthService) FindOrCreateUser(githubUser *GithubUserInfo) (*models.User, error) {
	db := database.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var user models.User
	err := db.Where("github_id = ?", githubUser.ID).First(&user).Error
	if err == nil {
		return &user, nil
	}

	username := githubUser.Login
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

	githubID := githubUser.ID
	newUser := models.User{
		Username:  username,
		Email:     githubUser.Email,
		Avatar:    githubUser.AvatarURL,
		Bio:       githubUser.Bio,
		Website:   githubUser.Blog,
		GithubID:  &githubID,
		PathAlias: utils.GenerateRandomString(16),
		Status:    1,
		Role:      3,
	}

	if err := db.Create(&newUser).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	return &newUser, nil
}
