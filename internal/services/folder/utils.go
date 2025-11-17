package folder

import (
	"encoding/json"
	"path/filepath"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"gorm.io/gorm"
)

const (
	DefaultUploadDir    = "uploads/files"
	DefaultThumbnailDir = "uploads/thumbnails"
)

func toResponse(folder *models.Folder) *FolderResponse {
	var fileCount int64
	database.DB.Model(&models.File{}).Where("folder_id = ?", folder.ID).Count(&fileCount)
	var childCount int64
	database.DB.Model(&models.Folder{}).Where("parent_id = ?", folder.ID).Count(&childCount)
	level := calculateFolderLevel(folder.UserID, folder.ID)
	return &FolderResponse{
		ID:          folder.ID,
		Name:        folder.Name,
		ParentID:    folder.ParentID,
		Permission:  folder.Permission,
		Description: folder.Description,
		FileCount:   fileCount,
		HasChildren: childCount > 0,
		SortOrder:   folder.SortOrder,
		Level:       level,
		CreatedAt:   folder.CreatedAt,
		UpdatedAt:   folder.UpdatedAt,
	}
}

func calculateFolderLevel(userID uint, folderID string) int {
	var folder models.Folder
	err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error
	if err != nil {
		return 1
	}
	if folder.ParentID == "" {
		return 1
	}
	return calculateFolderLevel(userID, folder.ParentID) + 1
}

func getFileAIInfo(fileID string) (*FolderFileAIInfo, error) {
	var aiInfo models.FileAIInfo
	if err := database.DB.Where("file_id = ?", fileID).First(&aiInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件AI信息失败")
	}
	var tagNames []string
	if len(aiInfo.Tags) > 0 {
		_ = json.Unmarshal(aiInfo.Tags, &tagNames)
	}
	if len(tagNames) == 0 {
		var tags []models.GlobalTag
		database.DB.Model(&models.GlobalTag{}).
			Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").
			Where("file_global_tag_relation.file_id = ?", fileID).
			Find(&tags)
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}
	return &FolderFileAIInfo{
		Description:    aiInfo.Description,
		Tags:           tagNames,
		DominantColor:  aiInfo.DominantColor,
		Resolution:     aiInfo.Resolution,
		IsNSFW:         aiInfo.IsNSFW,
		NSFWScore:      aiInfo.NSFWScore,
		NSFWEvaluation: aiInfo.NSFWEvaluation,
	}, nil
}

func getFolderPathFromDB(userID uint, folderID string) (string, error) {
	var folder models.Folder
	err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New(errors.CodeFolderNotFound, "父文件夹不存在")
		}
		return "", errors.Wrap(err, errors.CodeDBQueryFailed, "查询父文件夹失败")
	}
	if folder.ParentID == "" || folder.ParentID == "0" {
		return folder.Name, nil
	}
	parentPath, err := getFolderPathFromDB(userID, folder.ParentID)
	if err != nil {
		return "", err
	}
	return filepath.Join(parentPath, folder.Name), nil
}

func buildFolderPathForExistingFolder(userID uint, folder *models.Folder) (string, error) {
	if folder.ParentID == "" || folder.ParentID == "0" {
		return folder.Name, nil
	}
	parentPath, err := getFolderPathFromDB(userID, folder.ParentID)
	if err != nil {
		return "", err
	}
	return filepath.Join(parentPath, folder.Name), nil
}
