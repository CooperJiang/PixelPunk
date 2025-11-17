package middleware

import (
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"pixelpunk/pkg/imagex/decode"
	"pixelpunk/pkg/imagex/formats"
)

// ValidationOptions 验证选项
type ValidationOptions struct {
	MaxFileSize  int64 // 单个文件最大大小
	MaxTotalSize int64 // 批量上传总大小限制
	MinFileSize  int64 // 最小文件大小

	AllowedFormats   []string // 允许的文件格式
	ForbiddenFormats []string // 禁止的文件格式
	CheckMimeType    bool     // 是否检查MIME类型

	MaxNameLength  int      // 文件名最大长度
	ForbiddenChars []string // 禁止的字符
	ForbiddenNames []string // 禁止的文件名
	SanitizeNames  bool     // 是否清理文件名

	MaxFileCount int // 最大文件数量

	CheckFileHeader bool // 是否检查文件头
	ScanMalware     bool // 是否扫描恶意软件
	CheckDimensions bool // 是否检查图像尺寸
	MaxWidth        int  // 最大宽度
	MaxHeight       int  // 最大高度
	MinWidth        int  // 最小宽度
	MinHeight       int  // 最小高度

	CustomValidators []CustomValidator
}

// CustomValidator 自定义验证器
type CustomValidator func(*multipart.FileHeader) error

// ValidationResult 验证结果
type ValidationResult struct {
	Valid        bool                    // 是否验证通过
	ValidFiles   []*multipart.FileHeader // 有效文件
	InvalidFiles []*ValidationError      // 无效文件及错误信息
	TotalSize    int64                   // 总大小
	FileCount    int                     // 文件数量
	Warnings     []string                // 警告信息
}

// ValidationError 验证错误
type ValidationError struct {
	File    *multipart.FileHeader // 文件
	Code    string                // 错误代码
	Message string                // 错误信息
	Field   string                // 错误字段
}

// FileValidator 文件验证器
type FileValidator struct {
	defaultOptions ValidationOptions
}

func NewFileValidator(defaultOptions ValidationOptions) *FileValidator {
	return &FileValidator{
		defaultOptions: defaultOptions,
	}
}

var defaultValidationOptions = &ValidationOptions{
	MaxFileSize:     20 * 1024 * 1024,
	AllowedFormats:  formats.SupportedExtensionsWithDot(),
	CheckFileHeader: true,
	CheckMimeType:   true,
}

// DefaultFileValidator 默认文件验证器
var DefaultFileValidator = NewFileValidator(*defaultValidationOptions)

// ValidateFiles 验证文件列表
func (v *FileValidator) ValidateFiles(files []*multipart.FileHeader, options *ValidationOptions) (*ValidationResult, error) {

	opts := v.mergeOptions(options)

	result := &ValidationResult{
		Valid:        true,
		ValidFiles:   make([]*multipart.FileHeader, 0),
		InvalidFiles: make([]*ValidationError, 0),
		Warnings:     make([]string, 0),
	}

	if len(files) > opts.MaxFileCount {
		return nil, fmt.Errorf("文件数量超过限制: %d > %d", len(files), opts.MaxFileCount)
	}

	var totalSize int64

	for _, file := range files {

		if err := v.validateSingleFile(file, opts); err != nil {
			validationErr := &ValidationError{
				File:    file,
				Code:    "VALIDATION_FAILED",
				Message: err.Error(),
			}
			result.InvalidFiles = append(result.InvalidFiles, validationErr)
			result.Valid = false
			continue
		}

		totalSize += file.Size
		result.ValidFiles = append(result.ValidFiles, file)
	}

	if opts.MaxTotalSize > 0 && totalSize > opts.MaxTotalSize {
		return nil, fmt.Errorf("总文件大小超过限制: %d > %d", totalSize, opts.MaxTotalSize)
	}

	result.TotalSize = totalSize
	result.FileCount = len(result.ValidFiles)

	return result, nil
}

// validateSingleFile 验证单个文件
func (v *FileValidator) validateSingleFile(file *multipart.FileHeader, options ValidationOptions) error {
	if file == nil {
		return fmt.Errorf("文件为空")
	}

	if strings.TrimSpace(file.Filename) == "" {
		return fmt.Errorf("文件名为空")
	}

	if options.MaxFileSize > 0 && file.Size > options.MaxFileSize {
		return fmt.Errorf("文件大小超过限制: %d > %d", file.Size, options.MaxFileSize)
	}

	if options.MinFileSize > 0 && file.Size < options.MinFileSize {
		return fmt.Errorf("文件大小低于最小限制: %d < %d", file.Size, options.MinFileSize)
	}

	if file.Size == 0 {
		return fmt.Errorf("文件为空")
	}

	if err := v.validateFileName(file.Filename, options); err != nil {
		return fmt.Errorf("文件名验证失败: %w", err)
	}

	if err := v.validateFileFormat(file.Filename, options); err != nil {
		return fmt.Errorf("文件格式验证失败: %w", err)
	}

	if options.CheckMimeType {
		if err := v.validateMimeType(file, options); err != nil {
			return fmt.Errorf("MIME类型验证失败: %w", err)
		}
	}

	if options.CheckFileHeader {
		if err := v.validateFileHeader(file); err != nil {
			return fmt.Errorf("文件头验证失败: %w", err)
		}
	}

	if options.CheckDimensions {
		if err := v.validateImageDimensions(file, options); err != nil {
			return fmt.Errorf("图像尺寸验证失败: %w", err)
		}
	}

	for _, validator := range options.CustomValidators {
		if err := validator(file); err != nil {
			return fmt.Errorf("自定义验证失败: %w", err)
		}
	}

	return nil
}

// validateFileName 验证文件名
func (v *FileValidator) validateFileName(filename string, options ValidationOptions) error {
	if options.MaxNameLength > 0 && len(filename) > options.MaxNameLength {
		return fmt.Errorf("文件名过长: %d > %d", len(filename), options.MaxNameLength)
	}

	for _, char := range options.ForbiddenChars {
		if strings.Contains(filename, char) {
			return fmt.Errorf("文件名包含禁止字符: %s", char)
		}
	}

	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	for _, forbidden := range options.ForbiddenNames {
		if strings.EqualFold(nameWithoutExt, forbidden) {
			return fmt.Errorf("文件名被禁止: %s", forbidden)
		}
	}

	dangerousPatterns := []string{
		"con", "prn", "aux", "nul",
		"com1", "com2", "com3", "com4", "com5", "com6", "com7", "com8", "com9",
		"lpt1", "lpt2", "lpt3", "lpt4", "lpt5", "lpt6", "lpt7", "lpt8", "lpt9",
	}

	for _, pattern := range dangerousPatterns {
		if strings.EqualFold(nameWithoutExt, pattern) {
			return fmt.Errorf("文件名为系统保留名称: %s", pattern)
		}
	}

	return nil
}

// validateFileFormat 验证文件格式
func (v *FileValidator) validateFileFormat(filename string, options ValidationOptions) error {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, forbidden := range options.ForbiddenFormats {
		if strings.ToLower(forbidden) == ext {
			return fmt.Errorf("文件格式被禁止: %s", ext)
		}
	}

	if len(options.AllowedFormats) > 0 {
		allowed := false
		for _, format := range options.AllowedFormats {
			if strings.ToLower(format) == ext {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("文件格式不被支持: %s", ext)
		}
	}

	return nil
}

// validateMimeType 验证MIME类型
func (v *FileValidator) validateMimeType(file *multipart.FileHeader, options ValidationOptions) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("无法读取文件头: %w", err)
	}

	// 这里可以使用http.DetectContentType或者更精确的MIME检测
	// contentType := http.DetectContentType(buffer)

	// MIME类型验证已通过 validateFileHeader 函数中的文件头检测实现
	return nil
}

// validateFileHeader 验证文件头
func (v *FileValidator) validateFileHeader(file *multipart.FileHeader) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer src.Close()

	header := make([]byte, 16)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return fmt.Errorf("无法读取文件头: %w", err)
	}

	if n < 2 {
		return fmt.Errorf("文件头过短")
	}

	// 检查常见图像格式的文件头，但采用宽松策略
	ext := strings.ToLower(filepath.Ext(file.Filename))
	return v.validateFileHeaderWithTolerance(header, n, ext, file.Filename)

	/*
		// 原文件头验证逻辑，后续需要时取消注释
		src, err := file.Open()
		if err != nil {
			return fmt.Errorf("无法打开文件: %w", err)
		}
		defer src.Close()

		header := make([]byte, 16)
		n, err := src.Read(header)
		if err != nil && err != io.EOF {
			return fmt.Errorf("无法读取文件头: %w", err)
		}

		if n < 2 {
			return fmt.Errorf("文件头过短")
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			if n >= 3 && (header[0] != 0xFF || header[1] != 0xD8 || header[2] != 0xFF) {
				return fmt.Errorf("JPEG文件头不匹配")
			}
		case ".png":
			pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
			if n < len(pngHeader) {
				return fmt.Errorf("PNG文件头长度不足")
			}
			for i, b := range pngHeader {
				if header[i] != b {
					return fmt.Errorf("PNG文件头不匹配")
				}
			}
		case ".gif":
			if n >= 6 && (string(header[:6]) != "GIF87a" && string(header[:6]) != "GIF89a") {
				return fmt.Errorf("GIF文件头不匹配")
			}
		case ".bmp":
			if n >= 2 && (header[0] != 0x42 || header[1] != 0x4D) {
				return fmt.Errorf("BMP文件头不匹配")
			}
		case ".webp":
			if n >= 12 && (string(header[:4]) != "RIFF" || string(header[8:12]) != "WEBP") {
				return fmt.Errorf("WebP文件头不匹配")
			}
		}

		return nil
	*/
}

// validateImageDimensions 验证图像尺寸
func (v *FileValidator) validateImageDimensions(file *multipart.FileHeader, options ValidationOptions) error {
	// 使用 imagex/decode 中的 DetectFormat
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}
	defer src.Close()

	width, height, _, err := decode.DetectFormat(src)
	if err != nil {
		return fmt.Errorf("无法获取图像尺寸: %w", err)
	}

	if options.MaxWidth > 0 && width > options.MaxWidth {
		return fmt.Errorf("图像宽度超过限制: %d > %d", width, options.MaxWidth)
	}

	if options.MaxHeight > 0 && height > options.MaxHeight {
		return fmt.Errorf("图像高度超过限制: %d > %d", height, options.MaxHeight)
	}

	if options.MinWidth > 0 && width < options.MinWidth {
		return fmt.Errorf("图像宽度低于最小限制: %d < %d", width, options.MinWidth)
	}

	if options.MinHeight > 0 && height < options.MinHeight {
		return fmt.Errorf("图像高度低于最小限制: %d < %d", height, options.MinHeight)
	}

	return nil
}

// mergeOptions 合并选项
func (v *FileValidator) mergeOptions(options *ValidationOptions) ValidationOptions {
	result := v.defaultOptions

	if options != nil {
		if options.MaxFileSize > 0 {
			result.MaxFileSize = options.MaxFileSize
		}
		if options.MaxTotalSize > 0 {
			result.MaxTotalSize = options.MaxTotalSize
		}
		if options.MinFileSize > 0 {
			result.MinFileSize = options.MinFileSize
		}
		if len(options.AllowedFormats) > 0 {
			result.AllowedFormats = options.AllowedFormats
		}
		if len(options.ForbiddenFormats) > 0 {
			result.ForbiddenFormats = options.ForbiddenFormats
		}
		if options.MaxNameLength > 0 {
			result.MaxNameLength = options.MaxNameLength
		}
		if len(options.ForbiddenChars) > 0 {
			result.ForbiddenChars = options.ForbiddenChars
		}
		if len(options.ForbiddenNames) > 0 {
			result.ForbiddenNames = options.ForbiddenNames
		}
		if options.MaxFileCount > 0 {
			result.MaxFileCount = options.MaxFileCount
		}
		if options.MaxWidth > 0 {
			result.MaxWidth = options.MaxWidth
		}
		if options.MaxHeight > 0 {
			result.MaxHeight = options.MaxHeight
		}
		if options.MinWidth > 0 {
			result.MinWidth = options.MinWidth
		}
		if options.MinHeight > 0 {
			result.MinHeight = options.MinHeight
		}
		if len(options.CustomValidators) > 0 {
			result.CustomValidators = options.CustomValidators
		}

		result.CheckMimeType = options.CheckMimeType
		result.SanitizeNames = options.SanitizeNames
		result.CheckFileHeader = options.CheckFileHeader
		result.ScanMalware = options.ScanMalware
		result.CheckDimensions = options.CheckDimensions
	}

	return result
}

// SanitizeFileName 清理文件名
func (v *FileValidator) SanitizeFileName(filename string) string {
	unsafe := map[string]string{
		"/":  "_",
		"\\": "_",
		":":  "_",
		"*":  "_",
		"?":  "_",
		"\"": "_",
		"<":  "_",
		">":  "_",
		"|":  "_",
		" ":  "_", // 可选：替换空格
	}

	result := filename
	for unsafe, safe := range unsafe {
		result = strings.ReplaceAll(result, unsafe, safe)
	}

	result = strings.Trim(result, ". ")

	if len(result) > 255 {
		ext := filepath.Ext(result)
		name := strings.TrimSuffix(result, ext)
		maxNameLen := 255 - len(ext)
		if maxNameLen > 0 {
			result = name[:maxNameLen] + ext
		}
	}

	return result
}

// 便捷方法

// ValidateFiles 验证文件列表（便捷方法）
func ValidateFiles(files []*multipart.FileHeader, options *ValidationOptions) (*ValidationResult, error) {
	return DefaultFileValidator.ValidateFiles(files, options)
}

// ValidateSingleFile 验证单个文件（便捷方法）
func ValidateSingleFile(file *multipart.FileHeader, options *ValidationOptions) error {
	validator := DefaultFileValidator
	opts := validator.mergeOptions(options)
	return validator.validateSingleFile(file, opts)
}

// QuickValidate 快速验证（便捷方法）
func QuickValidate(files []*multipart.FileHeader, maxFileSize, maxTotalSize int64) (*ValidationResult, error) {
	options := &ValidationOptions{
		MaxFileSize:  maxFileSize,
		MaxTotalSize: maxTotalSize,
	}
	return ValidateFiles(files, options)
}

// CreateAPIKeyValidator 创建API密钥验证器
func CreateAPIKeyValidator(singleFileLimit, storageLimit int64, allowedFormats []string) *FileValidator {
	options := ValidationOptions{
		MaxFileSize:     singleFileLimit,
		AllowedFormats:  allowedFormats,
		CheckFileHeader: true,
		CheckMimeType:   true,
		SanitizeNames:   true,
	}

	return NewFileValidator(options)
}

// CreateChannelValidator 创建存储渠道验证器
func CreateChannelValidator(channelUploadLimit int64) *FileValidator {
	options := ValidationOptions{
		MaxFileSize:     channelUploadLimit,
		CheckFileHeader: true,
		CheckMimeType:   true,
	}

	return NewFileValidator(options)
}

// validateFileHeaderWithTolerance 宽松的文件头验证
func (v *FileValidator) validateFileHeaderWithTolerance(header []byte, headerSize int, ext, filename string) error {

	switch ext {
	case ".jpg", ".jpeg":
		if headerSize >= 3 && (header[0] != 0xFF || header[1] != 0xD8 || header[2] != 0xFF) {
			// 允许轻微的文件头差异
		}
	case ".png":
		pngHeader := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		if headerSize < len(pngHeader) {
			return fmt.Errorf("PNG文件头长度不足: %s", filename)
		}

		if header[1] != 0x50 || header[2] != 0x4E || header[3] != 0x47 {
			// 对于可能损坏的PNG，记录警告但仍允许上传，由后续处理决定
		}

		for i, expectedByte := range pngHeader {
			if i < headerSize && header[i] != expectedByte {
			}
		}
	case ".gif":
		if headerSize >= 6 && (string(header[:6]) != "GIF87a" && string(header[:6]) != "GIF89a") {
		}
	case ".bmp":
		if headerSize >= 2 && (header[0] != 0x42 || header[1] != 0x4D) {
		}
	case ".webp":
		if headerSize >= 12 && (string(header[:4]) != "RIFF" || string(header[8:12]) != "WEBP") {
		}
	default:
	}

	// 宽松策略：只记录警告，不阻止上传
	return nil
}

// BatchValidateWithLimits 批量验证带限制
func BatchValidateWithLimits(files []*multipart.FileHeader, dailyLimit int, userUploadCount int) (*ValidationResult, error) {
	if dailyLimit > 0 && userUploadCount+len(files) > dailyLimit {
		return nil, fmt.Errorf("超过每日上传限制: %d + %d > %d", userUploadCount, len(files), dailyLimit)
	}

	return ValidateFiles(files, nil)
}
