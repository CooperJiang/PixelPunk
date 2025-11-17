package url

import (
	"fmt"
	"strings"
)

// AccessMode 访问模式
type AccessMode string

const (
	AccessModeProxy  AccessMode = "proxy"  // 代理模式：通过服务器代理访问
	AccessModeDirect AccessMode = "direct" // 直链模式：直接访问云存储URL
)

// URLStrategy URL访问策略，统一管理URL生成逻辑
type URLStrategy struct {
	hideRemoteURL bool   // 是否隐藏远程URL（全局或渠道配置）
	accessControl string // 访问控制类型 (public/private)
	forceHTTPS    bool   // 是否强制HTTPS
	customDomain  string // 自定义域名
	cdnDomain     string // CDN域名
	baseURL       string // 基础URL

	accessMode   AccessMode // 实际使用的访问模式
	requiresAuth bool       // 是否需要认证
}

// StrategyConfig URL策略配置
type StrategyConfig struct {
	GlobalHideRemoteURL bool
	BaseURL             string
	ForceHTTPS          bool

	ChannelHideRemoteURL     bool
	ChannelHasHideURLSetting bool
	AccessControl            string
	CustomDomain             string
	CDNDomain                string
}

func NewURLStrategy(config StrategyConfig) *URLStrategy {
	strategy := &URLStrategy{
		baseURL:       config.BaseURL,
		forceHTTPS:    config.ForceHTTPS,
		customDomain:  config.CustomDomain,
		cdnDomain:     config.CDNDomain,
		accessControl: config.AccessControl,
	}

	// 计算是否隐藏远程URL的最终结果
	strategy.hideRemoteURL = strategy.calculateHideRemoteURL(config)

	strategy.accessMode = strategy.calculateAccessMode()

	strategy.requiresAuth = strategy.calculateRequiresAuth()

	return strategy
}

// calculateHideRemoteURL 计算最终的隐藏远程URL策略
// 优先级：私有访问 > 全局配置 > 渠道配置
func (s *URLStrategy) calculateHideRemoteURL(config StrategyConfig) bool {
	// 私有访问控制强制隐藏远程URL
	if config.AccessControl == "private" {
		return true
	}

	if config.GlobalHideRemoteURL {
		return true
	}

	// 如果渠道有独立配置，使用渠道配置
	if config.ChannelHasHideURLSetting {
		return config.ChannelHideRemoteURL
	}

	return false
}

// calculateAccessMode 计算访问模式
func (s *URLStrategy) calculateAccessMode() AccessMode {
	if s.hideRemoteURL {
		return AccessModeProxy
	}
	return AccessModeDirect
}

// calculateRequiresAuth 计算是否需要认证
func (s *URLStrategy) calculateRequiresAuth() bool {
	return s.accessControl == "private"
}

func (s *URLStrategy) GetAccessMode() AccessMode {
	return s.accessMode
}

// ShouldUseProxy 是否应该使用代理模式
func (s *URLStrategy) ShouldUseProxy() bool {
	return s.accessMode == AccessModeProxy
}

// RequiresAuth 是否需要认证
func (s *URLStrategy) RequiresAuth() bool {
	return s.requiresAuth
}


// BuildDirectURL 构建直链URL
func (s *URLStrategy) BuildDirectURL(remoteURL string) string {
	if remoteURL == "" {
		return ""
	}

	finalURL := s.applyProtocolSettings(remoteURL)

	// 应用自定义域名（如果配置了）
	if s.customDomain != "" {
		finalURL = s.applyCustomDomain(finalURL)
	}

	return finalURL
}

// BuildCDNURL 构建CDN URL
func (s *URLStrategy) BuildCDNURL(remoteURL string) string {
	if s.cdnDomain == "" {
		return s.BuildDirectURL(remoteURL)
	}

	// 提取路径部分，替换域名为CDN域名
	if remoteURL != "" {
		path := s.extractPathFromURL(remoteURL)
		protocol := "https"
		if !s.forceHTTPS && strings.HasPrefix(remoteURL, "http://") {
			protocol = "http"
		}
		return fmt.Sprintf("%s://%s%s", protocol, s.cdnDomain, path)
	}

	return ""
}

// applyProtocolSettings 应用协议设置
func (s *URLStrategy) applyProtocolSettings(url string) string {
	if s.forceHTTPS && strings.HasPrefix(url, "http://") {
		return strings.Replace(url, "http://", "https://", 1)
	}
	if !s.forceHTTPS && strings.HasPrefix(url, "https://") {
		return strings.Replace(url, "https://", "http://", 1)
	}
	return url
}

// applyCustomDomain 应用自定义域名
func (s *URLStrategy) applyCustomDomain(url string) string {
	if s.customDomain == "" {
		return url
	}

	path := s.extractPathFromURL(url)
	protocol := "https"
	if !s.forceHTTPS && strings.HasPrefix(url, "http://") {
		protocol = "http"
	}

	return fmt.Sprintf("%s://%s%s", protocol, s.customDomain, path)
}

// extractPathFromURL 从URL中提取路径部分
func (s *URLStrategy) extractPathFromURL(url string) string {
	if url == "" {
		return ""
	}

	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// 查找第一个斜杠，之后的部分就是路径
	if idx := strings.Index(url, "/"); idx != -1 {
		return url[idx:]
	}

	return "/"
}

func (s *URLStrategy) GetDebugInfo() map[string]interface{} {
	return map[string]interface{}{
		"hideRemoteURL": s.hideRemoteURL,
		"accessControl": s.accessControl,
		"accessMode":    s.accessMode,
		"requiresAuth":  s.requiresAuth,
		"forceHTTPS":    s.forceHTTPS,
		"customDomain":  s.customDomain,
		"cdnDomain":     s.cdnDomain,
		"baseURL":       s.baseURL,
	}
}

// String 返回策略的字符串表示
func (s *URLStrategy) String() string {
	return fmt.Sprintf("URLStrategy{mode=%s, hideRemote=%v, auth=%v}",
		s.accessMode, s.hideRemoteURL, s.requiresAuth)
}
