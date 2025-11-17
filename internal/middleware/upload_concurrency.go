package middleware

import (
	"crypto/md5"
	"fmt"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/utils"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestTracker struct {
	Count    int
	LastUsed time.Time
}

var (
	userRequestMap  = make(map[uint]*RequestTracker)
	guestRequestMap = make(map[string]*RequestTracker)
	mapMutex        = &sync.RWMutex{}
	cleanupTicker   *time.Ticker
)

func init() {
	cleanupTicker = time.NewTicker(5 * time.Minute)
	go func() {
		for range cleanupTicker.C {
			cleanupExpiredTrackers()
		}
	}()
}

func cleanupExpiredTrackers() {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	now := time.Now()
	expireDuration := 10 * time.Minute // 10分钟未使用的记录将被清理

	for userID, tracker := range userRequestMap {
		if now.Sub(tracker.LastUsed) > expireDuration && tracker.Count == 0 {
			delete(userRequestMap, userID)
		}
	}

	for fingerprint, tracker := range guestRequestMap {
		if now.Sub(tracker.LastUsed) > expireDuration && tracker.Count == 0 {
			delete(guestRequestMap, fingerprint)
		}
	}
}

func UploadConcurrencyLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetCurrentUserID(c)

		if userID > 0 {
			if !checkUserConcurrency(userID, c) {
				return
			}
		} else {
			if !checkGuestConcurrency(c) {
				return
			}
		}

		c.Next()

		if userID > 0 {
			releaseUserConcurrency(userID)
		} else {
			releaseGuestConcurrency(c)
		}
	}
}

func checkUserConcurrency(userID uint, c *gin.Context) bool {
	maxConcurrent := getUserMaxConcurrent()
	if maxConcurrent <= 0 {
		maxConcurrent = 3 // 默认值
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	tracker, exists := userRequestMap[userID]
	if !exists {
		tracker = &RequestTracker{Count: 0, LastUsed: time.Now()}
		userRequestMap[userID] = tracker
	}

	if tracker.Count >= maxConcurrent {
		errors.HandleError(c, errors.New(errors.CodeRateLimited,
			fmt.Sprintf("并发上传数量超出限制，最多允许%d个并发请求", maxConcurrent)))
		return false
	}

	tracker.Count++
	tracker.LastUsed = time.Now()

	return true
}

func checkGuestConcurrency(c *gin.Context) bool {
	const maxConcurrent = 1

	fingerprint := getGuestFingerprint(c)
	if fingerprint == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidRequest, "无法识别游客身份"))
		return false
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	tracker, exists := guestRequestMap[fingerprint]
	if !exists {
		tracker = &RequestTracker{Count: 0, LastUsed: time.Now()}
		guestRequestMap[fingerprint] = tracker
	}

	if tracker.Count >= maxConcurrent {
		errors.HandleError(c, errors.New(errors.CodeRateLimited,
			"游客模式下只能同时上传1个文件，请等待当前上传完成"))
		return false
	}

	tracker.Count++
	tracker.LastUsed = time.Now()

	return true
}

func releaseUserConcurrency(userID uint) {
	mapMutex.Lock()
	defer mapMutex.Unlock()

	if tracker, exists := userRequestMap[userID]; exists {
		tracker.Count--
		if tracker.Count < 0 {
			tracker.Count = 0
		}
		tracker.LastUsed = time.Now()
	}
}

func releaseGuestConcurrency(c *gin.Context) {
	fingerprint := getGuestFingerprint(c)
	if fingerprint == "" {
		return
	}

	mapMutex.Lock()
	defer mapMutex.Unlock()

	if tracker, exists := guestRequestMap[fingerprint]; exists {
		tracker.Count--
		if tracker.Count < 0 {
			tracker.Count = 0
		}
		tracker.LastUsed = time.Now()
	}
}

func getUserMaxConcurrent() int {
	uploadSettings, err := setting.GetSettingsByGroupAsMap("upload")
	if err != nil {
		logger.Error("获取上传设置失败: %v", err)
		return 3 // 返回默认值
	}

	if value, exists := uploadSettings.Settings["client_max_concurrent_uploads"]; exists {
		switch v := value.(type) {
		case float64:
			return int(v)
		case int:
			return v
		default:
			logger.Warn("client_max_concurrent_uploads配置类型错误: %T", value)
		}
	}

	return 3 // 默认值
}

func getGuestFingerprint(c *gin.Context) string {
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "unknown"
	}

	clientIP := utils.GetClientIP(c)

	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage == "" {
		acceptLanguage = "unknown"
	}

	fingerprint := fmt.Sprintf("%s|%s|%s", userAgent, clientIP, acceptLanguage)

	hash := md5.Sum([]byte(fingerprint))
	return fmt.Sprintf("%x", hash)
}

func GetConcurrencyStatus() gin.H {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	userStatus := make(map[string]int)
	for userID, tracker := range userRequestMap {
		userStatus[fmt.Sprintf("user_%d", userID)] = tracker.Count
	}

	guestStatus := make(map[string]int)
	for fingerprint, tracker := range guestRequestMap {
		guestStatus[fingerprint] = tracker.Count
	}

	return gin.H{
		"users":        userStatus,
		"guests":       guestStatus,
		"total_users":  len(userRequestMap),
		"total_guests": len(guestRequestMap),
	}
}
