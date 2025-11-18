package setting

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pixelpunk/internal/controllers/setting/dto"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/storage"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"strings"
	"time"
)

/* GetSettingsByGroupAsMap 按分组查询设置并以对象形式返回 */
func GetSettingsByGroupAsMap(group string) (*dto.SettingMapResponseDTO, error) {
	if cachedGroup, found := getSettingGroupFromCache(group); found {
		return cachedGroup, nil
	}

	db := database.GetDB()
	if db == nil {
		return &dto.SettingMapResponseDTO{Group: group, Settings: make(map[string]interface{}), UpdatedAt: ""}, nil
	}

	var settings []models.Setting
	query := db.Model(&models.Setting{}).Where("`group` = ?", group)

	if err := query.Order("`key` ASC").Find(&settings).Error; err != nil {
		return &dto.SettingMapResponseDTO{Group: group, Settings: make(map[string]interface{})}, nil
	}

	if len(settings) == 0 {
		result := &dto.SettingMapResponseDTO{Group: group, Settings: make(map[string]interface{})}
		setSettingGroupToCache(group, result)
		return result, nil
	}

	result := &dto.SettingMapResponseDTO{Group: group, Settings: make(map[string]interface{}, len(settings))}
	var latestUpdateTime time.Time
	for _, setting := range settings {
		value := parseSettingValue(setting)
		result.Settings[setting.Key] = value
		updateTime := time.Time(setting.UpdatedAt)
		if updateTime.After(latestUpdateTime) {
			latestUpdateTime = updateTime
		}
	}
	if !latestUpdateTime.IsZero() {
		result.UpdatedAt = latestUpdateTime.Format(time.RFC3339)
	}
	setSettingGroupToCache(group, result)
	return result, nil
}

/* GetGlobalSettingsGroups 获取多组全局设置 */
func GetGlobalSettingsGroups() (*dto.GlobalSettingsResponseDTO, error) {
	publicGroups := []string{"website", "website_info", "upload", "theme", "registration", "version", "ai", "vector", "guest", "appearance", "analytics"}
	result := &dto.GlobalSettingsResponseDTO{}

	for _, groupName := range publicGroups {
		groupSettings, err := GetSettingsByGroupAsMap(groupName)
		if err != nil {
			continue
		}

		switch groupName {
		case "website":
			result.Website = groupSettings.Settings
		case "website_info":
			result.WebsiteInfo = groupSettings.Settings
		case "upload":
			uploadConfig := make(map[string]interface{})
			allowedKeys := []string{"allowed_file_formats", "max_file_size", "max_batch_size", "content_detection_enabled", "sensitive_content_handling", "user_allowed_storage_durations", "user_default_storage_duration", "instant_upload_enabled"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					uploadConfig[key] = value
				}
			}

			uploadConfig["is_allow_chunk_upload"] = isLocalStorageDefault()

			result.Upload = uploadConfig
		case "theme":
			result.Theme = groupSettings.Settings
		case "registration":
			registrationConfig := make(map[string]interface{})
			allowedKeys := []string{"enable_registration", "email_verification"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					registrationConfig[key] = value
				}
			}
			result.Registration = registrationConfig
		case "version":
			versionConfig := make(map[string]interface{})
			allowedKeys := []string{"current_version", "build_time", "update_available", "last_update_check"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					versionConfig[key] = value
				}
			}
			result.Version = versionConfig
		case "ai":
			aiConfig := make(map[string]interface{})
			allowedKeys := []string{"ai_enabled"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					aiConfig[key] = value
				}
			}
			result.AI = aiConfig
		case "vector":
			vectorConfig := make(map[string]interface{})
			allowedKeys := []string{"vector_enabled"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					vectorConfig[key] = value
				}
			}
			result.Vector = vectorConfig
		case "guest":
			guestConfig := make(map[string]interface{})
			allowedKeys := []string{"enable_guest_upload", "guest_daily_limit", "guest_default_access_level", "guest_allowed_storage_durations", "guest_default_storage_duration"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					guestConfig[key] = value
				}
			}
			result.Guest = guestConfig
		case "appearance":
			result.Appearance = groupSettings.Settings
		case "analytics":
			analyticsConfig := make(map[string]interface{})
			allowedKeys := []string{"baidu_analytics_enabled", "baidu_analytics_site_id", "google_analytics_enabled", "google_analytics_measurement_id"}
			for _, key := range allowedKeys {
				if value, exists := groupSettings.Settings[key]; exists {
					analyticsConfig[key] = value
				}
			}
			result.Analytics = analyticsConfig
		}
	}

	oauthSettings, err := GetSettingsByGroupAsMap("oauth")
	if err == nil {
		if enabled, exists := oauthSettings.Settings["github_oauth_enabled"]; exists {
			if enabledBool, ok := enabled.(bool); ok {
				result.OAuthProviders.GithubEnabled = enabledBool
			}
		}
		if enabled, exists := oauthSettings.Settings["google_oauth_enabled"]; exists {
			if enabledBool, ok := enabled.(bool); ok {
				result.OAuthProviders.GoogleEnabled = enabledBool
			}
		}
		if enabled, exists := oauthSettings.Settings["linuxdo_oauth_enabled"]; exists {
			if enabledBool, ok := enabled.(bool); ok {
				result.OAuthProviders.LinuxdoEnabled = enabledBool
			}
		}
	}

	result.DeployMode = common.GetDeployMode()

	return result, nil
}

/* GetLegalDocuments 获取法律文档（隐私政策和服务条款） */
func GetLegalDocuments() (*dto.LegalDocumentsResponseDTO, error) {
	legalSettings, err := GetSettingsByGroupAsMap("legal")
	if err != nil {
		return &dto.LegalDocumentsResponseDTO{
			PrivacyPolicy:  "",
			TermsOfService: "",
		}, nil
	}

	result := &dto.LegalDocumentsResponseDTO{}

	if privacy, exists := legalSettings.Settings["privacy_policy_content"]; exists {
		if privacyStr, ok := privacy.(string); ok {
			result.PrivacyPolicy = privacyStr
		}
	}

	if terms, exists := legalSettings.Settings["terms_of_service_content"]; exists {
		if termsStr, ok := terms.(string); ok {
			result.TermsOfService = termsStr
		}
	}

	return result, nil
}

func isLocalStorageDefault() bool {
	// 快速检查：数据库不可用时返回true（启动时默认允许分片上传）
	db := database.GetDB()
	if db == nil {
		return true
	}

	// 检查数据库连接状态，避免阻塞
	if sqlDB, err := db.DB(); err != nil || sqlDB == nil {
		return true
	}

	defaultChannel, err := storage.GetDefaultChannel()
	if err != nil {
		return true
	}

	result := defaultChannel.Type == "local"
	return result
}

/* GetOAuthConfig 获取 OAuth 配置 */
func GetOAuthConfig() (*dto.OAuthConfigResponseDTO, error) {
	oauthSettings, err := GetSettingsByGroupAsMap("oauth")
	if err != nil {
		return &dto.OAuthConfigResponseDTO{
			Github: dto.GithubOAuthConfig{
				Enabled:      false,
				ClientID:     "",
				ClientSecret: "",
				RedirectURI:  "",
				Scope:        "user:email",
			},
		}, nil
	}

	result := &dto.OAuthConfigResponseDTO{}

	if enabled, exists := oauthSettings.Settings["github_oauth_enabled"]; exists {
		if enabledBool, ok := enabled.(bool); ok {
			result.Github.Enabled = enabledBool
		}
	}

	if clientID, exists := oauthSettings.Settings["github_oauth_client_id"]; exists {
		if clientIDStr, ok := clientID.(string); ok {
			result.Github.ClientID = clientIDStr
		}
	}

	if clientSecret, exists := oauthSettings.Settings["github_oauth_client_secret"]; exists {
		if clientSecretStr, ok := clientSecret.(string); ok {
			result.Github.ClientSecret = clientSecretStr
		}
	}

	if redirectURI, exists := oauthSettings.Settings["github_oauth_redirect_uri"]; exists {
		if redirectURIStr, ok := redirectURI.(string); ok {
			result.Github.RedirectURI = redirectURIStr
		}
	}

	if scope, exists := oauthSettings.Settings["github_oauth_scope"]; exists {
		if scopeStr, ok := scope.(string); ok {
			result.Github.Scope = scopeStr
		}
	}

	if result.Github.Scope == "" {
		result.Github.Scope = "user:email"
	}

	var sharedProxyConfig struct {
		ProxyType     string
		ProxyHost     string
		ProxyPort     string
		ProxyUsername string
		ProxyPassword string
		ProxyDynamic  bool
		ProxyAPIURL   string
	}

	if proxyType, exists := oauthSettings.Settings["oauth_proxy_type"]; exists {
		if proxyTypeStr, ok := proxyType.(string); ok {
			sharedProxyConfig.ProxyType = proxyTypeStr
		}
	}

	if proxyHost, exists := oauthSettings.Settings["oauth_proxy_host"]; exists {
		if proxyHostStr, ok := proxyHost.(string); ok {
			sharedProxyConfig.ProxyHost = proxyHostStr
		}
	}

	if proxyPort, exists := oauthSettings.Settings["oauth_proxy_port"]; exists {
		if proxyPortStr, ok := proxyPort.(string); ok {
			sharedProxyConfig.ProxyPort = proxyPortStr
		}
	}

	if proxyUsername, exists := oauthSettings.Settings["oauth_proxy_username"]; exists {
		if proxyUsernameStr, ok := proxyUsername.(string); ok {
			sharedProxyConfig.ProxyUsername = proxyUsernameStr
		}
	}

	if proxyPassword, exists := oauthSettings.Settings["oauth_proxy_password"]; exists {
		if proxyPasswordStr, ok := proxyPassword.(string); ok {
			sharedProxyConfig.ProxyPassword = proxyPasswordStr
		}
	}

	if proxyDynamic, exists := oauthSettings.Settings["oauth_proxy_dynamic"]; exists {
		if proxyDynamicBool, ok := proxyDynamic.(bool); ok {
			sharedProxyConfig.ProxyDynamic = proxyDynamicBool
		}
	}

	if proxyAPIURL, exists := oauthSettings.Settings["oauth_proxy_api_url"]; exists {
		if proxyAPIURLStr, ok := proxyAPIURL.(string); ok {
			sharedProxyConfig.ProxyAPIURL = proxyAPIURLStr
		}
	}

	// GitHub 代理配置：只读取是否启用代理的标志
	if proxyEnabled, exists := oauthSettings.Settings["github_oauth_proxy_enabled"]; exists {
		if proxyEnabledBool, ok := proxyEnabled.(bool); ok {
			result.Github.ProxyEnabled = proxyEnabledBool
			// 如果启用代理，使用统一的代理配置
			if proxyEnabledBool {
				result.Github.ProxyType = sharedProxyConfig.ProxyType
				result.Github.ProxyHost = sharedProxyConfig.ProxyHost
				result.Github.ProxyPort = sharedProxyConfig.ProxyPort
				result.Github.ProxyUsername = sharedProxyConfig.ProxyUsername
				result.Github.ProxyPassword = sharedProxyConfig.ProxyPassword
				result.Github.ProxyDynamic = sharedProxyConfig.ProxyDynamic
				result.Github.ProxyAPIURL = sharedProxyConfig.ProxyAPIURL
			}
		}
	}

	if enabled, exists := oauthSettings.Settings["google_oauth_enabled"]; exists {
		if enabledBool, ok := enabled.(bool); ok {
			result.Google.Enabled = enabledBool
		}
	}

	if clientID, exists := oauthSettings.Settings["google_oauth_client_id"]; exists {
		if clientIDStr, ok := clientID.(string); ok {
			result.Google.ClientID = clientIDStr
		}
	}

	if clientSecret, exists := oauthSettings.Settings["google_oauth_client_secret"]; exists {
		if clientSecretStr, ok := clientSecret.(string); ok {
			result.Google.ClientSecret = clientSecretStr
		}
	}

	if redirectURI, exists := oauthSettings.Settings["google_oauth_redirect_uri"]; exists {
		if redirectURIStr, ok := redirectURI.(string); ok {
			result.Google.RedirectURI = redirectURIStr
		}
	}

	if scope, exists := oauthSettings.Settings["google_oauth_scope"]; exists {
		if scopeStr, ok := scope.(string); ok {
			result.Google.Scope = scopeStr
		}
	}

	if result.Google.Scope == "" {
		result.Google.Scope = "openid email profile"
	}

	// Google 代理配置：只读取是否启用代理的标志
	if proxyEnabled, exists := oauthSettings.Settings["google_oauth_proxy_enabled"]; exists {
		if proxyEnabledBool, ok := proxyEnabled.(bool); ok {
			result.Google.ProxyEnabled = proxyEnabledBool
			// 如果启用代理，使用统一的代理配置
			if proxyEnabledBool {
				result.Google.ProxyType = sharedProxyConfig.ProxyType
				result.Google.ProxyHost = sharedProxyConfig.ProxyHost
				result.Google.ProxyPort = sharedProxyConfig.ProxyPort
				result.Google.ProxyUsername = sharedProxyConfig.ProxyUsername
				result.Google.ProxyPassword = sharedProxyConfig.ProxyPassword
				result.Google.ProxyDynamic = sharedProxyConfig.ProxyDynamic
				result.Google.ProxyAPIURL = sharedProxyConfig.ProxyAPIURL
			}
		}
	}

	// Linux DO 配置
	if enabled, exists := oauthSettings.Settings["linuxdo_oauth_enabled"]; exists {
		if enabledBool, ok := enabled.(bool); ok {
			result.Linuxdo.Enabled = enabledBool
		}
	}

	if clientID, exists := oauthSettings.Settings["linuxdo_oauth_client_id"]; exists {
		if clientIDStr, ok := clientID.(string); ok {
			result.Linuxdo.ClientID = clientIDStr
		}
	}

	if clientSecret, exists := oauthSettings.Settings["linuxdo_oauth_client_secret"]; exists {
		if clientSecretStr, ok := clientSecret.(string); ok {
			result.Linuxdo.ClientSecret = clientSecretStr
		}
	}

	if redirectURI, exists := oauthSettings.Settings["linuxdo_oauth_redirect_uri"]; exists {
		if redirectURIStr, ok := redirectURI.(string); ok {
			result.Linuxdo.RedirectURI = redirectURIStr
		}
	}

	if scope, exists := oauthSettings.Settings["linuxdo_oauth_scope"]; exists {
		if scopeStr, ok := scope.(string); ok {
			result.Linuxdo.Scope = scopeStr
		}
	}

	if result.Linuxdo.Scope == "" {
		result.Linuxdo.Scope = "user"
	}

	// LinuxDO 代理配置：只读取是否启用代理的标志
	if proxyEnabled, exists := oauthSettings.Settings["linuxdo_oauth_proxy_enabled"]; exists {
		if proxyEnabledBool, ok := proxyEnabled.(bool); ok {
			result.Linuxdo.ProxyEnabled = proxyEnabledBool
			// 如果启用代理，使用统一的代理配置
			if proxyEnabledBool {
				result.Linuxdo.ProxyType = sharedProxyConfig.ProxyType
				result.Linuxdo.ProxyHost = sharedProxyConfig.ProxyHost
				result.Linuxdo.ProxyPort = sharedProxyConfig.ProxyPort
				result.Linuxdo.ProxyUsername = sharedProxyConfig.ProxyUsername
				result.Linuxdo.ProxyPassword = sharedProxyConfig.ProxyPassword
				result.Linuxdo.ProxyDynamic = sharedProxyConfig.ProxyDynamic
				result.Linuxdo.ProxyAPIURL = sharedProxyConfig.ProxyAPIURL
			}
		}
	}

	return result, nil
}

/* DynamicProxyInfo 动态代理信息 */
type DynamicProxyInfo struct {
	Type     string // 代理类型: socks5
	Host     string // 代理地址
	Port     int    // 代理端口
	Username string // 用户名
	Password string // 密码
}

func FetchDynamicProxy(apiURL, username, password string) (*DynamicProxyInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("调用动态代理API失败: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取API响应失败: %v", err)
	}

	proxyInfo, err := parseProxyResponse(body, username, password)
	if err != nil {
		return nil, err
	}

	return proxyInfo, nil
}

func parseProxyResponse(body []byte, username, password string) (*DynamicProxyInfo, error) {
	if proxy, err := parseZhanDaYeFormat(body, username, password); err == nil {
		return proxy, nil
	}

	if proxy, err := parseGenericNestedFormat(body, username, password); err == nil {
		return proxy, nil
	}

	if proxy, err := parseSimpleArrayFormat(body, username, password); err == nil {
		return proxy, nil
	}

	if proxy, err := parseTextFormat(body, username, password); err == nil {
		return proxy, nil
	}

	return nil, fmt.Errorf("无法解析API响应，不支持的格式")
}

func parseZhanDaYeFormat(body []byte, username, password string) (*DynamicProxyInfo, error) {
	var result struct {
		Code string `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Count     int `json:"count"`
			ProxyList []struct {
				IP   string `json:"ip"`
				Port int    `json:"port"`
			} `json:"proxy_list"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Code != "10001" {
		return nil, fmt.Errorf("API返回错误: %s", result.Msg)
	}

	if result.Data.Count == 0 || len(result.Data.ProxyList) == 0 {
		return nil, fmt.Errorf("没有可用的代理")
	}

	proxy := result.Data.ProxyList[0]
	return &DynamicProxyInfo{
		Type:     "socks5",
		Host:     proxy.IP,
		Port:     proxy.Port,
		Username: username,
		Password: password,
	}, nil
}

func parseGenericNestedFormat(body []byte, username, password string) (*DynamicProxyInfo, error) {
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	var proxyList []interface{}
	for _, dataKey := range []string{"data", "proxies", "list", "result"} {
		if data, ok := result[dataKey]; ok {
			if dataMap, ok := data.(map[string]interface{}); ok {
				for _, listKey := range []string{"list", "proxy_list", "proxies", "items"} {
					if list, ok := dataMap[listKey]; ok {
						if arr, ok := list.([]interface{}); ok {
							proxyList = arr
							break
						}
					}
				}
			}
			if arr, ok := data.([]interface{}); ok {
				proxyList = arr
			}
			if len(proxyList) > 0 {
				break
			}
		}
	}

	if len(proxyList) == 0 {
		return nil, fmt.Errorf("未找到代理列表")
	}

	if proxyMap, ok := proxyList[0].(map[string]interface{}); ok {
		var host string
		var port int

		for _, hostKey := range []string{"ip", "host", "address", "proxy_ip", "server"} {
			if h, ok := proxyMap[hostKey].(string); ok && h != "" {
				host = h
				break
			}
		}

		for _, portKey := range []string{"port", "proxy_port"} {
			if p, ok := proxyMap[portKey].(float64); ok {
				port = int(p)
				break
			}
		}

		if host != "" && port > 0 {
			return &DynamicProxyInfo{
				Type:     "socks5",
				Host:     host,
				Port:     port,
				Username: username,
				Password: password,
			}, nil
		}
	}

	return nil, fmt.Errorf("解析代理信息失败")
}

func parseSimpleArrayFormat(body []byte, username, password string) (*DynamicProxyInfo, error) {
	var proxyList []map[string]interface{}
	if err := json.Unmarshal(body, &proxyList); err != nil {
		return nil, err
	}

	if len(proxyList) == 0 {
		return nil, fmt.Errorf("代理列表为空")
	}

	proxyMap := proxyList[0]
	var host string
	var port int

	for _, hostKey := range []string{"ip", "host", "address"} {
		if h, ok := proxyMap[hostKey].(string); ok {
			host = h
			break
		}
	}

	for _, portKey := range []string{"port"} {
		if p, ok := proxyMap[portKey].(float64); ok {
			port = int(p)
			break
		}
	}

	if host != "" && port > 0 {
		return &DynamicProxyInfo{
			Type:     "socks5",
			Host:     host,
			Port:     port,
			Username: username,
			Password: password,
		}, nil
	}

	return nil, fmt.Errorf("解析失败")
}

func parseTextFormat(body []byte, username, password string) (*DynamicProxyInfo, error) {
	text := string(body)
	lines := strings.Split(strings.TrimSpace(text), "\n")

	if len(lines) == 0 {
		return nil, fmt.Errorf("文本为空")
	}

	firstLine := strings.TrimSpace(lines[0])
	parts := strings.Split(firstLine, ":")

	if len(parts) >= 2 {
		host := parts[0]
		port := 0
		fmt.Sscanf(parts[1], "%d", &port)

		// 如果文本格式包含用户名密码，使用文本中的
		if len(parts) >= 4 {
			username = parts[2]
			password = parts[3]
		}

		if host != "" && port > 0 {
			return &DynamicProxyInfo{
				Type:     "socks5",
				Host:     host,
				Port:     port,
				Username: username,
				Password: password,
			}, nil
		}
	}

	return nil, fmt.Errorf("文本格式解析失败")
}

func TestProxyConnection(req *dto.TestProxyDTO) (*dto.TestProxyResponseDTO, error) {
	var proxyHost string
	var proxyPort int
	var proxyType string

	if req.ProxyDynamic && req.ProxyAPIURL != "" {
		dynamicProxy, err := FetchDynamicProxy(req.ProxyAPIURL, req.ProxyUsername, req.ProxyPassword)
		if err != nil {
			return &dto.TestProxyResponseDTO{
				Success: false,
				Message: fmt.Sprintf("获取动态代理失败: %v", err),
				Latency: 0,
			}, nil
		}

		proxyHost = dynamicProxy.Host
		proxyPort = dynamicProxy.Port
		proxyType = dynamicProxy.Type
	} else {
		proxyHost = req.ProxyHost
		proxyType = req.ProxyType
		fmt.Sscanf(req.ProxyPort, "%d", &proxyPort)
	}

	var proxyURL string
	if req.ProxyUsername != "" && req.ProxyPassword != "" {
		proxyURL = fmt.Sprintf("%s://%s:%s@%s:%d",
			proxyType,
			url.QueryEscape(req.ProxyUsername),
			url.QueryEscape(req.ProxyPassword),
			proxyHost,
			proxyPort,
		)
	} else {
		proxyURL = fmt.Sprintf("%s://%s:%d",
			proxyType,
			proxyHost,
			proxyPort,
		)
	}

	proxy, err := url.Parse(proxyURL)
	if err != nil {
		return &dto.TestProxyResponseDTO{
			Success: false,
			Message: fmt.Sprintf("代理 URL 格式错误: %v", err),
			Latency: 0,
		}, nil
	}

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}
	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}

	startTime := time.Now()
	resp, err := client.Get("https://api.github.com/")
	latency := time.Since(startTime).Milliseconds()

	if err != nil {
		return &dto.TestProxyResponseDTO{
			Success: false,
			Message: fmt.Sprintf("代理连接失败: %v", err),
			Latency: latency,
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusForbidden {
		return &dto.TestProxyResponseDTO{
			Success: false,
			Message: fmt.Sprintf("GitHub API 返回错误状态码: %d", resp.StatusCode),
			Latency: latency,
		}, nil
	}

	return &dto.TestProxyResponseDTO{
		Success: true,
		Message: "代理连接成功，可以访问 GitHub API",
		Latency: latency,
	}, nil
}
