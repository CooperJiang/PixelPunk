package file

/* Persistence helpers split from upload_service.go (no behavior change). */

import (
	"fmt"
	"path/filepath"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/user"
	"pixelpunk/pkg/errors"
	"strings"

	"gorm.io/gorm"
)

func createFileModel(ctx *UploadContext) *models.File {
	ratio := 0.0
	if ctx.Result.Height > 0 {
		ratio = float64(ctx.Result.Width) / float64(ctx.Result.Height)
	}
	sizeFormatted := formatFileSize(ctx.File.Size)
	resolutionType := func(width, height int) string {
		pixels := width * height
		switch {
		case pixels >= 7680*4320:
			return "8K"
		case pixels >= 3840*2160:
			return "4K"
		case pixels >= 2560*1440:
			return "2K"
		case pixels >= 1920*1080:
			return "1080p"
		case pixels >= 1280*720:
			return "720p"
		case pixels >= 854*480:
			return "480p"
		default:
			return "SD"
		}
	}(ctx.Result.Width, ctx.Result.Height)
	thumbURL := ctx.Result.ThumbUrl
	if thumbURL == "" && ctx.Result.LocalThumbPath != "" {
		p := ctx.Result.LocalThumbPath
		p = strings.TrimPrefix(p, "uploads/thumbnails/")
		if strings.HasPrefix(p, "user_") {
			if idx := strings.Index(p, "/"); idx >= 0 {
				p = p[idx+1:]
			} else {
				p = ""
			}
		}
		thumbURL = p
	}
	return &models.File{
		ID:                        ctx.FileID,
		UserID:                    ctx.UserID,
		FolderID:                  ctx.FolderID,
		OriginalName:              ctx.File.Filename,
		DisplayName:               ctx.DisplayName,
		FileName:                  filepath.Base(ctx.Result.URL),
		FilePath:                  ctx.Result.URL,
		FullPath:                  ctx.Result.RemoteUrl,
		LocalFilePath:             ctx.Result.LocalUrlPath,
		LocalThumbPath:            ctx.Result.LocalThumbPath,
		URL:                       ctx.Result.URL,
		ThumbURL:                  thumbURL,
		RemoteURL:                 ctx.Result.RemoteUrl,
		RemoteThumbURL:            ctx.Result.RemoteThumbUrl,
		MD5Hash:                   ctx.FileHash,
		Size:                      ctx.File.Size,
		SizeFormatted:             sizeFormatted,
		Width:                     ctx.Result.Width,
		Height:                    ctx.Result.Height,
		Ratio:                     ratio,
		Format:                    strings.TrimPrefix(ctx.FileExt, "."),
		Mime:                      ctx.File.Header.Get("Content-Type"),
		Resolution:                resolutionType,
		Description:               getDescriptionFromContext(ctx),
		NSFW:                      false,
		AccessLevel:               ctx.AccessLevel,
		AccessKey:                 ctx.AccessKey,
		IsDuplicate:               ctx.IsDuplicate,
		OriginalFileID:            ctx.OriginalFileID,
		StorageProviderID:         ctx.StorageChannel.ID,
		StorageType:               ctx.StorageChannel.Type,
		AITaggingStatus:           "none",
		AITaggingTries:            0,
		AITaggingDuration:         0,
		StorageDuration:           ctx.StorageDuration,
		ExpiresAt:                 ctx.ExpiresAt,
		IsGuestUpload:             ctx.IsGuestUpload,
		GuestFingerprint:          ctx.GuestFingerprint,
		GuestIP:                   ctx.GuestIP,
		ThumbnailGenerationFailed: ctx.Result.ThumbnailGenerationFailed,
		ThumbnailFailureReason:    ctx.Result.ThumbnailFailureReason,
	}
}

func formatFileSize(size int64) string {
	const (
		B  = 1
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)
	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	default:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	}
}

func saveFileRecord(tx *gorm.DB, file *models.File) error {
	if err := tx.Create(file).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBCreateFailed, "保存文件记录失败")
	}
	if err := InitFileStats(tx, file.ID); err != nil {
		return err
	}
	return nil
}

func updateUserStats(tx *gorm.DB, ctx *UploadContext) error {
	return user.UpdateFileUploadStats(tx, ctx.UserID, ctx.File.Size)
}

func updateStatisticsAsync(ctx *UploadContext) {
	statsChannel <- StatsEvent{Type: "file_created", UserID: ctx.UserID, FileID: ctx.FileID, Size: ctx.File.Size}
}
