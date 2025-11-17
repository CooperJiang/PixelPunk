package file

/* Query helpers (no behavior change). */

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"strings"

	"gorm.io/gorm"
)

/* GetImageList 获取文件列表 */
func GetFileList(userID uint, folderID string, page, size int, sort string, accessLevel string, keyword string, tags []string, dominantColor []string, resolution string, minWidth, maxWidth, minHeight, maxHeight int) ([]FileDetailResponse, int64, error) {
	var total int64
	var images []models.File
	var responses []FileDetailResponse

	query := database.DB.Where("user_id = ?", userID).Where("status <> ?", StatusPendingDeletion).Joins("LEFT JOIN file_ai_info ON file_ai_info.file_id = file.id")
	if folderID != "" {
		query = query.Where("folder_id = ?", folderID)
	}
	if accessLevel != "" {
		query = query.Where("access_level = ?", accessLevel)
	}
	if keyword != "" {
		query = query.Where("original_name LIKE ? OR file_path LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if len(tags) > 0 {
		var tagIDs []uint
		if err := database.DB.Model(&models.GlobalTag{}).Where("name IN ?", tags).Pluck("id", &tagIDs).Error; err != nil {
			return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询标签失败")
		}
		if len(tagIDs) == 0 {
			return []FileDetailResponse{}, 0, nil
		}
		sub := database.DB.Model(&models.FileGlobalTagRelation{}).Select("DISTINCT file_id").Where("tag_id IN ?", tagIDs)
		if len(tags) > 1 {
			sub = sub.Group("file_id").Having("COUNT(DISTINCT tag_id) = ?", len(tagIDs))
		}
		query = query.Where("id IN (?)", sub)
	}

	if len(dominantColor) > 0 {
		var all []string
		for _, c := range dominantColor {
			if strings.HasPrefix(c, "#") {
				all = append(all, c, strings.TrimPrefix(c, "#"))
			} else {
				all = append(all, c, "#"+c)
			}
		}
		query = query.Where("file_ai_info.dominant_color IN ?", all)
	}
	if resolution != "" {
		query = query.Where("file_ai_info.resolution = ?", resolution)
	}
	if minWidth > 0 {
		query = query.Where("width >= ?", minWidth)
	}
	if maxWidth > 0 {
		query = query.Where("width <= ?", maxWidth)
	}
	if minHeight > 0 {
		query = query.Where("height >= ?", minHeight)
	}
	if maxHeight > 0 {
		query = query.Where("height <= ?", maxHeight)
	}

	if err := query.Model(&models.File{}).Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取文件总数失败")
	}
	switch sort {
	case "newest":
		query = query.Order("created_at DESC")
	case "oldest":
		query = query.Order("created_at ASC")
	case "name":
		query = query.Order("original_name ASC")
	case "size":
		query = query.Order("size DESC")
	case "width":
		query = query.Order("width DESC")
	case "height":
		query = query.Order("height DESC")
	case "quality":
		query = query.Order("file_ai_info.nsfw_score ASC")
	case "nsfw_score":
		query = query.Order("file_ai_info.nsfw_score DESC")
	default:
		query = query.Order("created_at DESC")
	}
	offset := (page - 1) * size
	if err := query.Offset(offset).Limit(size).Find(&images).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}
	for _, file := range images {
		aiInfo, _ := GetFileAIInfo(file.ID)
		responses = append(responses, BuildFileDetailResponse(file, 0, aiInfo))
	}
	return responses, total, nil
}

/* GetFileDetail 获取单个文件详情 */
func GetFileDetail(userID uint, fileID string) (*FileDetailResponse, error) {
	var file models.File
	if err := database.DB.Where("id = ? AND user_id = ?", fileID, userID).Where("status <> ?", StatusPendingDeletion).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeFileNotFound, "文件不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
	}
	var stats models.FileStats
	if err := database.DB.Where("file_id = ?", fileID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			stats = models.FileStats{FileID: fileID}
		}
	}
	aiInfo, _ := GetFileAIInfo(file.ID)

	var exifInfo models.FileEXIF
	resp := BuildFileDetailResponse(file, stats.Views, aiInfo)
	if err := database.DB.Where("file_id = ?", fileID).First(&exifInfo).Error; err == nil {
		resp.EXIFInfo = &exifInfo
	}

	return &resp, nil
}
