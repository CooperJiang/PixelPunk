package user_tag

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

/* GetUserTagList 获取用户标签列表（分页） */
func GetUserTagList(userID uint, keyword, sortBy, sortOrder string, page, pageSize int) ([]TagWithCount, int64, error) {
	db := database.GetDB()

	// 基础查询：获取用户引用的标签
	query := db.Model(&models.UserTagReference{}).
		Select("global_tag.*, user_tag_reference.source, user_tag_reference.created_at as reference_created_at, "+
			"COUNT(DISTINCT file_global_tag_relation.id) as file_count").
		Joins("JOIN global_tag ON global_tag.id = user_tag_reference.tag_id").
		Joins("LEFT JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id AND file_global_tag_relation.user_id = ?", userID).
		Where("user_tag_reference.user_id = ?", userID).
		Group("global_tag.id, user_tag_reference.source, user_tag_reference.created_at")

	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		query = query.Where("global_tag.name LIKE ?", "%"+keyword+"%")
	}

	// 获取总数（需要子查询）
	var total int64
	countQuery := db.Model(&models.UserTagReference{}).
		Joins("JOIN global_tag ON global_tag.id = user_tag_reference.tag_id").
		Where("user_tag_reference.user_id = ?", userID)
	if keyword != "" {
		countQuery = countQuery.Where("global_tag.name LIKE ?", "%"+keyword+"%")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		logger.Error("获取标签总数失败: %v", err)
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签列表失败")
	}

	validSortBy := map[string]bool{
		"name":       true,
		"file_count": true,
		"created_at": true,
	}
	if !validSortBy[sortBy] {
		sortBy = "file_count"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)
	query = query.Order(orderClause)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	var tags []TagWithCount
	if err := query.Scan(&tags).Error; err != nil {
		logger.Error("获取标签列表失败: %v", err)
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签列表失败")
	}

	return tags, total, nil
}

/* GetAllUserTags 获取用户所有标签（不分页） */
func GetAllUserTags(userID uint) ([]models.GlobalTag, error) {
	db := database.GetDB()

	var tags []models.GlobalTag
	err := db.Model(&models.GlobalTag{}).
		Joins("JOIN user_tag_reference ON user_tag_reference.tag_id = global_tag.id").
		Where("user_tag_reference.user_id = ?", userID).
		Order("global_tag.name ASC").
		Find(&tags).Error

	if err != nil {
		logger.Error("获取用户所有标签失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签列表失败")
	}

	return tags, nil
}

/* CreateUserTag 创建用户标签 */
func CreateUserTag(userID uint, name string) (*models.GlobalTag, error) {
	db := database.GetDB()

	var existingTag models.GlobalTag
	err := db.Where("name = ?", name).First(&existingTag).Error

	if err == nil {
		// 标签已存在，检查用户是否已引用
		var ref models.UserTagReference
		err = db.Where("user_id = ? AND tag_id = ?", userID, existingTag.ID).First(&ref).Error
		if err == nil {
			return nil, errors.New(errors.CodeDBDuplicate, "标签已存在")
		}

		ref = models.UserTagReference{
			UserID: userID,
			TagID:  existingTag.ID,
			Source: "manual",
		}
		if err := db.Create(&ref).Error; err != nil {
			logger.Error("创建标签引用失败: %v", err)
			return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建标签失败")
		}

		return &existingTag, nil
	}

	// 标签不存在，创建新标签
	if err != gorm.ErrRecordNotFound {
		logger.Error("查询标签失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "创建标签失败")
	}

	// 使用事务创建标签和引用
	var newTag models.GlobalTag
	err = db.Transaction(func(tx *gorm.DB) error {
		newTag = models.GlobalTag{
			Name:      name,
			CreatorID: userID,
			IsSystem:  false,
		}
		if err := tx.Create(&newTag).Error; err != nil {
			logger.Error("创建全局标签失败: %v", err)
			return err
		}

		ref := models.UserTagReference{
			UserID: userID,
			TagID:  newTag.ID,
			Source: "manual",
		}
		if err := tx.Create(&ref).Error; err != nil {
			logger.Error("创建标签引用失败: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建标签失败")
	}

	return &newTag, nil
}

/* UpdateUserTag 更新标签（只有创建者可以更新） */
func UpdateUserTag(userID, tagID uint, name string) (*models.GlobalTag, error) {
	db := database.GetDB()

	var tag models.GlobalTag
	if err := db.Where("id = ?", tagID).First(&tag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "标签不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签失败")
	}

	if tag.CreatorID != userID {
		return nil, errors.New(errors.CodeForbidden, "只有创建者可以修改标签")
	}

	if name != tag.Name {
		var existingTag models.GlobalTag
		if err := db.Where("name = ? AND id != ?", name, tagID).First(&existingTag).Error; err == nil {
			return nil, errors.New(errors.CodeDBDuplicate, "标签名称已存在")
		}
	}

	if err := db.Model(&tag).Update("name", name).Error; err != nil {
		logger.Error("更新标签失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新标签失败")
	}

	if err := db.Where("id = ?", tagID).First(&tag).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取更新后的标签失败")
	}

	return &tag, nil
}

/* DeleteUserTag 删除用户标签引用 */
func DeleteUserTag(userID, tagID uint) error {
	db := database.GetDB()

	var ref models.UserTagReference
	if err := db.Where("user_id = ? AND tag_id = ?", userID, tagID).First(&ref).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "标签引用不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签引用失败")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ? AND tag_id = ?", userID, tagID).
			Delete(&models.FileGlobalTagRelation{}).Error; err != nil {
			logger.Error("删除文件标签关联失败: %v", err)
			return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除标签失败")
		}

		if err := tx.Delete(&ref).Error; err != nil {
			logger.Error("删除标签引用失败: %v", err)
			return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除标签失败")
		}

		return nil
	})
}

/* BatchDeleteUserTags 批量删除用户标签引用 */
func BatchDeleteUserTags(userID uint, tagIDs []uint) (int, error) {
	db := database.GetDB()

	deletedCount := 0

	return deletedCount, db.Transaction(func(tx *gorm.DB) error {
		var refs []models.UserTagReference
		if err := tx.Where("user_id = ? AND tag_id IN ?", userID, tagIDs).Find(&refs).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "获取标签引用失败")
		}

		if len(refs) == 0 {
			return errors.New(errors.CodeNotFound, "未找到可删除的标签")
		}

		validIDs := make([]uint, len(refs))
		for i, ref := range refs {
			validIDs[i] = ref.TagID
		}

		if err := tx.Where("user_id = ? AND tag_id IN ?", userID, validIDs).
			Delete(&models.FileGlobalTagRelation{}).Error; err != nil {
			logger.Error("批量删除文件标签关联失败: %v", err)
			return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除标签失败")
		}

		result := tx.Where("user_id = ? AND tag_id IN ?", userID, validIDs).
			Delete(&models.UserTagReference{})
		if result.Error != nil {
			logger.Error("批量删除标签引用失败: %v", result.Error)
			return errors.Wrap(result.Error, errors.CodeDBDeleteFailed, "批量删除标签失败")
		}

		deletedCount = int(result.RowsAffected)
		return nil
	})
}

/* MergeUserTags 合并标签（将多个标签合并到目标标签） */
func MergeUserTags(userID uint, sourceTagIDs []uint, targetTagID uint) error {
	db := database.GetDB()

	var targetTag models.GlobalTag
	if err := db.Where("id = ?", targetTagID).First(&targetTag).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "目标标签不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取目标标签失败")
	}

	var targetRef models.UserTagReference
	if err := db.Where("user_id = ? AND tag_id = ?", userID, targetTagID).First(&targetRef).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			targetRef = models.UserTagReference{
				UserID: userID,
				TagID:  targetTagID,
				Source: "manual",
			}
			if err := db.Create(&targetRef).Error; err != nil {
				return errors.Wrap(err, errors.CodeDBCreateFailed, "创建目标标签引用失败")
			}
		} else {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "获取目标标签引用失败")
		}
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 将源标签的文件关联转移到目标标签
		for _, sourceID := range sourceTagIDs {
			if sourceID == targetTagID {
				continue
			}

			var relations []models.FileGlobalTagRelation
			if err := tx.Where("user_id = ? AND tag_id = ?", userID, sourceID).
				Find(&relations).Error; err != nil {
				logger.Error("获取源标签文件关联失败: %v", err)
				return errors.Wrap(err, errors.CodeDBQueryFailed, "合并标签失败")
			}

			for _, rel := range relations {
				var existingRel models.FileGlobalTagRelation
				err := tx.Where("file_id = ? AND tag_id = ? AND user_id = ?", rel.FileID, targetTagID, userID).
					First(&existingRel).Error

				if err == gorm.ErrRecordNotFound {
					if err := tx.Model(&rel).Update("tag_id", targetTagID).Error; err != nil {
						logger.Error("转移文件标签关联失败: %v", err)
						return errors.Wrap(err, errors.CodeDBUpdateFailed, "合并标签失败")
					}
				} else if err == nil {
					if err := tx.Delete(&rel).Error; err != nil {
						logger.Error("删除重复文件标签关联失败: %v", err)
						return errors.Wrap(err, errors.CodeDBDeleteFailed, "合并标签失败")
					}
				} else {
					return errors.Wrap(err, errors.CodeDBQueryFailed, "检查文件标签关联失败")
				}
			}

			if err := tx.Where("user_id = ? AND tag_id = ?", userID, sourceID).
				Delete(&models.UserTagReference{}).Error; err != nil {
				logger.Error("删除源标签引用失败: %v", err)
				return errors.Wrap(err, errors.CodeDBDeleteFailed, "合并标签失败")
			}
		}

		return nil
	})
}

/* GetUserTagStats 获取用户标签统计 */
func GetUserTagStats(userID uint) (*TagStats, error) {
	db := database.GetDB()

	var stats TagStats

	var totalTags int64
	if err := db.Model(&models.UserTagReference{}).
		Where("user_id = ?", userID).
		Count(&totalTags).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取统计失败")
	}
	stats.TotalTags = totalTags

	var userCreatedTags int64
	if err := db.Model(&models.UserTagReference{}).
		Where("user_id = ? AND source = ?", userID, "manual").
		Count(&userCreatedTags).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取统计失败")
	}
	stats.UserCreatedTags = userCreatedTags

	var aiCreatedTags int64
	if err := db.Model(&models.UserTagReference{}).
		Where("user_id = ? AND source = ?", userID, "ai").
		Count(&aiCreatedTags).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取统计失败")
	}
	stats.AICreatedTags = aiCreatedTags

	var totalUsage int64
	if err := db.Model(&models.FileGlobalTagRelation{}).
		Where("user_id = ?", userID).
		Count(&totalUsage).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取统计失败")
	}
	stats.TotalUsage = totalUsage

	return &stats, nil
}

/* TagWithCount 带文件计数的标签结构 */
type TagWithCount struct {
	ID                 uint   `json:"id"`
	Name               string `json:"name"`
	Slug               string `json:"slug"`
	Description        string `json:"description"`
	IsSystem           bool   `json:"is_system"`
	CreatorID          uint   `json:"creator_id"`
	UsageCount         int    `json:"usage_count"`
	SortOrder          int    `json:"sort_order"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
	Source             string `json:"source"`
	ReferenceCreatedAt string `json:"reference_created_at"`
	FileCount          int64  `json:"file_count"`
}

/* TagStats 标签统计结构 */
type TagStats struct {
	TotalTags       int64 `json:"total_tags"`
	UserCreatedTags int64 `json:"user_created_tags"`
	AICreatedTags   int64 `json:"ai_created_tags"`
	TotalUsage      int64 `json:"total_usage"`
}
