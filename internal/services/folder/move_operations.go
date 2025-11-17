package folder

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"strings"

	"gorm.io/gorm"
)

const (
	MAX_FOLDER_DEPTH = 10 // 文件夹最大层级深度
)

func MoveFolders(userID uint, folderIDs []string, newParentID string) error {
	if len(folderIDs) == 0 {
		return errors.New(errors.CodeInvalidParameter, "文件夹ID列表不能为空")
	}

	if len(folderIDs) > 100 {
		return errors.New(errors.CodeInvalidParameter, "单次最多移动100个文件夹")
	}

	var count int64
	if err := database.DB.Model(&models.Folder{}).
		Where("id IN ? AND user_id = ?", folderIDs, userID).
		Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "验证文件夹所属失败")
	}

	if count != int64(len(folderIDs)) {
		return errors.New(errors.CodeInvalidParameter, "部分文件夹不存在或无权限")
	}

	if newParentID != "" {
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", newParentID, userID).First(&parentFolder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New(errors.CodeFolderNotFound, "目标父文件夹不存在")
			}
			return errors.Wrap(err, errors.CodeDBQueryFailed, "查询目标父文件夹失败")
		}
	}

	for _, folderID := range folderIDs {
		if folderID == newParentID {
			return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己")
		}

		if newParentID != "" && isSubFolder(userID, folderID, newParentID) {
			return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己的子文件夹")
		}
	}

	if err := database.DB.Model(&models.Folder{}).
		Where("id IN ? AND user_id = ?", folderIDs, userID).
		Update("parent_id", newParentID).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "批量移动文件夹失败")
	}

	return nil
}

/* MoveFolder 单个文件夹移动（保留向后兼容） */
func MoveFolder(userID uint, folderID, newParentID string) error {
	if folderID == "" {
		return errors.New(errors.CodeInvalidParameter, "文件夹ID不能为空")
	}

	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeFolderNotFound, "源文件夹不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询源文件夹失败")
	}

	if folderID == newParentID {
		return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己")
	}

	if newParentID != "" {
		var parentFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", newParentID, userID).First(&parentFolder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New(errors.CodeFolderNotFound, "目标父文件夹不存在")
			}
			return errors.Wrap(err, errors.CodeDBQueryFailed, "查询目标父文件夹失败")
		}

		if isSubFolder(userID, folderID, newParentID) {
			return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己的子文件夹")
		}

	}

	if err := database.DB.Model(&folder).Update("parent_id", newParentID).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "移动文件夹失败")
	}

	return nil
}

func ReorderFolders(userID uint, parentID string, folderIDs []string) error {
	if len(folderIDs) == 0 {
		return errors.New(errors.CodeInvalidParameter, "文件夹ID列表不能为空")
	}

	var count int64
	query := database.DB.Model(&models.Folder{}).Where("user_id = ? AND id IN ?", userID, folderIDs)

	if parentID == "" {
		query = query.Where("parent_id = '' OR parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", parentID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "验证文件夹所属失败")
	}

	if count != int64(len(folderIDs)) {
		return errors.New(errors.CodeInvalidParameter, "部分文件夹不存在或不属于指定父目录")
	}

	if len(folderIDs) == 0 {
		return nil
	}

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("UPDATE folder SET sort_order = CASE id ")

	args := make([]interface{}, 0, len(folderIDs)*2+1)
	for i, folderID := range folderIDs {
		sqlBuilder.WriteString("WHEN ? THEN ? ")
		args = append(args, folderID, i+1)
	}
	sqlBuilder.WriteString("END WHERE id IN (?)")
	args = append(args, folderIDs)

	return database.DB.Exec(sqlBuilder.String(), args...).Error
}

func isSubFolder(userID uint, parentID, childID string) bool {
	var childFolders []models.Folder
	if err := database.DB.Where("parent_id = ? AND user_id = ?", parentID, userID).Find(&childFolders).Error; err != nil {
		return false
	}

	for _, folder := range childFolders {
		if folder.ID == childID {
			return true
		}
		if isSubFolder(userID, folder.ID, childID) {
			return true
		}
	}

	return false
}

func calculateFolderDepth(userID uint, folderID string) int {
	if folderID == "" || folderID == "0" {
		return 0
	}

	var folder models.Folder
	if err := database.DB.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		return 0
	}

	return 1 + calculateFolderDepth(userID, folder.ParentID)
}

func GetFolderDepth(userID uint, folderID string) int {
	return calculateFolderDepth(userID, folderID)
}

func ValidateFolderMove(userID uint, folderID, newParentID string) error {
	if folderID == newParentID {
		return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己")
	}

	var count int64
	if err := database.DB.Model(&models.Folder{}).
		Where("id = ? AND user_id = ?", folderID, userID).
		Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "查询源文件夹失败")
	}
	if count == 0 {
		return errors.New(errors.CodeFolderNotFound, "源文件夹不存在")
	}

	if newParentID != "" {
		if err := database.DB.Model(&models.Folder{}).
			Where("id = ? AND user_id = ?", newParentID, userID).
			Count(&count).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "查询目标文件夹失败")
		}
		if count == 0 {
			return errors.New(errors.CodeFolderNotFound, "目标父文件夹不存在")
		}

		if isSubFolder(userID, folderID, newParentID) {
			return errors.New(errors.CodeInvalidParameter, "不能将文件夹移动到自己的子文件夹")
		}

	}

	return nil
}
