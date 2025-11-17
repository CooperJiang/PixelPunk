package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strings"

	"github.com/google/uuid"
)

// 字符集
const (
	letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	if length <= 0 {
		return ""
	}

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return generateFallbackRandomString(length)
		}
		result[i] = letters[num.Int64()]
	}

	return string(result)
}

// generateFallbackRandomString 生成备用随机字符串（当加密随机数生成失败时使用）
func generateFallbackRandomString(length int) string {
	b := make([]byte, length*2)
	if _, err := rand.Read(b); err != nil {
		// 如果随机读取也失败，使用一个简单的字符串
		return "fallback_random_string"
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}

// GenerateFileID 生成32位文件ID
func GenerateFileID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}
