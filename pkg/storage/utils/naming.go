package utils

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// MakeFileName 基于 UUID+短时间戳 与扩展名生成名称
func MakeFileName(ext string) string {
	u := strings.ReplaceAll(uuid.New().String(), "-", "")
	ts := fmt.Sprintf("%04d", time.Now().Unix()%10000)
	ext = strings.TrimPrefix(strings.ToLower(ext), ".")
	if ext == "" {
		ext = "jpg"
	}
	return u + ts + "." + ext
}

// MakeThumbName 由原文件名派生缩略图名（优先 _thumb 后缀）
func MakeThumbName(originalFileName, thumbExt string) string {
	base := strings.TrimSuffix(originalFileName, filepath.Ext(originalFileName))
	thumbExt = strings.TrimPrefix(strings.ToLower(thumbExt), ".")
	if thumbExt == "" {
		thumbExt = "jpg"
	}
	return base + "_thumb." + thumbExt
}


// BuildLogicalPath (folder)/file - 不含 user_N，且去除开头斜杠
func BuildLogicalPath(folderPath, fileName string) string {
	clean := strings.TrimPrefix(folderPath, "/")
	if clean == "" {
		return fileName
	}
	return filepath.Join(clean, fileName)
}
