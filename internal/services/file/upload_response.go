package file

/* Response builders split from upload_service.go (no behavior change). */

import (
	"fmt"
	"pixelpunk/pkg/storage"
	"time"
)

func createFileResponse(ctx *UploadContext) *FileDetailResponse {
	fullURL, fullThumbURL, _ := storage.GetFullURLs(*ctx.SavedFile)
	return &FileDetailResponse{
		ID:                ctx.SavedFile.ID,
		FullURL:           fullURL,
		FullThumbURL:      fullThumbURL,
		OriginalName:      ctx.SavedFile.OriginalName,
		DisplayName:       ctx.SavedFile.DisplayName,
		Size:              ctx.SavedFile.Size,
		Width:             ctx.SavedFile.Width,
		Height:            ctx.SavedFile.Height,
		Format:            ctx.SavedFile.Format,
		AccessLevel:       ctx.SavedFile.AccessLevel,
		FolderID:          ctx.SavedFile.FolderID,
		CreatedAt:         ctx.SavedFile.CreatedAt,
		UpdatedAt:         ctx.SavedFile.UpdatedAt,
		IsDuplicate:       ctx.SavedFile.IsDuplicate,
		MD5Hash:           ctx.SavedFile.MD5Hash,
		IsRecommended:     ctx.SavedFile.IsRecommended,
		StorageProviderID: ctx.SavedFile.StorageProviderID,
	}
}

func createExternalAPIResponse(ctx *UploadContext) *ExternalAPIFileResponse {
	fullURL, fullThumbURL, _ := storage.GetFullURLs(*ctx.SavedFile)
	return &ExternalAPIFileResponse{
		ID:           ctx.SavedFile.ID,
		URL:          fullURL,
		ThumbURL:     fullThumbURL,
		OriginalName: ctx.SavedFile.OriginalName,
		Size:         ctx.SavedFile.Size,
		Width:        ctx.SavedFile.Width,
		Height:       ctx.SavedFile.Height,
		Format:       ctx.SavedFile.Format,
		AccessLevel:  ctx.SavedFile.AccessLevel,
		CreatedAt:    ctx.SavedFile.CreatedAt,
	}
}

func buildUploadResponse(ctx *UploadContext) *FileUploadResponse {
	var fullURL, fullThumbURL, shortLinkURL string
	var createdAt, updatedAt *time.Time
	var md5Hash string
	var isRecommended bool
	if ctx.SavedFile != nil {
		fullURL, fullThumbURL, shortLinkURL = storage.GetFullURLs(*ctx.SavedFile)
		createdAt = (*time.Time)(&ctx.SavedFile.CreatedAt)
		updatedAt = (*time.Time)(&ctx.SavedFile.UpdatedAt)
		md5Hash = ctx.SavedFile.MD5Hash
		isRecommended = ctx.SavedFile.IsRecommended
	}
	var shortURL, imageIDURL, thumbIDURL string
	if ctx.SavedFile != nil {
		shortURL = ctx.SavedFile.GenerateShortURL()
		imageIDURL = fmt.Sprintf("/f/%s", ctx.SavedFile.ID)
		thumbIDURL = fmt.Sprintf("/t/%s", ctx.SavedFile.ID)
	}

	response := &FileUploadResponse{
		ID:           ctx.FileID,
		URL:          ctx.Result.URL,
		FullURL:      fullURL,
		FullThumbURL: fullThumbURL,

		ShortURL:     shortURL,
		FileIDURL:    imageIDURL,
		ThumbIDURL:   thumbIDURL,
		ShortLinkURL: shortLinkURL,

		Width:                     ctx.Result.Width,
		Height:                    ctx.Result.Height,
		Size:                      ctx.FileSize,
		SizeFormatted:             fmt.Sprintf("%.2f KB", float64(ctx.FileSize)/1024),
		Ratio:                     float64(ctx.Result.Width) / float64(ctx.Result.Height),
		Format:                    ctx.FileFormat,
		OriginalName:              ctx.OriginalName,
		DisplayName:               ctx.DisplayName,
		AccessLevel:               ctx.AccessLevel,
		AccessKey:                 ctx.AccessKey,
		IsDuplicate:               ctx.IsDuplicate,
		OriginalFileID:            ctx.OriginalFileID,
		StorageProvider:           ctx.ActualChannelID,
		StorageProviderID:         ctx.ActualChannelID,
		StorageChannel:            ctx.StorageChannel,
		MD5Hash:                   md5Hash,
		IsRecommended:             isRecommended,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
		StorageDuration:           ctx.StorageDuration,
		ExpiresAt:                 ctx.ExpiresAt,
		IsGuestUpload:             ctx.IsGuestUpload,
		ThumbnailGenerationFailed: ctx.Result.ThumbnailGenerationFailed,
		ThumbnailFailureReason:    ctx.Result.ThumbnailFailureReason,
		WatermarkApplied:          ctx.WatermarkApplied,
		WatermarkFailureReason:    ctx.WatermarkFailureReason,
	}
	if ctx.Result.Height == 0 || response.Ratio != response.Ratio {
		response.Ratio = 1.0
	}
	return response
}
