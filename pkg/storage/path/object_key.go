package path

import (
	"path/filepath"
	"strings"

	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage/tenant"
)

// EnsureObjectKey 将任意输入（对象键/逻辑路径/完整URL）规范为对象键
// - userID: 用于构建对象键（游客=0 会自动落在 files|thumbnails/guest 下）
// - input: 对象键、逻辑路径（可能含 thumb/）、或 http(s) URL
// - isThumb: 提示是否缩略图（可为 false，内部会根据文件名 _thumb. 再次判定）
func EnsureObjectKey(userID uint, input string, isThumb bool) string {
	if input == "" {
		return ""
	}

	if strings.HasPrefix(input, "files/") || strings.HasPrefix(input, "thumbnails/") {
		return CleanObjectPath(input)
	}

	// 如果是URL，先提取对象键再处理
	if IsHTTPURL(input) {
		if p := ExtractObjectPathFromURL(input); p != "" {
			return EnsureObjectKey(userID, p, isThumb)
		}
		// 无法识别URL格式，返回空
		logger.Warn("EnsureObjectKey: 无法从URL提取对象键: %s", input)
		return ""
	}

	rel := strings.TrimPrefix(input, "/")
	if strings.HasPrefix(rel, "thumb/") {
		rel = strings.TrimPrefix(rel, "thumb/")
		isThumb = true
	}
	dir := filepath.Dir(rel)
	if dir == "." {
		dir = ""
	}
	name := filepath.Base(rel)

	// 文件名包含 _thumb. 也视为缩略图
	if !isThumb && strings.Contains(strings.ToLower(name), "_thumb.") {
		isThumb = true
	}

	// 基于 userID 构建对象键（guest=0 简化为固定目录，tenant 中已处理）
	if isThumb {
		if key, err := tenant.BuildThumbObjectKey(userID, dir, name); err == nil {
			return CleanObjectPath(key)
		}
		return CleanObjectPath(filepath.Join("thumbnails", dir, name))
	}
	if key, err := tenant.BuildObjectKey(userID, dir, name); err == nil {
		return CleanObjectPath(key)
	}
	return CleanObjectPath(filepath.Join("files", dir, name))
}
