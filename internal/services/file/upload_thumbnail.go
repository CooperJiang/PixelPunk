package file

/* Thumbnail capture helpers split from upload_service.go (no behavior change). */

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	pkgStorage "pixelpunk/pkg/storage"
	pathutil "pixelpunk/pkg/storage/path"
	"strings"
)

func captureThumbnailBase64(ctx *UploadContext) error {
	if ctx.ActualChannelID != "" {
		return captureThumbnailBase64FromStorage(ctx)
	}
	return captureThumbnailBase64FromLocalFile(ctx)
}

func captureThumbnailBase64FromStorage(ctx *UploadContext) error {
	storageService := pkgStorage.NewGlobalStorage()
	var thumbPath string

	if ctx.Result.ThumbUrl != "" {
		thumbPath = ctx.Result.ThumbUrl
		if pathutil.IsHTTPURL(thumbPath) {
			if extractedPath := pathutil.ExtractObjectPathFromURL(thumbPath); extractedPath != "" {
				thumbPath = extractedPath
			}
		}
	} else if ctx.Result.LocalThumbPath != "" {
		p := ctx.Result.LocalThumbPath
		p = strings.TrimPrefix(p, "uploads/thumbnails/")
		thumbPath = p
	} else if ctx.Result.URL != "" {
		thumbPath = ctx.Result.URL
	} else {
		return fmt.Errorf("无法获取文件路径进行Base64读取")
	}

	thumbPath = ensureFullObjectPath(ctx, thumbPath)

	ctxBg := context.Background()
	base64Data, err := storageService.GetThumbnailBase64(ctxBg, ctx.ActualChannelID, thumbPath)
	if err != nil {
		if ctx.Result.URL != "" {
			originalPath := ctx.Result.URL
			if pathutil.IsHTTPURL(originalPath) {
				if extractedPath := pathutil.ExtractObjectPathFromURL(originalPath); extractedPath != "" {
					originalPath = extractedPath
				}
			}
			originalPath = ensureFullObjectPath(ctx, originalPath)
			base64Data, err = storageService.GetBase64(ctxBg, ctx.ActualChannelID, originalPath)
		}
		if err != nil {
			return fmt.Errorf("存储适配器读取Base64失败: %v", err)
		}
	}
	ctx.ThumbnailBase64 = base64Data
	return nil
}

func captureThumbnailBase64FromLocalFile(ctx *UploadContext) error {
	var filePath string
	if ctx.Result.LocalThumbPath != "" {
		filePath = ctx.Result.LocalThumbPath
	} else if ctx.Result.LocalUrlPath != "" {
		filePath = ctx.Result.LocalUrlPath
	} else {
		return fmt.Errorf("无法找到可用的本地文件路径")
	}
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("文件不存在: %s", filePath)
	}
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("读取文件内容失败: %v", err)
	}
	ctx.ThumbnailBase64 = base64.StdEncoding.EncodeToString(data)
	return nil
}
