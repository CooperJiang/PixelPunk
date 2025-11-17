package ai

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"pixelpunk/internal/models"
	sadapter "pixelpunk/pkg/storage/adapter"
)

func (s *TaggingService) processImageWithBase64(file models.File) error {
	base64Data, imageFormat, err := s.readImageAsBase64(file)

	if err != nil {
		return fmt.Errorf("读取文件文件失败: %v", err)
	}

	imageData := &CompressedFileData{
		ThumbnailBase64: base64Data,
	}

	err = AiImageTaggingAndSaveWithCompressedFile(file, imageData, imageFormat)

	return err
}

// readImageAsBase64 使用存储适配器读取文件文件并转换为Base64
func (s *TaggingService) readImageAsBase64(file models.File) (string, string, error) {
	ctx := context.Background()

	if file.StorageProviderID == "" {
		return "", "", fmt.Errorf("文件缺少存储渠道ID: %s", file.ID)
	}

	// 优先尝试缩略图，如果没有则使用原图
	var base64Data string
	var err error
	var useThumb bool = false
	var readAttempts []string // 记录尝试的路径

	// 直接使用物理地址字段，不需要复杂的路径拼接
	var thumbPath, originalPath string

	// 使用LocalThumbPath作为缩略图的物理路径（本地存储时是本地路径，三方存储时是三方物理路径）
	if file.LocalThumbPath != "" {
		thumbPath = file.LocalThumbPath
	}

	// 使用LocalFilePath作为原图的物理路径
	if file.LocalFilePath != "" {
		if file.StorageType == "local" {
			if strings.HasPrefix(file.LocalFilePath, "uploads/files/") {
				originalPath = strings.TrimPrefix(file.LocalFilePath, "uploads/files/")
			} else {
				originalPath = file.LocalFilePath
			}
		} else {
			originalPath = file.LocalFilePath
		}
	}

	// 本地路径归一化辅助：将物理路径映射为对象键，避免本地适配器二次拼接导致路径重复
	normalizeLocalObjectKey := func(p string, isThumb bool) string {
		if p == "" {
			return ""
		}
		clean := strings.TrimPrefix(p, "/")

		if strings.HasPrefix(clean, "files/") || strings.HasPrefix(clean, "thumbnails/") {
			return clean
		}

		if isThumb {
			if strings.HasPrefix(clean, "uploads/thumbnails/") {
				rel := strings.TrimPrefix(clean, "uploads/thumbnails/")
				return "thumbnails/" + rel
			} else {
				return "thumbnails/" + clean
			}
		} else {
			if strings.HasPrefix(clean, "uploads/files/") {
				rel := strings.TrimPrefix(clean, "uploads/files/")
				return "files/" + rel
			} else if strings.HasPrefix(clean, "uploads/images1/") {
				rel := strings.TrimPrefix(clean, "uploads/images1/")
				return "files/" + rel
			} else {
				return "files/" + clean
			}
		}
	}

	if thumbPath != "" {
		readAttempts = append(readAttempts, fmt.Sprintf("缩略图路径: %s", thumbPath))
		// 对 local 做对象键归一化，避免重复前缀
		normalized := thumbPath
		if file.StorageType == "local" {
			normalized = normalizeLocalObjectKey(thumbPath, true)
		}
		base64Data, err = s.storage.GetThumbnailBase64(ctx, file.StorageProviderID, normalized)
		if err == nil && len(base64Data) > 0 {
			useThumb = true
		} else {
			useThumb = false
			base64Data = ""
		}
	}

	// 回退原图（仅在缩略图失败且原图路径存在时）
	if !useThumb && originalPath != "" {
		readAttempts = append(readAttempts, fmt.Sprintf("原图路径: %s", originalPath))
		normalized := originalPath
		if file.StorageType == "local" {
			normalized = normalizeLocalObjectKey(originalPath, false)
		}
		base64Data, err = s.storage.GetBase64(ctx, file.StorageProviderID, normalized)
		if err != nil {
			attemptInfo := strings.Join(readAttempts, ", ")
			if file.StorageType == "local" {
				candidates := []string{}
				if file.LocalThumbPath != "" {
					candidates = append(candidates, file.LocalThumbPath)
				}
				if file.LocalFilePath != "" {
					candidates = append(candidates, file.LocalFilePath)
				}
				for _, c := range []string{file.LocalThumbPath, file.LocalFilePath} {
					if strings.HasPrefix(c, "uploads/images1/") {
						candidates = append(candidates, "uploads/files/"+strings.TrimPrefix(c, "uploads/images1/"))
					}
				}
				for _, p := range candidates {
					if strings.TrimSpace(p) == "" {
						continue
					}
					if data, ferr := s.readLocalFileAsBase64(p); ferr == nil && data != "" {
						base64Data = data
						err = nil
						break
					}
				}
				if err == nil && base64Data != "" {
				} else {
					// 若读取错误显式为 not found，则返回带哨兵错误，便于上层策略处理
					if sadapter.IsNotFoundError(err) {
						return "", "", fmt.Errorf("missing file: %w", errMissingFile)
					}
					return "", "", fmt.Errorf("文件读取失败 - 存储类型: %s, 存储渠道ID: %s, 尝试路径: [%s], 最终错误: %w",
						file.StorageType, file.StorageProviderID, attemptInfo, err)
				}
			} else {
				if sadapter.IsNotFoundError(err) {
					return "", "", fmt.Errorf("missing file: %w", errMissingFile)
				}
				return "", "", fmt.Errorf("文件读取失败 - 存储类型: %s, 存储渠道ID: %s, 尝试路径: [%s], 最终错误: %w",
					file.StorageType, file.StorageProviderID, attemptInfo, err)
			}
		}
	}

	if thumbPath == "" && originalPath == "" {
		return "", "", fmt.Errorf("文件文件路径信息缺失 - 缩略图路径: '%s', 原图路径: '%s'",
			file.LocalThumbPath, file.LocalFilePath)
	}

	// 根据实际使用的文件路径确定 MIME 扩展
	var imageFormat string
	if useThumb {
		imageFormat = strings.TrimPrefix(strings.ToLower(filepath.Ext(thumbPath)), ".")
	} else {
		imageFormat = strings.TrimPrefix(strings.ToLower(filepath.Ext(originalPath)), ".")
	}
	// 回退：数据表格式 -> jpg
	if strings.TrimSpace(imageFormat) == "" {
		if strings.TrimSpace(file.Format) != "" {
			imageFormat = strings.TrimPrefix(strings.ToLower(file.Format), ".")
		} else {
			imageFormat = "jpg"
		}
	}

	// 最终检查：确保获得了有效的base64数据
	if len(base64Data) == 0 {
		attemptInfo := strings.Join(readAttempts, ", ")
		return "", "", fmt.Errorf("获取的Base64数据为空 - 存储类型: %s, 存储渠道ID: %s, 尝试路径: [%s]",
			file.StorageType, file.StorageProviderID, attemptInfo)
	}

	return base64Data, imageFormat, nil
}

// readLocalFileAsBase64 直接从本地文件读取并转为base64（用于本地适配器兜底）
func (s *TaggingService) readLocalFileAsBase64(p string) (string, error) {
	if strings.TrimSpace(p) == "" {
		return "", fmt.Errorf("empty path")
	}
	// 允许相对路径（以当前工作目录为根）
	f, err := os.Open(p)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// ProcessSingleFile 处理单张文件
