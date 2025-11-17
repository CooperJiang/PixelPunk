package assets

import (
	"bytes"
	"embed"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed embedded/review.webp embedded/test-connect.webp embedded/fail.webp embedded/not-found.webp embedded/unauthorized.webp embedded/bandwidthLimit.webp
var defaultAssets embed.FS

// DefaultFileType 默认文件类型枚举
type DefaultFileType string

const (
	// ImageTypeReview 审核中的文件
	FileTypeReview DefaultFileType = "review"
	// ImageTypeTestConnect 存储渠道连接测试文件
	FileTypeTestConnect DefaultFileType = "test_connect"
	// ImageTypeFail 缩略图生成失败的默认文件
	FileTypeFail DefaultFileType = "fail"
	// ImageTypeNotFound 文件不存在（文件丢失）
	FileTypeNotFound DefaultFileType = "not_found"
	// ImageTypeUnauthorized 没有权限访问
	FileTypeUnauthorized DefaultFileType = "unauthorized"
	// ImageTypeError 文件加载错误
	FileTypeError DefaultFileType = "error"
	// ImageTypeBlocked 被阻止访问的文件
	FileTypeBlocked DefaultFileType = "blocked"
	// ImageTypeBandwidthLimit 带宽用尽的提示文件
	FileTypeBandwidthLimit DefaultFileType = "bandwidth_limit"
)

// DefaultFile 默认文件信息
type DefaultFile struct {
	Type        DefaultFileType
	FileName    string
	ContentType string
	CacheAge    int // 缓存时间（秒）
}

// 默认文件配置
var defaultFileConfigs = map[DefaultFileType]DefaultFile{
	FileTypeReview: {
		Type:        FileTypeReview,
		FileName:    "embedded/review.webp",
		ContentType: "image/webp",
		CacheAge:    0, // 不缓存，确保审核状态变化后立即生效
	},
	FileTypeTestConnect: {
		Type:        FileTypeTestConnect,
		FileName:    "embedded/test-connect.webp",
		ContentType: "image/webp",
		CacheAge:    3600, // 1小时缓存，测试文件可以缓存
	},
	FileTypeFail: {
		Type:        FileTypeFail,
		FileName:    "embedded/fail.webp",
		ContentType: "image/webp",
		CacheAge:    3600, // 1小时缓存，缩略图生成失败文件
	},
	FileTypeNotFound: {
		Type:        FileTypeNotFound,
		FileName:    "embedded/not-found.webp", // 文件文件不存在时显示
		ContentType: "image/webp",
		CacheAge:    86400, // 24小时
	},
	FileTypeUnauthorized: {
		Type:        FileTypeUnauthorized,
		FileName:    "embedded/unauthorized.webp", // 没有权限访问时显示
		ContentType: "image/webp",
		CacheAge:    3600, // 1小时缓存，权限可能会变化
	},
	FileTypeError: {
		Type:        FileTypeError,
		FileName:    "embedded/review.webp", // 暂时复用，后续可替换
		ContentType: "image/webp",
		CacheAge:    3600, // 1小时
	},
	FileTypeBlocked: {
		Type:        FileTypeBlocked,
		FileName:    "embedded/unauthorized.webp", // 复用权限拒绝文件
		ContentType: "image/webp",
		CacheAge:    3600, // 1小时
	},
	FileTypeBandwidthLimit: {
		Type:        FileTypeBandwidthLimit,
		FileName:    "embedded/bandwidthLimit.webp", // 带宽用尽专用提示图
		ContentType: "image/webp",
		CacheAge:    0, // 不缓存，确保带宽状态变化后立即生效
	},
}

func GetDefaultFileData(fileType DefaultFileType) ([]byte, error) {
	config, exists := defaultFileConfigs[fileType]
	if !exists {
		return nil, fmt.Errorf("未知的默认文件类型: %s", fileType)
	}

	data, err := defaultAssets.ReadFile(config.FileName)
	if err != nil {
		return nil, fmt.Errorf("读取默认文件失败: %w", err)
	}

	return data, nil
}

// ServeDefaultFile 直接向HTTP响应输出默认文件
func ServeDefaultFile(c *gin.Context, fileType DefaultFileType) {
	config, exists := defaultFileConfigs[fileType]
	if !exists {
		c.AbortWithStatusJSON(500, gin.H{"error": "未知的默认文件类型"})
		return
	}

	data, err := GetDefaultFileData(fileType)
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "无法读取默认文件"})
		return
	}

	c.Header("Content-Type", config.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	if config.CacheAge == 0 {
		// 不缓存，确保状态变化立即生效
		c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")
	} else {
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", config.CacheAge))
	}

	c.Data(http.StatusOK, config.ContentType, data)
	c.Abort() // 停止后续处理
}

func GetDefaultFileConfig(fileType DefaultFileType) (DefaultFile, bool) {
	config, exists := defaultFileConfigs[fileType]
	return config, exists
}

// RegisterDefaultFile 注册新的默认文件类型（用于扩展）
func RegisterDefaultFile(fileType DefaultFileType, config DefaultFile) {
	defaultFileConfigs[fileType] = config
}

func GetAllDefaultFileTypes() []DefaultFileType {
	types := make([]DefaultFileType, 0, len(defaultFileConfigs))
	for fileType := range defaultFileConfigs {
		types = append(types, fileType)
	}
	return types
}

func GetTestFileHeader() (*multipart.FileHeader, error) {
	data, err := GetDefaultFileData(FileTypeTestConnect)
	if err != nil {
		return nil, fmt.Errorf("读取测试文件失败: %w", err)
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "test-connect.webp")
	if err != nil {
		return nil, fmt.Errorf("创建文件字段失败: %w", err)
	}

	_, err = part.Write(data)
	if err != nil {
		return nil, fmt.Errorf("写入文件数据失败: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭writer失败: %w", err)
	}

	// 解析multipart数据
	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(32 << 20) // 32MB limit
	if err != nil {
		return nil, fmt.Errorf("读取表单失败: %w", err)
	}
	defer form.RemoveAll()

	files := form.File["file"]
	if len(files) == 0 {
		return nil, fmt.Errorf("文件字段为空")
	}

	return files[0], nil
}
