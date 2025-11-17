package file

/* Admin taxonomy queries (tags/colors) — no behavior change. */

import (
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
)

/* AdminGetTagList 管理员获取标签列表 */
func AdminGetTagList(page, size int, userID ...uint) ([]TagWithCount, int64, error) {
	var tags []models.GlobalTag
	var total int64
	query := database.DB.Model(&models.GlobalTag{})

	var specificUserID uint = 0
	var onlyRecommended bool = false
	if len(userID) > 0 && userID[0] > 0 {
		specificUserID = userID[0]
	}
	if len(userID) > 1 && userID[1] == 1 {
		onlyRecommended = true
	}

	imageQuery := database.DB.Model(&models.File{})
	if specificUserID > 0 {
		imageQuery = imageQuery.Where("user_id = ?", specificUserID)
	}
	if onlyRecommended {
		imageQuery = imageQuery.Where("is_recommended = ? AND access_level = ?", true, "public")
	}

	var imageIDs []string
	if specificUserID > 0 || onlyRecommended {
		if err := imageQuery.Pluck("id", &imageIDs).Error; err != nil {
			return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败")
		}
		if len(imageIDs) > 0 {
			var tagIDs []uint
			if err := database.DB.Model(&models.FileGlobalTagRelation{}).Where("file_id IN ?", imageIDs).Pluck("DISTINCT tag_id", &tagIDs).Error; err != nil {
				return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件标签关系失败")
			}
			if len(tagIDs) > 0 {
				query = query.Where("id IN ?", tagIDs)
			} else {
				return []TagWithCount{}, 0, nil
			}
		} else {
			return []TagWithCount{}, 0, nil
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签总数失败")
	}
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}
	query = query.Order("created_at DESC")
	if err := query.Find(&tags).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询标签列表失败")
	}

	var tagsWithCount []TagWithCount
	for _, tag := range tags {
		var count int64
		countQuery := database.DB.Model(&models.FileGlobalTagRelation{}).Where("tag_id = ?", tag.ID)
		if specificUserID > 0 || onlyRecommended {
			countQuery = countQuery.Joins("JOIN file ON file_global_tag_relation.file_id = file.id")
			if specificUserID > 0 {
				countQuery = countQuery.Where("file.user_id = ?", specificUserID)
			}
			if onlyRecommended {
				countQuery = countQuery.Where("file.is_recommended = ? AND file.access_level = ?", true, "public")
			}
		}
		if err := countQuery.Count(&count).Error; err != nil {
			count = 0
		}
		tagsWithCount = append(tagsWithCount, TagWithCount{GlobalTag: tag, Count: count})
	}
	return tagsWithCount, total, nil
}

/* GetPublicTagList 获取公开标签列表（用户端） */
func GetPublicTagList(page, size int, keyword string) ([]TagWithCount, int64, error) {
	var tagsWithCount []TagWithCount
	var total int64

	db := database.DB

	query := db.Model(&models.GlobalTag{}).
		Select("global_tag.id, global_tag.name, global_tag.slug, global_tag.description, global_tag.is_system, global_tag.creator_id, global_tag.sort_order, global_tag.created_at, global_tag.updated_at, COUNT(DISTINCT file_global_tag_relation.file_id) as count").
		Joins("JOIN file_global_tag_relation ON global_tag.id = file_global_tag_relation.tag_id").
		Joins("JOIN file ON file_global_tag_relation.file_id = file.id").
		Where("file.access_level = ?", "public")

	if keyword != "" {
		query = query.Where("global_tag.name LIKE ?", "%"+keyword+"%")
	}

	query = query.Group("global_tag.id").Having("count > 0")

	countQuery := db.Model(&models.GlobalTag{}).
		Select("COUNT(DISTINCT global_tag.id)").
		Joins("JOIN file_global_tag_relation ON global_tag.id = file_global_tag_relation.tag_id").
		Joins("JOIN file ON file_global_tag_relation.file_id = file.id").
		Where("file.access_level = ?", "public")

	if keyword != "" {
		countQuery = countQuery.Where("global_tag.name LIKE ?", "%"+keyword+"%")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "计算总数失败")
	}

	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}

	query = query.Order("global_tag.created_at DESC")

	rows, err := query.Rows()
	if err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询公开标签失败")
	}
	defer rows.Close()

	for rows.Next() {
		var tag TagWithCount
		if scanErr := rows.Scan(&tag.ID, &tag.Name, &tag.Slug, &tag.Description, &tag.IsSystem, &tag.CreatorID, &tag.SortOrder, &tag.CreatedAt, &tag.UpdatedAt, &tag.Count); scanErr != nil {
			return nil, 0, errors.Wrap(errors.New(errors.CodeDBQueryFailed, scanErr.Error()), errors.CodeDBQueryFailed, "扫描公开标签失败")
		}
		tagsWithCount = append(tagsWithCount, tag)
	}

	return tagsWithCount, total, nil
}

/* AdminGetColorList 管理员获取主色调列表 */
func AdminGetColorList(page, size int, userID ...uint) ([]string, int64, error) {
	var colors []string
	var total int64
	var specificUserID uint = 0
	var onlyRecommended bool = false
	if len(userID) > 0 && userID[0] > 0 {
		specificUserID = userID[0]
	}
	if len(userID) > 1 && userID[1] == 1 {
		onlyRecommended = true
	}
	query := database.DB.Model(&models.FileAIInfo{}).Select("DISTINCT dominant_color").Where("dominant_color IS NOT NULL AND dominant_color != ''")
	if specificUserID > 0 || onlyRecommended {
		query = query.Joins("JOIN file ON file_ai_info.file_id = file.id")
		if specificUserID > 0 {
			query = query.Where("file.user_id = ?", specificUserID)
		}
		if onlyRecommended {
			query = query.Where("file.is_recommended = ? AND file.access_level = ?", true, "public")
		}
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取颜色总数失败")
	}
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		query = query.Offset(offset).Limit(size)
	}
	if err := query.Pluck("dominant_color", &colors).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询颜色列表失败")
	}
	return colors, total, nil
}
