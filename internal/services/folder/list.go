package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
)

func ListFolders(userID uint, parentID string) ([]*FolderResponse, error) {
	var folders []models.Folder
	query := database.DB.Where("user_id = ?", userID)
	if parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	} else {
		query = query.Where("parent_id = '' OR parent_id IS NULL")
	}
	if err := query.Order("sort_order ASC, name ASC").Find(&folders).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹列表失败")
	}
	var result []*FolderResponse
	for i := range folders {
		result = append(result, toResponse(&folders[i]))
	}
	return result, nil
}

func ListFolderContents(userID uint, parentID string) (*FolderContentResponse, error) {
	var folders []models.Folder
	folderQuery := database.DB.Where("user_id = ?", userID)
	if parentID != "" {
		folderQuery = folderQuery.Where("parent_id = ?", parentID)
	} else {
		folderQuery = folderQuery.Where("parent_id = '' OR parent_id IS NULL")
	}
	if err := folderQuery.Order("sort_order ASC, name ASC").Find(&folders).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹列表失败")
	}
	var folderResponses []*FolderResponse
	for i := range folders {
		folderResponses = append(folderResponses, toResponse(&folders[i]))
	}

	var images []models.File
	imageQuery := database.DB.Where("user_id = ?", userID).Where("status <> ?", "pending_deletion")
	if parentID != "" {
		imageQuery = imageQuery.Where("folder_id = ?", parentID)
	} else {
		imageQuery = imageQuery.Where("folder_id = '' OR folder_id IS NULL")
	}
	if err := imageQuery.Order("created_at DESC").Find(&images).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}
	var fileResponses []*FileResponse
	for i := range images {
		file := &images[i]
		fullURL, fullThumbURL, shortURL := storage.GetFullURLs(*file)
		aiInfo, _ := getFileAIInfo(file.ID)
		fileResponses = append(fileResponses, &FileResponse{
			ID: file.ID, URL: file.URL, ThumbnailURL: file.ThumbURL, FullURL: fullURL, FullThumbURL: fullThumbURL, ShortURL: shortURL,
			OriginalName: file.OriginalName, DisplayName: file.DisplayName, Size: file.Size, Width: file.Width, Height: file.Height,
			Format: file.Format, AccessLevel: file.AccessLevel, FolderID: file.FolderID,
			CreatedAt: file.CreatedAt, UpdatedAt: file.UpdatedAt, IsDuplicate: file.IsDuplicate, MD5Hash: file.MD5Hash,
			IsRecommended: file.IsRecommended, AIInfo: aiInfo, StorageDuration: file.StorageDuration, ExpiresAt: (*common.JSONTime)(file.ExpiresAt), IsTimeLimited: file.IsTimeLimitedStorage(),
		})
	}
	return &FolderContentResponse{Folders: folderResponses, Files: fileResponses}, nil
}

/* ListFolderContentsWithoutPagination 兼容旧调用名（带过滤但不分页） */
func ListFolderContentsWithoutPagination(userID uint, parentID, keyword, accessLevel, sortBy, sortOrder string) (*FolderContentResponse, error) {
	return ListFolderContentsWithFilters(userID, parentID, keyword, 0, 0, accessLevel, sortBy, sortOrder)
}
