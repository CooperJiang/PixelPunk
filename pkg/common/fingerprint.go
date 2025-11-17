package common

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// GenerateFingerprint 生成浏览器指纹
func GenerateFingerprint(r *http.Request) string {
	var components []string

	ip := getClientIP(r)
	if ip != "" {
		privacyIP := maskIP(ip)
		components = append(components, privacyIP)
	}

	ua := r.UserAgent()
	if ua != "" {
		// 简单提取浏览器和操作系统信息
		browser := extractBrowser(ua)
		os := extractOS(ua)
		components = append(components, browser, os)
	}

	if lang := r.Header.Get("Accept-Language"); lang != "" {
		// 只取前两个字符，如 zh-CN -> zh
		if len(lang) > 2 {
			lang = lang[:2]
		}
		components = append(components, lang)
	}

	hour := fmt.Sprintf("%02d", time.Now().Hour())
	components = append(components, hour)

	components = append(components, "pixelpunk-host")

	fingerprint := generateHash(strings.Join(components, "|"))

	return fingerprint
}

// getClientIP 获取客户端真实IP
func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		// X-Forwarded-For 可能包含多个IP，取第一个
		ips := strings.Split(ip, ",")
		return strings.TrimSpace(ips[0])
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	if ip := r.Header.Get("CF-Connecting-IP"); ip != "" {
		return ip
	}

	// 直接使用 RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}

	return ip
}

// maskIP 隐藏IP地址的最后一段，保护隐私
func maskIP(ip string) string {
	if strings.Contains(ip, ":") {
		parts := strings.Split(ip, ":")
		if len(parts) >= 3 {
			parts[len(parts)-1] = "xxxx"
			return strings.Join(parts, ":")
		}
	} else {
		parts := strings.Split(ip, ".")
		if len(parts) == 4 {
			parts[3] = "xxx"
			return strings.Join(parts, ".")
		}
	}
	return ip
}

// generateHash 生成哈希值作为指纹
func generateHash(data string) string {
	// 使用SHA256生成指纹
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hash := hasher.Sum(nil)

	// 使用MD5进行二次哈希，生成更短的指纹
	md5Hasher := md5.New()
	md5Hasher.Write(hash)
	finalHash := md5Hasher.Sum(nil)

	return hex.EncodeToString(finalHash)[:16]
}

// ValidateFingerprint 验证指纹格式
func ValidateFingerprint(fingerprint string) bool {
	if len(fingerprint) != 16 {
		return false
	}

	_, err := hex.DecodeString(fingerprint)
	return err == nil
}

// 用于处理同一用户可能产生不同指纹的情况
func GetFingerprintKey(fingerprint string) string {
	// 可以实现更复杂的逻辑，比如：
	// - 生成多个可能的指纹变体
	// - 使用模糊匹配
	// - 添加时间窗口

	return fingerprint
}

// 用于检测可能的相同用户
func GetSimilarFingerprints(fingerprint string) []string {
	var similar []string

	hour := time.Now().Hour()
	for h := hour - 1; h <= hour+1; h++ {
		if h >= 0 && h < 24 {
			// 这里需要重新生成指纹，简化实现
			// 实际使用时应该保存原始数据以便重新生成
		}
	}

	// 如果指纹包含IP信息，可以生成同一IP段的指纹

	return similar
}

// extractBrowser 从User-Agent中提取浏览器信息
func extractBrowser(ua string) string {
	browsers := []struct {
		pattern string
		name    string
	}{
		{`Firefox/`, "Firefox"},
		{`Chrome/`, "Chrome"},
		{`Safari/`, "Safari"},
		{`Edge/`, "Edge"},
		{`Opera/`, "Opera"},
		{`PostmanRuntime/`, "Postman"},
		{`curl/`, "Curl"},
	}

	for _, b := range browsers {
		if matched, _ := regexp.MatchString(b.pattern, ua); matched {
			return b.name
		}
	}

	return "Unknown"
}

// extractOS 从User-Agent中提取操作系统信息
func extractOS(ua string) string {
	// 简单的操作系统检测
	osList := []struct {
		pattern string
		name    string
	}{
		{`Windows NT`, "Windows"},
		{`Mac OS X`, "macOS"},
		{`Linux`, "Linux"},
		{`Android`, "Android"},
		{`iPhone OS`, "iOS"},
		{`iPad OS`, "iPadOS"},
	}

	for _, os := range osList {
		if matched, _ := regexp.MatchString(os.pattern, ua); matched {
			return os.name
		}
	}

	return "Unknown"
}

// IsSuspiciousFingerprint 检查指纹是否可疑
// 用于检测可能的攻击或滥用
func IsSuspiciousFingerprint(fingerprint string, requestCount int) bool {
	if requestCount > 100 { // 5分钟内超过100次请求
		return true
	}

	if !ValidateFingerprint(fingerprint) {
		return true
	}

	return false
}

// FingerprintInfo 指纹信息
type FingerprintInfo struct {
	Fingerprint string `json:"fingerprint"`
	IP          string `json:"ip"`
	UserAgent   string `json:"user_agent"`
	Language    string `json:"language"`
	Hour        string `json:"hour"`
	GeneratedAt int64  `json:"generated_at"`
}

// ParseFingerprint 解析指纹信息（用于调试）
// 注意：这是一个单向过程，无法从指纹还原原始信息
func ParseFingerprint(r *http.Request) *FingerprintInfo {
	return &FingerprintInfo{
		Fingerprint: GenerateFingerprint(r),
		IP:          getClientIP(r),
		UserAgent:   r.UserAgent(),
		Language:    r.Header.Get("Accept-Language"),
		Hour:        fmt.Sprintf("%02d", time.Now().Hour()),
		GeneratedAt: time.Now().Unix(),
	}
}
