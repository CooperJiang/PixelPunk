package utils

import (
	"path/filepath"
	"strings"
)

// 文件类型常量
const (
	FileTypeImage    = "image"
	FileTypeVideo    = "video"
	FileTypeDocument = "document"
	FileTypeArchive  = "archive"
	FileTypeAudio    = "audio"
	FileTypeOther    = "other"
)

// GetFileType 根据文件名检测文件类型
func GetFileType(filename string) string {
	if filename == "" {
		return FileTypeOther
	}

	ext := strings.ToLower(filepath.Ext(filename))
	ext = strings.TrimPrefix(ext, ".")

	switch ext {
	case "jpg", "jpeg", "png", "gif", "bmp", "svg", "webp", "ico", "tiff", "tif":
		return FileTypeImage
	case "mp4", "avi", "mov", "wmv", "flv", "webm", "mkv", "m4v", "3gp", "ogv", "mpg", "mpeg":
		return FileTypeVideo
	case "pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "rtf", "odt", "ods", "odp":
		return FileTypeDocument
	case "zip", "rar", "7z", "tar", "gz", "bz2", "xz", "cab", "iso", "deb", "rpm":
		return FileTypeArchive
	case "mp3", "wav", "flac", "aac", "ogg", "wma", "m4a", "opus", "ape":
		return FileTypeAudio
	default:
		return FileTypeOther
	}
}

// GetMimeType 根据文件名获取MIME类型
func GetMimeType(filename string) string {
	if filename == "" {
		return "application/octet-stream"
	}

	ext := strings.ToLower(filepath.Ext(filename))

	mimeMap := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".svg":  "image/svg+xml",
		".webp": "image/webp",
		".ico":  "image/x-icon",
		".tiff": "image/tiff",
		".tif":  "image/tiff",

		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".wmv":  "video/x-ms-wmv",
		".flv":  "video/x-flv",
		".webm": "video/webm",
		".mkv":  "video/x-matroska",
		".m4v":  "video/x-m4v",
		".3gp":  "video/3gpp",
		".ogv":  "video/ogg",
		".mpg":  "video/mpeg",
		".mpeg": "video/mpeg",

		".pdf":  "application/pdf",
		".doc":  "application/msword",
		".docx": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		".xls":  "application/vnd.ms-excel",
		".xlsx": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		".ppt":  "application/vnd.ms-powerpoint",
		".pptx": "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		".txt":  "text/plain",
		".rtf":  "application/rtf",
		".odt":  "application/vnd.oasis.opendocument.text",
		".ods":  "application/vnd.oasis.opendocument.spreadsheet",
		".odp":  "application/vnd.oasis.opendocument.presentation",

		".zip": "application/zip",
		".rar": "application/vnd.rar",
		".7z":  "application/x-7z-compressed",
		".tar": "application/x-tar",
		".gz":  "application/gzip",
		".bz2": "application/x-bzip2",
		".xz":  "application/x-xz",
		".cab": "application/vnd.ms-cab-compressed",
		".iso": "application/x-iso9660-image",

		".mp3":  "audio/mpeg",
		".wav":  "audio/wav",
		".flac": "audio/flac",
		".aac":  "audio/aac",
		".ogg":  "audio/ogg",
		".wma":  "audio/x-ms-wma",
		".m4a":  "audio/x-m4a",
		".opus": "audio/opus",
		".ape":  "audio/x-ape",
	}

	if mime, exists := mimeMap[ext]; exists {
		return mime
	}

	return "application/octet-stream"
}

// IsImageFile 检查是否为文件文件
func IsImageFile(filename string) bool {
	return GetFileType(filename) == FileTypeImage
}

// IsVideoFile 检查是否为视频文件
func IsVideoFile(filename string) bool {
	return GetFileType(filename) == FileTypeVideo
}

// IsDocumentFile 检查是否为文档文件
func IsDocumentFile(filename string) bool {
	return GetFileType(filename) == FileTypeDocument
}

// IsArchiveFile 检查是否为压缩包文件
func IsArchiveFile(filename string) bool {
	return GetFileType(filename) == FileTypeArchive
}

// IsAudioFile 检查是否为音频文件
func IsAudioFile(filename string) bool {
	return GetFileType(filename) == FileTypeAudio
}

// HasDimensions 检查文件类型是否可能有尺寸信息
func HasDimensions(fileType string) bool {
	return fileType == FileTypeImage || fileType == FileTypeVideo
}

func GetFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if len(ext) > 0 {
		return strings.ToLower(ext[1:]) // 去掉开头的点号
	}
	return ""
}
