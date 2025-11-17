package folder

import (
	"fmt"
	"path/filepath"
	"strings"

	"pixelpunk/internal/controllers/folder/dto"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	file "pixelpunk/pkg/storage"

	"gorm.io/gorm"
)

func GetFolderPath(baseDir string, userID uint, folderID string) (string, error) {
	userDir := filepath.Join(baseDir, fmt.Sprintf("user_%d", userID))
	if folderID == "" {
		return userDir, nil
	}
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		return "", err
	}
	folderName := file.SanitizeFileName(folder.Name)
	if folder.ParentID == "" {
		return filepath.Join(userDir, folderName), nil
	}
	parentPath, err := GetFolderPath(baseDir, userID, folder.ParentID)
	if err != nil {
		return "", err
	}
	return filepath.Join(parentPath, folderName), nil
}

func CreateFolderByPath(userID uint, filePath string) (string, error) {
	filePath = strings.Trim(filePath, "/")
	if filePath == "" {
		return "", nil
	}
	parts := strings.Split(filePath, "/")
	var currentParentID string
	for _, folderName := range parts {
		folderName = strings.TrimSpace(folderName)
		if folderName == "" {
			continue
		}
		var existingFolder models.Folder
		result := database.DB.Where("user_id = ? AND parent_id = ? AND name = ?", userID, currentParentID, folderName).First(&existingFolder)
		if result.Error == nil {
			currentParentID = existingFolder.ID
			continue
		}
		if result.Error != gorm.ErrRecordNotFound {
			return "", errors.Wrap(result.Error, errors.CodeDBQueryFailed, "查询文件夹失败")
		}
		newFolder, err := CreateFolder(userID, folderName, currentParentID, "private", "")
		if err != nil {
			return "", errors.Wrap(err, errors.CodeFolderCreateFailed, fmt.Sprintf("创建文件夹路径失败: %s", folderName))
		}
		currentParentID = newFolder.ID
	}
	return currentParentID, nil
}

func GetFolderPathChain(folderID string, userID uint) (*dto.FolderPathChainResponseDTO, error) {
	folderModel, err := models.GetFolderByIDAndUserID(database.DB, folderID, userID)
	if err != nil {
		return nil, errors.New(errors.CodeFolderNotFound, "文件夹不存在或无权访问")
	}
	pathChain := []dto.FolderPathChainDTO{}
	currentFolder := folderModel
	level := 1
	for currentFolder != nil {
		pathChain = append([]dto.FolderPathChainDTO{{ID: currentFolder.ID, Name: currentFolder.Name, ParentID: currentFolder.ParentID, Level: level}}, pathChain...)
		if currentFolder.ParentID == "" {
			break
		}
		parentFolder, err := models.GetFolderByIDAndUserID(database.DB, currentFolder.ParentID, userID)
		if err != nil {
			break
		}
		currentFolder = parentFolder
		level++
	}
	var pathNames []string
	for _, f := range pathChain {
		pathNames = append(pathNames, f.Name)
	}
	fullPath := strings.Join(pathNames, " / ")
	return &dto.FolderPathChainResponseDTO{FolderID: folderID, FullPath: fullPath, PathChain: pathChain, TotalLevels: len(pathChain)}, nil
}

func GetBatchFolderPathChains(folderIDs []string, userID uint) ([]dto.FolderPathChainResponseDTO, error) {
	var results []dto.FolderPathChainResponseDTO
	for _, folderID := range folderIDs {
		pc, err := GetFolderPathChain(folderID, userID)
		if err != nil {
			continue
		}
		results = append(results, *pc)
	}
	return results, nil
}
