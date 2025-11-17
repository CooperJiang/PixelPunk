package folder

import (
	"log"
	"strings"

	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
)

func SearchFolders(userID uint, keyword string) ([]FolderResponse, error) {
	var folders []models.Folder
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []FolderResponse{}, nil
	}
	query := database.DB.Where("user_id = ?", userID).Where("LOWER(name) LIKE ?", "%"+strings.ToLower(keyword)+"%")
	if err := query.Order("name ASC").Find(&folders).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "搜索文件夹失败")
	}
	var folderResponses []FolderResponse
	for _, folder := range folders {
		fullPath, err := buildFolderPathForExistingFolder(userID, &folder)
		if err != nil {
			log.Printf("构建文件夹路径失败: %v", err)
			fullPath = folder.Name
		}
		var childCount int64
		database.DB.Model(&models.Folder{}).Where("parent_id = ? AND user_id = ?", folder.ID, userID).Count(&childCount)
		var imageCount int64
		database.DB.Model(&models.File{}).Where("folder_id = ? AND user_id = ?", folder.ID, userID).Count(&imageCount)
		folderResponses = append(folderResponses, FolderResponse{
			ID: folder.ID, Name: folder.Name, ParentID: folder.ParentID, Permission: folder.Permission,
			Description: fullPath, FileCount: imageCount, HasChildren: childCount > 0, SortOrder: folder.SortOrder,
			CreatedAt: common.JSONTime(folder.CreatedAt), UpdatedAt: common.JSONTime(folder.UpdatedAt),
		})
	}
	return folderResponses, nil
}
