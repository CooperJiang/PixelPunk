package file

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"strings"

	"gorm.io/gorm"
)

const (
	MAX_BATCH_MOVE_FILES = 500 // 批量移动文件最大数量
	MAX_FOLDER_DEPTH     = 10  // 文件夹最大层级深度
)

/* MoveFiles 批量移动文件到指定文件夹 */
func MoveFiles(userID uint, fileIDs []string, targetFolderID string) error {
	if len(fileIDs) == 0 {
		return errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空")
	}

	if len(fileIDs) > MAX_BATCH_MOVE_FILES {
		return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("一次最多移动%d个文件", MAX_BATCH_MOVE_FILES))
	}

	if targetFolderID != "" {
		var targetFolder models.Folder
		if err := database.DB.Where("id = ? AND user_id = ?", targetFolderID, userID).First(&targetFolder).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return errors.New(errors.CodeFolderNotFound, "目标文件夹不存在")
			}
			return errors.Wrap(err, errors.CodeDBQueryFailed, "查询目标文件夹失败")
		}
	}

	var count int64
	if err := database.DB.Model(&models.File{}).
		Where("id IN ? AND user_id = ?", fileIDs, userID).
		Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "验证文件所属失败")
	}

	if count != int64(len(fileIDs)) {
		return errors.New(errors.CodeInvalidParameter, "部分文件不存在或无权限")
	}

	result := database.DB.Model(&models.File{}).
		Where("id IN ? AND user_id = ?", fileIDs, userID).
		Update("folder_id", targetFolderID)
	if result.Error != nil {
		return errors.Wrap(result.Error, errors.CodeDBUpdateFailed, "移动文件失败")
	}
	if result.RowsAffected == 0 {
		return errors.New(errors.CodeInvalidParameter, "没有可移动的文件或无权限")
	}

	return nil
}

/* ReorderFiles 重新排序文件 */
func ReorderFiles(userID uint, folderID string, fileIDs []string) error {
	if len(fileIDs) == 0 {
		return errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空")
	}

	if len(fileIDs) > MAX_BATCH_MOVE_FILES {
		return errors.New(errors.CodeInvalidParameter, fmt.Sprintf("一次最多排序%d个文件", MAX_BATCH_MOVE_FILES))
	}

	var count int64
	query := database.DB.Model(&models.File{}).Where("user_id = ? AND id IN ?", userID, fileIDs)

	if folderID == "" {
		query = query.Where("folder_id = '' OR folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}

	if err := query.Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "验证文件所属失败")
	}

	if count != int64(len(fileIDs)) {
		return errors.New(errors.CodeInvalidParameter, "部分文件不存在或不属于指定文件夹")
	}

	if len(fileIDs) == 0 {
		return nil
	}

	sqlBuilder := strings.Builder{}
	sqlBuilder.WriteString("UPDATE file SET sort_order = CASE id ")

	args := make([]interface{}, 0, len(fileIDs)*2+1)
	for i, fileID := range fileIDs {
		sqlBuilder.WriteString("WHEN ? THEN ? ")
		args = append(args, fileID, i+1)
	}
	sqlBuilder.WriteString("END WHERE id IN (?)")
	args = append(args, fileIDs)

	return database.DB.Exec(sqlBuilder.String(), args...).Error
}

/* GetMovableImageCount 获取可移动的文件数量 */
func GetMovableFileCount(userID uint, folderID string) (int64, error) {
	var count int64
	query := database.DB.Model(&models.File{}).Where("user_id = ?", userID)

	if folderID == "" {
		query = query.Where("folder_id = '' OR folder_id IS NULL")
	} else {
		query = query.Where("folder_id = ?", folderID)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取可移动文件数量失败")
	}

	return count, nil
}
