package file

/* Common response builders to reduce duplication (no behavior change). */

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/storage"
)

/* BuildFileDetailResponse 统一构建用户侧文件响应 */
func BuildFileDetailResponse(file models.File, views int64, aiInfo *AIInfoResponse) FileDetailResponse {
	fullURL, fullThumbURL, _ := storage.GetFullURLs(file)
	return FileDetailResponse{
		ID:                file.ID,
		FullURL:           fullURL,
		FullThumbURL:      fullThumbURL,
		OriginalName:      file.OriginalName,
		DisplayName:       file.DisplayName,
		Size:              file.Size,
		Width:             file.Width,
		Height:            file.Height,
		Format:            file.Format,
		AccessLevel:       file.AccessLevel,
		FolderID:          file.FolderID,
		CreatedAt:         file.CreatedAt,
		UpdatedAt:         file.UpdatedAt,
		Views:             views,
		IsDuplicate:       file.IsDuplicate,
		MD5Hash:           file.MD5Hash,
		IsRecommended:     file.IsRecommended,
		StorageProviderID: file.StorageProviderID,
		AIInfo:            aiInfo,
	}
}

/* BuildAdminFileDetailResponse 统一构建管理员侧文件响应 */
func BuildAdminFileDetailResponse(file models.File, views int64, userName string, aiInfo *AIInfoResponse) AdminFileDetailResponse {
	fullURL, fullThumbURL, shortURL := storage.GetFullURLs(file)
	return AdminFileDetailResponse{
		ID:                file.ID,
		URL:               file.URL,
		ThumbnailURL:      file.ThumbURL,
		FullURL:           fullURL,
		FullThumbURL:      fullThumbURL,
		ShortURL:          shortURL,
		OriginalName:      file.OriginalName,
		DisplayName:       file.DisplayName,
		Size:              file.Size,
		Width:             file.Width,
		Height:            file.Height,
		Format:            file.Format,
		AccessLevel:       file.AccessLevel,
		FolderID:          file.FolderID,
		CreatedAt:         file.CreatedAt,
		UpdatedAt:         file.UpdatedAt,
		Views:             views,
		IsRecommended:     file.IsRecommended,
		StorageProviderID: file.StorageProviderID,
		IsDuplicate:       file.IsDuplicate,
		MD5Hash:           file.MD5Hash,
		UserID:            file.UserID,
		UserName:          userName,
		AIInfo:            aiInfo,
		StorageDuration:   file.StorageDuration,
		ExpiresAt:         (*common.JSONTime)(file.ExpiresAt),
		IsTimeLimited:     file.IsTimeLimitedStorage(),
	}
}
