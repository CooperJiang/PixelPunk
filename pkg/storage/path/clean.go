package path

import (
	"path"
	"strings"
)

// CleanObjectPath 规范化对象键（云存储对象路径）
// - 统一分隔符为 '/'
// - 去掉多余的 '/'（包含前后缀）
// - 保持相对路径（不以 '/' 开头）
// - 不做大小写与编码转换
func CleanObjectPath(p string) string {
	if p == "" {
		return ""
	}
	p = strings.ReplaceAll(p, "\\", "/")
	// 使用 path.Clean 去除重复与 . ..
	p = path.Clean(p)
	// path.Clean("//a//b/") -> "/a/b"，去掉前导 '/'
	p = strings.TrimPrefix(p, "/")
	// 移除多余的尾部 '/'
	if p == "." {
		return ""
	}
	return strings.TrimSuffix(p, "/")
}
