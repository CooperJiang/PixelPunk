package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	file "pixelpunk/pkg/storage"

	"gorm.io/gorm"
)

func GetOrCreateGuestFolder() (string, error) {
	const guestFolderName = "guest"
	const guestFolderDescription = "系统文件夹，用于存放游客上传的文件"

	var guestFolder models.Folder
	err := database.DB.Where("name = ? AND parent_id = ''", guestFolderName).First(&guestFolder).Error
	if err == nil {
		return guestFolder.ID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return "", errors.Wrap(err, errors.CodeDBQueryFailed, "查询游客文件夹失败")
	}

	newFolder := models.Folder{ID: file.GenerateFolderID(), Name: guestFolderName, ParentID: "", Permission: "public", Description: guestFolderDescription, UserID: 0, SortOrder: 999}
	if err := database.DB.Create(&newFolder).Error; err != nil {
		return "", errors.Wrap(err, errors.CodeDBCreateFailed, "创建游客文件夹失败")
	}
	return newFolder.ID, nil
}
