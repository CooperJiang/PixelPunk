package storage

import (
	"fmt"

	su "pixelpunk/pkg/storage/utils"

	"github.com/google/uuid"
)

// IsValidFolderName checks if a name is safe (no forbidden characters).
func IsValidFolderName(name string) bool { return su.SanitizeFileName(name) == name }

// SanitizeFileName exposes utils sanitizer at storage package.
func SanitizeFileName(name string) string { return su.SanitizeFileName(name) }

// GenerateFolderID returns a short unique folder id (16 chars).
func GenerateFolderID() string { return uuid.NewString()[:8] + uuid.NewString()[9:17] }

// EnsureUserDirectories 确保用户目录结构存在（扁平存储）
func EnsureUserDirectories(userID uint) error {
	return su.EnsureUserDirectories(fmt.Sprintf("%d", userID))
}

// CreateUserFilePaths 生成用户文件的扁平存储路径
func CreateUserFilePaths(userID uint, fileName string) (string, string) {
	return su.CreateUserFilePaths(fmt.Sprintf("%d", userID), fileName)
}
