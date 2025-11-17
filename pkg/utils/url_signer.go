package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
)

// URLSigner URL签名器
type URLSigner struct {
	secret string
	cache  map[string]CachedSignature
	mutex  sync.RWMutex
}

// CachedSignature 缓存的签名
type CachedSignature struct {
	URL     string
	Expires time.Time
}

// URLSigner相关常量
const (
	SIGNATURE_DURATION     = 10 * time.Minute // 10分钟有效期
	TIME_WINDOW_ALIGN      = 5 * time.Minute  // 5分钟对齐窗口
	CACHE_CLEANUP_INTERVAL = 5 * time.Minute  // 缓存清理间隔
)

// 全局签名器实例
var (
	globalSigner *URLSigner
	signerOnce   sync.Once
)

func NewURLSigner(secret string) *URLSigner {
	signer := &URLSigner{
		secret: secret,
		cache:  make(map[string]CachedSignature),
	}

	go signer.startCacheCleanup()

	return signer
}

// SignFileURL 为文件URL生成签名
func (s *URLSigner) SignFileURL(fileID string, duration time.Duration) string {
	// 使用时间窗口对齐减少重复签名
	now := time.Now()
	alignedTime := now.Truncate(TIME_WINDOW_ALIGN).Add(duration)
	expiry := alignedTime.Unix()

	cacheKey := fmt.Sprintf("file:%s:%d", fileID, expiry)
	if cachedURL := s.getCachedURL(cacheKey); cachedURL != "" {
		return cachedURL
	}

	message := fmt.Sprintf("%s:%d", fileID, expiry)
	signature := s.generateSignature(message)

	// 构建带签名的相对路径 - 使用新架构的FileID格式
	relativeURL := fmt.Sprintf("/f/%s?t=%d&s=%s",
		fileID, expiry, signature)

	// 使用现有的GetSystemFileURL函数添加基础域名
	signedURL := GetSystemFileURL(relativeURL)

	s.setCachedURL(cacheKey, signedURL, alignedTime)

	return signedURL
}

// SignThumbURL 为缩略图URL生成签名
func (s *URLSigner) SignThumbURL(fileID string, duration time.Duration) string {
	now := time.Now()
	alignedTime := now.Truncate(TIME_WINDOW_ALIGN).Add(duration)
	expiry := alignedTime.Unix()

	cacheKey := fmt.Sprintf("thumb:%s:%d", fileID, expiry)
	if cachedURL := s.getCachedURL(cacheKey); cachedURL != "" {
		return cachedURL
	}

	// 缩略图使用不同的签名前缀
	message := fmt.Sprintf("thumb:%s:%d", fileID, expiry)
	signature := s.generateSignature(message)

	// 构建带签名的相对路径 - 使用新架构的缩略图格式
	relativeURL := fmt.Sprintf("/t/%s?t=%d&s=%s",
		fileID, expiry, signature)

	// 使用现有的GetSystemFileURL函数添加基础域名
	signedURL := GetSystemFileURL(relativeURL)

	s.setCachedURL(cacheKey, signedURL, alignedTime)

	return signedURL
}

// VerifyFileURL 验证文件URL签名
func (s *URLSigner) VerifyFileURL(path, timeParam, signatureParam string) bool {
	if timeParam == "" || signatureParam == "" {
		return false
	}

	expiry, err := strconv.ParseInt(timeParam, 10, 64)
	if err != nil {
		return false
	}

	if time.Now().Unix() > expiry {
		return false // 已过期
	}

	fileID := s.extractFileIDFromPath(path)
	if fileID == "" {
		return false
	}

	// 生成期望的签名（新架构：/f/{id} 与 /t/{id}）
	var message string
	if strings.HasPrefix(path, "/t/") {
		message = fmt.Sprintf("thumb:%s:%d", fileID, expiry)
	} else {
		message = fmt.Sprintf("%s:%d", fileID, expiry)
	}

	expectedSignature := s.generateSignature(message)

	// 对比签名（使用恒定时间比较防止时序攻击）
	return hmac.Equal([]byte(expectedSignature), []byte(signatureParam))
}

// generateSignature 生成HMAC签名
func (s *URLSigner) generateSignature(message string) string {
	mac := hmac.New(sha256.New, []byte(s.secret))
	mac.Write([]byte(message))

	// 生成base64编码，取前16位作为签名
	fullSignature := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	if len(fullSignature) > 16 {
		return fullSignature[:16]
	}
	return fullSignature
}

// extractFileIDFromPath 从路径中提取文件ID（/f/:id 或 /t/:id）
func (s *URLSigner) extractFileIDFromPath(path string) string {
	// 处理新架构路径: /f/:id 或 /t/:id
	parts := strings.Split(path, "/")

	if len(parts) >= 3 {
		if parts[1] == "f" || parts[1] == "t" {
			return parts[2]
		}
	}

	return ""
}

// getCachedURL 获取缓存的URL
func (s *URLSigner) getCachedURL(cacheKey string) string {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if cached, exists := s.cache[cacheKey]; exists && time.Now().Before(cached.Expires) {
		return cached.URL
	}

	return ""
}

// setCachedURL 设置缓存的URL
func (s *URLSigner) setCachedURL(cacheKey, url string, expires time.Time) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[cacheKey] = CachedSignature{
		URL:     url,
		Expires: expires,
	}
}

// startCacheCleanup 启动缓存清理
func (s *URLSigner) startCacheCleanup() {
	ticker := time.NewTicker(CACHE_CLEANUP_INTERVAL)
	defer ticker.Stop()

	for range ticker.C {
		s.cleanupExpiredCache()
	}
}

// cleanupExpiredCache 清理过期缓存
func (s *URLSigner) cleanupExpiredCache() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for key, cached := range s.cache {
		if now.After(cached.Expires) {
			delete(s.cache, key)
		}
	}
}

func GetURLSigner() *URLSigner {
	signerOnce.Do(func() {
		secret := GetSigningSecret()
		globalSigner = NewURLSigner(secret)
	})
	return globalSigner
}

func GetSigningSecret() string {
	// 优先从环境变量获取
	if secret := os.Getenv("URL_SIGNING_SECRET"); secret != "" {
		return secret
	}

	if secret := getSecretFromDatabase(); secret != "" {
		return secret
	}

	// 生成默认密钥（生产环境应该更换）
	return "default-url-signing-secret-change-in-production"
}

// getSecretFromDatabase 从数据库获取密钥
func getSecretFromDatabase() string {
	// 优先查找专用的URL签名密钥
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", "url_signing_secret").First(&setting).Error; err == nil {
		return setting.GetStringValue()
	}

	// 如果没有，使用JWT密钥作为基础
	if err := database.DB.Where("`key` = ?", "jwt_secret").First(&setting).Error; err == nil {
		jwtSecret := setting.GetStringValue()
		if jwtSecret != "" {
			return jwtSecret + "-url-signing"
		}
	}

	return ""
}

// ResetURLSigner 重置全局签名器（用于测试或密钥更新）
func ResetURLSigner() {
	globalSigner = nil
	signerOnce = sync.Once{}
}
