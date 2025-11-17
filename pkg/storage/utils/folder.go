package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pixelpunk/pkg/logger"
)

// EnsureUserDirectories 确保用户基础目录存在（扁平结构）
func EnsureUserDirectories(userID string) error {
	userDir := fmt.Sprintf("user_%s", userID)

	userPath := filepath.Join("uploads", userDir)
	if err := os.MkdirAll(userPath, 0755); err != nil {
		logger.Error("创建用户目录失败: %v", err)
		return fmt.Errorf("创建用户目录失败: %w", err)
	}

	imageDir := filepath.Join(userPath, "files")
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		logger.Error("创建用户文件目录失败: %v", err)
		return fmt.Errorf("创建用户文件目录失败: %w", err)
	}

	thumbDir := filepath.Join(userPath, "thumbnails")
	if err := os.MkdirAll(thumbDir, 0755); err != nil {
		logger.Error("创建用户缩略图目录失败: %v", err)
		return fmt.Errorf("创建用户缩略图目录失败: %w", err)
	}

	return nil
}

// CreateUserFilePaths 生成用户文件的扁平存储路径
func CreateUserFilePaths(userID string, fileName string) (string, string) {
	userDir := fmt.Sprintf("user_%s", userID)

	// 文件路径：uploads/user_N/files/filename.jpg
	imagePath := filepath.Join("uploads", userDir, "files", fileName)

	// 缩略图路径：uploads/user_N/thumbnails/filename.jpg
	thumbPath := filepath.Join("uploads", userDir, "thumbnails", fileName)

	return imagePath, thumbPath
}

// 注意：在新的扁平存储结构下，不再需要物理文件夹的移动和删除操作
// 所有文件夹操作都是纯逻辑的，仅在数据库层面进行
// 物理文件始终存储在用户的 files/ 和 thumbnails/ 目录下

// DirExists 检查目录是否存在
func DirExists(dirPath string) bool {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// FileExists 检查文件是否存在
func FileExists(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// SanitizeFileName 清理文件名，确保安全
func SanitizeFileName(filename string) string {
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", "\n", "\r", "\t"}
	sanitized := filename

	for _, char := range unsafe {
		sanitized = strings.ReplaceAll(sanitized, char, "_")
	}

	sanitized = strings.TrimSpace(sanitized)

	if sanitized == "" {
		sanitized = "unnamed"
	}

	return sanitized
}

// CopyFile 复制文件
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("读取源文件失败: %w", err)
	}

	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("创建目标目录失败: %w", err)
	}

	if err := os.WriteFile(dst, data, 0644); err != nil {
		return fmt.Errorf("写入目标文件失败: %w", err)
	}

	return nil
}
