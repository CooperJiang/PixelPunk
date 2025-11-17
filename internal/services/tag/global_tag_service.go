package tag

import (
	"context"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

type GlobalTagResponseItem struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsSystem    bool            `json:"is_system"`
	UsageCount  int             `json:"usage_count"`
	Count       int64           `json:"count"` // 文件引用数量
	CreatedAt   common.JSONTime `json:"created_at"`
	UpdatedAt   common.JSONTime `json:"updated_at"`
}

type TagsPaginationResult struct {
	Items []GlobalTagResponseItem `json:"items"`
	Total int64                   `json:"total"`
}

type GlobalTagService struct {
	db *gorm.DB
}

/* NewGlobalTagService 创建全局标签服务实例 */
func NewGlobalTagService() *GlobalTagService {
	return &GlobalTagService{
		db: database.GetDB(),
	}
}

func (s *GlobalTagService) MergeGlobalTags(operatorID uint, sourceTagIDs []uint, targetTagID uint) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}
	if len(sourceTagIDs) == 0 {
		return fmt.Errorf("源标签列表不能为空")
	}
	for _, id := range sourceTagIDs {
		if id == targetTagID {
			return fmt.Errorf("目标标签不能与源标签相同")
		}
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		var target models.GlobalTag
		if err := tx.First(&target, targetTagID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("目标标签不存在")
			}
			return fmt.Errorf("查询目标标签失败: %v", err)
		}
		var sources []models.GlobalTag
		if err := tx.Where("id IN ?", sourceTagIDs).Find(&sources).Error; err != nil {
			return fmt.Errorf("查询源标签失败: %v", err)
		}
		if len(sources) != len(sourceTagIDs) {
			return fmt.Errorf("部分源标签不存在")
		}
		for _, t := range sources {
			if t.IsSystem {
				return fmt.Errorf("不能合并系统标签: %s", t.Name)
			}
		}

		var existingTargetRels []models.FileGlobalTagRelation
		if err := tx.Where("tag_id = ?", targetTagID).Find(&existingTargetRels).Error; err != nil {
			return fmt.Errorf("查询目标标签关联失败: %v", err)
		}
		existingMap := make(map[string]bool, len(existingTargetRels))
		for _, r := range existingTargetRels {
			existingMap[r.FileID] = true
		}

		var sourceRels []models.FileGlobalTagRelation
		if err := tx.Where("tag_id IN ?", sourceTagIDs).Find(&sourceRels).Error; err != nil {
			return fmt.Errorf("查询源标签关联失败: %v", err)
		}

		newRels := make([]models.FileGlobalTagRelation, 0)
		for _, r := range sourceRels {
			if !existingMap[r.FileID] {
				newRels = append(newRels, models.FileGlobalTagRelation{
					FileID:      r.FileID,
					TagID:       targetTagID,
					UserID:      r.UserID,
					AccessLevel: r.AccessLevel,
					Source:      r.Source,
					Confidence:  r.Confidence,
					CreatedAt:   r.CreatedAt,
				})
				existingMap[r.FileID] = true
			}
		}

		if len(newRels) > 0 {
			if err := tx.Create(&newRels).Error; err != nil {
				return fmt.Errorf("批量创建文件标签关联失败: %v", err)
			}
		}

		if err := tx.Where("tag_id IN ?", sourceTagIDs).Delete(&models.FileGlobalTagRelation{}).Error; err != nil {
			return fmt.Errorf("删除源标签关联失败: %v", err)
		}

		var sourceUserRefs []models.UserTagReference
		if err := tx.Where("tag_id IN ?", sourceTagIDs).Find(&sourceUserRefs).Error; err != nil {
			return fmt.Errorf("查询用户标签引用失败: %v", err)
		}

		if len(sourceUserRefs) > 0 {
			var existingRefs []models.UserTagReference
			if err := tx.Where("tag_id = ?", targetTagID).Find(&existingRefs).Error; err != nil {
				return fmt.Errorf("查询目标用户标签引用失败: %v", err)
			}
			existUser := make(map[uint]bool, len(existingRefs))
			for _, ur := range existingRefs {
				existUser[ur.UserID] = true
			}

			newUserRefs := make([]models.UserTagReference, 0)
			for _, ur := range sourceUserRefs {
				if !existUser[ur.UserID] {
					newUserRefs = append(newUserRefs, models.UserTagReference{
						UserID: ur.UserID,
						TagID:  targetTagID,
						Source: ur.Source,
					})
					existUser[ur.UserID] = true
				}
			}

			if len(newUserRefs) > 0 {
				if err := tx.Create(&newUserRefs).Error; err != nil {
					return fmt.Errorf("批量创建用户标签引用失败: %v", err)
				}
			}

			if err := tx.Where("tag_id IN ?", sourceTagIDs).Delete(&models.UserTagReference{}).Error; err != nil {
				return fmt.Errorf("删除源用户标签引用失败: %v", err)
			}
		}

		var sourceCat []models.TagCategoryRelation
		if err := tx.Where("tag_id IN ?", sourceTagIDs).Find(&sourceCat).Error; err != nil {
			return fmt.Errorf("查询分类关联失败: %v", err)
		}

		if len(sourceCat) > 0 {
			var targetCat []models.TagCategoryRelation
			if err := tx.Where("tag_id = ?", targetTagID).Find(&targetCat).Error; err != nil {
				return fmt.Errorf("查询目标分类关联失败: %v", err)
			}
			byCategory := make(map[uint]*models.TagCategoryRelation, len(targetCat))
			for i := range targetCat {
				byCategory[targetCat[i].CategoryID] = &targetCat[i]
			}

			toUpdate := make([]models.TagCategoryRelation, 0)
			toInsert := make([]models.TagCategoryRelation, 0)

			for _, sc := range sourceCat {
				if tc, ok := byCategory[sc.CategoryID]; ok {
					newWeight := tc.Weight
					if sc.Weight > newWeight {
						newWeight = sc.Weight
					}
					newCore := tc.IsCore || sc.IsCore
					tc.Weight = newWeight
					tc.IsCore = newCore
					toUpdate = append(toUpdate, *tc)
				} else {
					nr := models.TagCategoryRelation{
						TagID:      targetTagID,
						CategoryID: sc.CategoryID,
						Weight:     sc.Weight,
						IsCore:     sc.IsCore,
					}
					toInsert = append(toInsert, nr)
					byCategory[sc.CategoryID] = &nr
				}
			}

			for _, item := range toUpdate {
				if err := tx.Model(&models.TagCategoryRelation{}).
					Where("tag_id = ? AND category_id = ?", item.TagID, item.CategoryID).
					Updates(map[string]interface{}{
						"weight":  item.Weight,
						"is_core": item.IsCore,
					}).Error; err != nil {
					return fmt.Errorf("更新分类关联失败: %v", err)
				}
			}

			if len(toInsert) > 0 {
				if err := tx.Create(&toInsert).Error; err != nil {
					return fmt.Errorf("批量创建分类关联失败: %v", err)
				}
			}

			if err := tx.Where("tag_id IN ?", sourceTagIDs).Delete(&models.TagCategoryRelation{}).Error; err != nil {
				return fmt.Errorf("删除源分类关联失败: %v", err)
			}
		}

		if err := tx.Where("id IN ?", sourceTagIDs).Delete(&models.GlobalTag{}).Error; err != nil {
			return fmt.Errorf("删除源标签失败: %v", err)
		}

		if err := tx.Where("tag_id IN ?", append(sourceTagIDs, targetTagID)).Delete(&models.GlobalTagStatsCache{}).Error; err != nil {
			logger.Warn("清理标签统计缓存失败: %v", err)
		}

		details := fmt.Sprintf("merge %v -> %d", sourceTagIDs, targetTagID)
		log := &models.GlobalTagOperationLog{OperationType: "merge", TagID: &targetTagID, UserID: &operatorID, NewValue: &details}
		_ = tx.Create(log).Error

		return nil
	})
}

/* CreateTag 创建全局标签（显式接口，支持 sort_order） */
func (s *GlobalTagService) CreateTag(name, description string, creatorID uint, isSystem bool, sortOrder *int) (*models.GlobalTag, error) {
	tag, err := s.CreateOrGetGlobalTag(strings.TrimSpace(name), description, creatorID, isSystem)
	if err != nil {
		return nil, err
	}
	if sortOrder != nil && tag.SortOrder != *sortOrder {
		if uerr := s.db.Model(tag).Update("sort_order", *sortOrder).Error; uerr != nil {
			logger.Warn("更新标签排序失败: %v", uerr)
		} else {
			tag.SortOrder = *sortOrder
		}
	}
	return tag, nil
}

/* UpdateGlobalTag 更新全局标签（名称/描述/排序） */
func (s *GlobalTagService) UpdateGlobalTag(tagID uint, name *string, description *string, sortOrder *int) (*models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tag models.GlobalTag
	if err := s.db.First(&tag, tagID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("标签不存在")
		}
		return nil, fmt.Errorf("查询标签失败: %v", err)
	}

	updates := map[string]interface{}{}

	if name != nil {
		newName := strings.TrimSpace(*name)
		if newName == "" {
			return nil, fmt.Errorf("标签名称不能为空")
		}
		if newName != tag.Name {
			var cnt int64
			if err := s.db.Model(&models.GlobalTag{}).Where("name = ? AND id != ?", newName, tag.ID).Count(&cnt).Error; err != nil {
				return nil, fmt.Errorf("检查重名失败: %v", err)
			}
			if cnt > 0 {
				return nil, fmt.Errorf("同名标签已存在")
			}
			updates["name"] = newName
		}
	}

	if description != nil {
		updates["description"] = *description
	}
	if sortOrder != nil {
		updates["sort_order"] = *sortOrder
	}

	if len(updates) == 0 {
		return &tag, nil
	}

	if err := s.db.Model(&tag).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新标签失败: %v", err)
	}

	if err := s.db.First(&tag, tagID).Error; err != nil {
		return nil, fmt.Errorf("获取更新后的标签失败: %v", err)
	}
	return &tag, nil
}

/* CreateOrGetGlobalTag 创建或获取全局标签 */
func (s *GlobalTagService) CreateOrGetGlobalTag(name, description string, creatorID uint, isSystem bool) (*models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var existingTag models.GlobalTag
	err := s.db.Where("name = ?", name).First(&existingTag).Error
	if err == nil {
		return &existingTag, nil
	}

	newTag := &models.GlobalTag{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
		IsSystem:    isSystem,
		UsageCount:  0,
	}

	err = s.db.Create(newTag).Error
	if err != nil {
		return nil, fmt.Errorf("创建全局标签失败: %v", err)
	}

	return newTag, nil
}

/* AddUserTagReference 为用户添加标签引用 */
func (s *GlobalTagService) AddUserTagReference(userID, tagID uint, source string) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	var existingRef models.UserTagReference
	err := s.db.Where("user_id = ? AND tag_id = ?", userID, tagID).First(&existingRef).Error
	if err == nil {
		return nil
	}

	ref := &models.UserTagReference{
		UserID: userID,
		TagID:  tagID,
		Source: source,
	}

	err = s.db.Create(ref).Error
	if err != nil {
		return fmt.Errorf("创建用户标签引用失败: %v", err)
	}

	return nil
}

/* GetUserTags 获取用户的所有标签 */
func (s *GlobalTagService) GetUserTags(userID uint) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag
	err := s.db.Table("global_tag").
		Joins("JOIN user_tag_reference ON global_tag.id = user_tag_reference.tag_id").
		Where("user_tag_reference.user_id = ?", userID).
		Order("global_tag.usage_count DESC, global_tag.name ASC").
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("获取用户标签失败: %v", err)
	}

	return tags, nil
}

/* GetCategoryTags 获取分类关联的标签 */
func (s *GlobalTagService) GetCategoryTags(categoryID uint) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag
	err := s.db.Table("global_tag").
		Joins("JOIN tag_category_relation ON global_tag.id = tag_category_relation.tag_id").
		Where("tag_category_relation.category_id = ?", categoryID).
		Order("tag_category_relation.weight DESC, global_tag.usage_count DESC").
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("获取分类标签失败: %v", err)
	}

	return tags, nil
}

/* AddTagToCategory 将标签关联到分类 */
func (s *GlobalTagService) AddTagToCategory(tagID, categoryID uint, weight int, isCore bool) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	var existingRelation models.TagCategoryRelation
	err := s.db.Where("tag_id = ? AND category_id = ?", tagID, categoryID).First(&existingRelation).Error
	if err == nil {
		err = s.db.Model(&existingRelation).Updates(map[string]interface{}{
			"weight":  weight,
			"is_core": isCore,
		}).Error
		return err
	}

	relation := &models.TagCategoryRelation{
		TagID:      tagID,
		CategoryID: categoryID,
		Weight:     weight,
		IsCore:     isCore,
	}

	err = s.db.Create(relation).Error
	if err != nil {
		return fmt.Errorf("创建标签分类关联失败: %v", err)
	}

	return nil
}

/* GetPopularTags 获取热门标签 */
func (s *GlobalTagService) GetPopularTags(limit int) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag
	err := s.db.Where("usage_count > 0").
		Order("usage_count DESC, created_at DESC").
		Limit(limit).
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("获取热门标签失败: %v", err)
	}

	return tags, nil
}

/* SearchTags 搜索标签 */
func (s *GlobalTagService) SearchTags(keyword string, limit int) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag
	query := s.db.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Order("usage_count DESC, name ASC").Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("搜索标签失败: %v", err)
	}

	return tags, nil
}

/* GetTagBySlug 通过slug获取标签 */
func (s *GlobalTagService) GetTagBySlug(slug string) (*models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tag models.GlobalTag
	err := s.db.Where("slug = ?", slug).First(&tag).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("标签不存在")
		}
		return nil, fmt.Errorf("获取标签失败: %v", err)
	}

	return &tag, nil
}

/* DeleteGlobalTag 删除全局标签（仅在没有任何引用时允许） */
func (s *GlobalTagService) DeleteGlobalTag(tagID uint) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	var userRefCount int64
	err := s.db.Model(&models.UserTagReference{}).Where("tag_id = ?", tagID).Count(&userRefCount).Error
	if err != nil {
		return fmt.Errorf("检查用户引用失败: %v", err)
	}

	if userRefCount > 0 {
		return fmt.Errorf("标签仍被 %d 个用户引用，无法删除", userRefCount)
	}

	var imageRefCount int64
	err = s.db.Model(&models.FileGlobalTagRelation{}).Where("tag_id = ?", tagID).Count(&imageRefCount).Error
	if err != nil {
		return fmt.Errorf("检查文件关联失败: %v", err)
	}

	if imageRefCount > 0 {
		return fmt.Errorf("标签仍被 %d 张文件关联，无法删除", imageRefCount)
	}

	err = s.db.Where("tag_id = ?", tagID).Delete(&models.TagCategoryRelation{}).Error
	if err != nil {
		return fmt.Errorf("删除分类关联失败: %v", err)
	}

	err = s.db.Delete(&models.GlobalTag{}, tagID).Error
	if err != nil {
		return fmt.Errorf("删除标签失败: %v", err)
	}

	return nil
}

/* RemoveUserTagReference 移除用户标签引用 */
func (s *GlobalTagService) RemoveUserTagReference(userID, tagID uint) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	err := s.db.Where("user_id = ? AND tag_id = ?", userID, tagID).Delete(&models.UserTagReference{}).Error
	if err != nil {
		return fmt.Errorf("移除用户标签引用失败: %v", err)
	}

	return nil
}

/* UpdateTagUsageCount 更新标签使用次数（通常由钩子函数自动调用） */
func (s *GlobalTagService) UpdateTagUsageCount(tagID uint, increment bool) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	var expr string
	if increment {
		expr = "usage_count + 1"
	} else {
		expr = "CASE WHEN usage_count > 0 THEN usage_count - 1 ELSE 0 END"
	}

	err := s.db.Model(&models.GlobalTag{}).
		Where("id = ?", tagID).
		UpdateColumn("usage_count", gorm.Expr(expr)).Error

	if err != nil {
		return fmt.Errorf("更新标签使用次数失败: %v", err)
	}

	return nil
}

/* GetTagsByIDs 批量获取标签 */
func (s *GlobalTagService) GetTagsByIDs(tagIDs []uint) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	if len(tagIDs) == 0 {
		return []models.GlobalTag{}, nil
	}

	var tags []models.GlobalTag
	err := s.db.Where("id IN ?", tagIDs).Find(&tags).Error
	if err != nil {
		return nil, fmt.Errorf("批量获取标签失败: %v", err)
	}

	return tags, nil
}

/* CreateTagsFromNames 从标签名称列表创建或获取标签 */
func (s *GlobalTagService) CreateTagsFromNames(names []string, creatorID uint, source string) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag

	for _, name := range names {
		if name == "" {
			continue
		}

		tag, err := s.CreateOrGetGlobalTag(name, "", creatorID, false)
		if err != nil {
			logger.Warn("创建标签失败 [%s]: %v", name, err)
			continue
		}

		err = s.AddUserTagReference(creatorID, tag.ID, source)
		if err != nil {
			logger.Warn("添加用户标签引用失败 [%s]: %v", name, err)
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

/* GetTagFileCount 获取标签的文件引用计数 */
func (s *GlobalTagService) GetTagFileCount(tagID uint) (int64, error) {
	if s.db == nil {
		return 0, fmt.Errorf("数据库连接失败")
	}

	var count int64
	err := s.db.Model(&models.FileGlobalTagRelation{}).Where("tag_id = ?", tagID).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("获取标签文件计数失败: %v", err)
	}

	return count, nil
}

/* ConvertToResponseItem 将GlobalTag转换为响应项，包含文件计数 */
func (s *GlobalTagService) ConvertToResponseItem(tag models.GlobalTag) (GlobalTagResponseItem, error) {
	count, err := s.GetTagFileCount(tag.ID)
	if err != nil {
		logger.Warn("获取标签 %d 的文件计数失败: %v", tag.ID, err)
		count = 0
	}

	return GlobalTagResponseItem{
		ID:          tag.ID,
		Name:        tag.Name,
		Description: tag.Description,
		IsSystem:    tag.IsSystem,
		UsageCount:  tag.UsageCount,
		Count:       count,
		CreatedAt:   tag.CreatedAt,
		UpdatedAt:   tag.UpdatedAt,
	}, nil
}

/* GetTagsWithPagination 分页获取标签列表（支持搜索和排序） */
func (s *GlobalTagService) GetTagsWithPagination(userID uint, keyword, sortBy, sortOrder string, page, size int) (*TagsPaginationResult, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	query := s.db.Table("global_tag").
		Select(`global_tag.id, global_tag.name, global_tag.description,
			global_tag.is_system, global_tag.usage_count,
			global_tag.created_at, global_tag.updated_at,
			CASE
				WHEN user_tags.tag_id IS NOT NULL THEN 1
				ELSE 0
			END as is_user_tag`).
		Joins(`LEFT JOIN user_tag_reference user_tags ON global_tag.id = user_tags.tag_id AND user_tags.user_id = ?`, userID)

	if keyword != "" {
		query = query.Where("global_tag.name LIKE ? OR global_tag.description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	orderClause := "is_user_tag DESC, "
	switch sortBy {
	case "count":
		orderClause += "global_tag.usage_count "
	case "name":
		orderClause += "global_tag.name "
	case "usage_count":
		orderClause += "global_tag.usage_count "
	default:
		orderClause += "global_tag.usage_count "
	}

	if sortOrder == "asc" {
		orderClause += "ASC"
	} else {
		orderClause += "DESC"
	}

	query = query.Order(orderClause)

	var total int64
	countQuery := s.db.Table("global_tag")

	if keyword != "" {
		countQuery = countQuery.Where("global_tag.name LIKE ? OR global_tag.description LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%")
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, fmt.Errorf("获取标签总数失败: %v", err)
	}

	offset := (page - 1) * size
	var results []struct {
		ID          uint            `json:"id"`
		Name        string          `json:"name"`
		Description string          `json:"description"`
		IsSystem    bool            `json:"is_system"`
		UsageCount  int             `json:"usage_count"`
		CreatedAt   common.JSONTime `json:"created_at"`
		UpdatedAt   common.JSONTime `json:"updated_at"`
		IsUserTag   int             `json:"is_user_tag"`
	}

	err = query.Offset(offset).Limit(size).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("查询标签列表失败: %v", err)
	}

	items := make([]GlobalTagResponseItem, len(results))
	for i, result := range results {
		items[i] = GlobalTagResponseItem{
			ID:          result.ID,
			Name:        result.Name,
			Description: result.Description,
			IsSystem:    result.IsSystem,
			UsageCount:  result.UsageCount,
			Count:       int64(result.UsageCount),
			CreatedAt:   result.CreatedAt,
			UpdatedAt:   result.UpdatedAt,
		}
	}

	return &TagsPaginationResult{
		Items: items,
		Total: total,
	}, nil
}

/* CleanupOrphanedSystemTags 清理无关联的系统标签 */
func (s *GlobalTagService) CleanupOrphanedSystemTags(ctx context.Context) (int, error) {
	var deletedCount int

	return deletedCount, s.db.Transaction(func(tx *gorm.DB) error {
		// 使用 GORM ORM 方法替代 Raw SQL，确保 SQLite 兼容性
		var orphanedTags []models.GlobalTag

		// 查询无关联的系统标签：使用 LEFT JOIN 和 WHERE 条件
		err := tx.Model(&models.GlobalTag{}).
			Select("global_tag.*").
			Joins("LEFT JOIN file_global_tag_relation ON global_tag.id = file_global_tag_relation.tag_id").
			Where("file_global_tag_relation.tag_id IS NULL AND global_tag.is_system = ?", true).
			Find(&orphanedTags).Error

		if err != nil {
			return fmt.Errorf("查询无关联系统标签失败: %v", err)
		}

		if len(orphanedTags) == 0 {
			return nil
		}

		var tagIDs []uint
		for _, tag := range orphanedTags {
			tagIDs = append(tagIDs, tag.ID)
		}

		err = tx.Where("id IN ? AND is_system = true", tagIDs).Delete(&models.GlobalTag{}).Error
		if err != nil {
			return fmt.Errorf("删除无关联系统标签失败: %v", err)
		}

		deletedCount = len(orphanedTags)
		return nil
	})
}
