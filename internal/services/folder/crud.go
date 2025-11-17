package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/stats"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	file "pixelpunk/pkg/storage"

	"gorm.io/gorm"
)

func CreateFolder(userID uint, name, parentID, permission, description string) (*FolderResponse, error) {
	if !file.IsValidFolderName(name) {
		return nil, errors.New(errors.CodeInvalidParameter, "文件夹名称无效：不能为空或包含 / \\ : * ? \" < > | 等特殊字符")
	}
	if parentID != "" {
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", parentID, userID).First(&parentFolder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New(errors.CodeFolderNotFound, "父文件夹不存在")
			}
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询父文件夹失败")
		}
	}
	var count int64
	if err := database.DB.Model(&models.Folder{}).
		Where("user_id = ? AND parent_id = ? AND name = ?", userID, parentID, name).
		Count(&count).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	if count > 0 {
		return nil, errors.New(errors.CodeFolderNameDuplicate, "同级目录下已存在同名文件夹")
	}

	folder := models.Folder{ID: file.GenerateFolderID(), UserID: userID, ParentID: parentID, Name: name, Permission: permission, Description: description}
	if err := database.DB.Create(&folder).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeFolderCreateFailed, "创建文件夹失败")
	}
	stats.GetStatsAdapter().RecordFolderCreated()

	return toResponse(&folder), nil
}

func UpdateFolder(userID uint, folderID, name, parentID, permission, description string) (*FolderResponse, error) {
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFolderNotFound, "文件夹不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	if name != "" && name != folder.Name {
		if !file.IsValidFolderName(name) {
			return nil, errors.New(errors.CodeInvalidParameter, "文件夹名称无效：不能为空或包含 / \\ : * ? \" < > | 等特殊字符")
		}
	}
	if parentID != "" && parentID != folder.ParentID {
		if folderID == parentID {
			return nil, errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己下面")
		}
		var subFolders []models.Folder
		if err := database.DB.Where("parent_id = ?", folderID).Find(&subFolders).Error; err != nil {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询子文件夹失败")
		}
		var isSubFolder func(string) bool
		isSubFolder = func(id string) bool {
			for _, sub := range subFolders {
				if sub.ID == id {
					return true
				}
				var childFolders []models.Folder
				if err := database.DB.Where("parent_id = ?", sub.ID).Find(&childFolders).Error; err == nil {
					subFolders = append(subFolders, childFolders...)
					if isSubFolder(id) {
						return true
					}
				}
			}
			return false
		}
		if isSubFolder(parentID) {
			return nil, errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己的子文件夹下")
		}
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", parentID, userID).First(&parentFolder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return nil, errors.New(errors.CodeFolderNotFound, "父文件夹不存在")
			}
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询父文件夹失败")
		}
	}
	if name != "" && name != folder.Name {
		var count int64
		if err := database.DB.Model(&models.Folder{}).
			Where("user_id = ? AND parent_id = ? AND name = ? AND id != ?", userID, folder.ParentID, name, folderID).
			Count(&count).Error; err != nil {
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
		}
		if count > 0 {
			return nil, errors.New(errors.CodeFolderNameDuplicate, "同级目录下已存在同名文件夹")
		}
	}
	if name != "" {
		folder.Name = name
	}
	if parentID != "" {
		folder.ParentID = parentID
	}
	if permission != "" {
		folder.Permission = permission
	}
	if description != "" {
		folder.Description = description
	}
	if err := database.DB.Save(&folder).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeFolderUpdateFailed, "更新文件夹失败")
	}
	return toResponse(&folder), nil
}

func DeleteFolder(userID uint, folderID string) error {
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeFolderNotFound, "文件夹不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	var imageCount int64
	if err := database.DB.Model(&models.File{}).Where("folder_id = ?", folderID).Count(&imageCount).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	if imageCount > 0 {
		return errors.New(errors.CodeFolderNotEmpty, "请先删除文件夹内的文件")
	}
	var subFolderCount int64
	if err := database.DB.Model(&models.Folder{}).Where("parent_id = ?", folderID).Count(&subFolderCount).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询子文件夹失败")
	}
	if subFolderCount > 0 {
		return errors.New(errors.CodeFolderNotEmpty, "请先删除子文件夹")
	}
	if err := database.DB.Delete(&folder).Error; err != nil {
		return errors.Wrap(err, errors.CodeFolderDeleteFailed, "删除文件夹失败")
	}
	return nil
}

func GetFolderDetail(userID uint, folderID string) (*FolderResponse, error) {
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFolderNotFound, "文件夹不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	return toResponse(&folder), nil
}

/* ToggleFolderAccessLevel 切换文件夹权限状态（private/public） */
func ToggleFolderAccessLevel(userID uint, folderID string) (*FolderResponse, error) {
	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFolderNotFound, "文件夹不存在或无权限访问")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件夹失败")
	}
	newPermission := "public"
	if folder.Permission == "public" {
		newPermission = "private"
	}
	if err := database.DB.Model(&folder).Update("permission", newPermission).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新文件夹权限失败")
	}
	folder.Permission = newPermission
	return toResponse(&folder), nil
}
