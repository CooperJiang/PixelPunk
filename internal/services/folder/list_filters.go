package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/storage"
)

func ListFolderContentsWithFilters(userID uint, folderID, keyword string, page, size int, accessLevel, sortBy, sortOrder string) (*FolderContentResponse, error) {
	if sortBy == "" {
		sortBy = "created_at"
	}
	if sortOrder == "" {
		sortOrder = "desc"
	}

	var folders []models.Folder
	folderQuery := database.DB.Where("user_id = ?", userID)
	if folderID != "" {
		folderQuery = folderQuery.Where("parent_id = ?", folderID)
	} else {
		folderQuery = folderQuery.Where("parent_id = '' OR parent_id IS NULL")
	}
	if keyword != "" {
		folderQuery = folderQuery.Where("name LIKE ?", "%"+keyword+"%")
	}

	var folderOrder string
	switch sortBy {
	case "name":
		folderOrder = "name " + sortOrder
	case "created_at":
		folderOrder = "created_at " + sortOrder
	case "custom":
		folderOrder = "sort_order ASC, name ASC"
	default:
		folderOrder = "created_at " + sortOrder
	}
	if err := folderQuery.Order(folderOrder).Find(&folders).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹列表失败")
	}
	var folderResponses []*FolderResponse
	for i := range folders {
		folderResponses = append(folderResponses, toResponse(&folders[i]))
	}

	var images []models.File
	imageQuery := database.DB.Where("user_id = ?", userID).Where("status <> ?", "pending_deletion")
	if folderID != "" {
		imageQuery = imageQuery.Where("folder_id = ?", folderID)
	} else {
		imageQuery = imageQuery.Where("folder_id = '' OR folder_id IS NULL")
	}
	if accessLevel != "" {
		imageQuery = imageQuery.Where("access_level = ?", accessLevel)
	}
	if keyword != "" {
		imageQuery = imageQuery.Where("original_name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	var imageOrder string
	switch sortBy {
	case "name":
		imageOrder = "original_name " + sortOrder
	case "created_at":
		imageOrder = "created_at " + sortOrder
	case "size":
		imageOrder = "size " + sortOrder
	case "custom":
		imageOrder = "sort_order ASC, created_at DESC"
	default:
		imageOrder = "created_at " + sortOrder
	}
	if err := imageQuery.Order(imageOrder).Find(&images).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}
	var fileResponses []*FileResponse
	for i := range images {
		file := &images[i]
		fullURL, fullThumbURL, shortURL := storage.GetFullURLs(*file)
		aiInfo, _ := getFileAIInfo(file.ID)

		var exifInfo models.FileEXIF
		var exifPtr *models.FileEXIF
		if err := database.DB.Where("file_id = ?", file.ID).First(&exifInfo).Error; err == nil {
			exifPtr = &exifInfo
		}

		fileResponses = append(fileResponses, &FileResponse{
			ID: file.ID, URL: file.URL, ThumbnailURL: file.ThumbURL, FullURL: fullURL, FullThumbURL: fullThumbURL, ShortURL: shortURL,
			OriginalName: file.OriginalName, DisplayName: file.DisplayName, Size: file.Size, Width: file.Width, Height: file.Height,
			Format: file.Format, AccessLevel: file.AccessLevel, FolderID: file.FolderID,
			CreatedAt: file.CreatedAt, UpdatedAt: file.UpdatedAt, IsDuplicate: file.IsDuplicate, MD5Hash: file.MD5Hash,
			IsRecommended: file.IsRecommended, AIInfo: aiInfo, EXIFInfo: exifPtr, StorageDuration: file.StorageDuration, ExpiresAt: (*common.JSONTime)(file.ExpiresAt), IsTimeLimited: file.IsTimeLimitedStorage(),
		})
	}
	return &FolderContentResponse{Folders: folderResponses, Files: fileResponses}, nil
}
