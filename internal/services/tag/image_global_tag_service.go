package tag

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/ai"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* FileGlobalTagService 文件全局标签关联服务 */
type FileGlobalTagService struct {
	db *gorm.DB
}

/* NewFileGlobalTagService 创建文件全局标签关联服务实例 */
func NewFileGlobalTagService() *FileGlobalTagService {
	return &FileGlobalTagService{
		db: database.GetDB(),
	}
}

/* AddTagsToFile 为文件添加标签 */
func (s *FileGlobalTagService) AddTagsToFile(fileID string, tagIDs []uint, source string, confidence float64) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	var file models.File
	err := s.db.Where("id = ?", fileID).First(&file).Error
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	if len(tagIDs) == 0 {
		return nil
	}

	var existingRelations []models.FileGlobalTagRelation
	s.db.Where("file_id = ? AND tag_id IN ?", fileID, tagIDs).Find(&existingRelations)

	existingMap := make(map[uint]bool)
	for _, rel := range existingRelations {
		existingMap[rel.TagID] = true
	}

	newRelations := make([]models.FileGlobalTagRelation, 0)
	actuallyAddedTagIDs := make([]uint, 0) // 实际新增的标签ID
	for _, tagID := range tagIDs {
		if !existingMap[tagID] {
			newRelations = append(newRelations, models.FileGlobalTagRelation{
				FileID:      fileID,
				TagID:       tagID,
				UserID:      file.UserID,
				AccessLevel: file.AccessLevel,
				Source:      source,
				Confidence:  confidence,
			})
			actuallyAddedTagIDs = append(actuallyAddedTagIDs, tagID)
		}
	}

	if len(newRelations) == 0 {
		return nil
	}

	if err := s.db.Create(&newRelations).Error; err != nil {
		return fmt.Errorf("批量创建文件标签关联失败: %v", err)
	}

	// 手动批量更新受影响标签的usage_count
	if len(actuallyAddedTagIDs) > 0 {
		for _, tagID := range actuallyAddedTagIDs {
			var count int64
			if err := s.db.Model(&models.FileGlobalTagRelation{}).
				Where("tag_id = ?", tagID).
				Count(&count).Error; err != nil {
				continue
			}

			if err := s.db.Model(&models.GlobalTag{}).
				Where("id = ?", tagID).
				Update("usage_count", count).Error; err != nil {
			}
		}
	}

	return nil
}

/* RemoveTagsFromFile 移除文件的标签 */
func (s *FileGlobalTagService) RemoveTagsFromFile(fileID string, tagIDs []uint) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	if len(tagIDs) == 0 {
		return nil
	}

	// 删除标签关联
	err := s.db.Where("file_id = ? AND tag_id IN ?", fileID, tagIDs).Delete(&models.FileGlobalTagRelation{}).Error
	if err != nil {
		return fmt.Errorf("移除文件标签失败: %v", err)
	}

	// 手动批量更新受影响标签的usage_count
	if len(tagIDs) > 0 {
		for _, tagID := range tagIDs {
			var count int64
			if err := s.db.Model(&models.FileGlobalTagRelation{}).
				Where("tag_id = ?", tagID).
				Count(&count).Error; err != nil {
				continue
			}

			if err := s.db.Model(&models.GlobalTag{}).
				Where("id = ?", tagID).
				Update("usage_count", count).Error; err != nil {
			}
		}
	}

	return nil
}

/* GetFileTags 获取文件的所有标签 */
func (s *FileGlobalTagService) GetFileTags(fileID string) ([]models.GlobalTag, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var tags []models.GlobalTag
	err := s.db.Table("global_tag").
		Joins("JOIN file_global_tag_relation ON global_tag.id = file_global_tag_relation.tag_id").
		Where("file_global_tag_relation.file_id = ?", fileID).
		Order("file_global_tag_relation.confidence DESC, global_tag.usage_count DESC").
		Find(&tags).Error

	if err != nil {
		return nil, fmt.Errorf("获取文件标签失败: %v", err)
	}

	return tags, nil
}

/* GetFileTagsWithDetails 获取文件标签详细信息（包含关联信息） */
func (s *FileGlobalTagService) GetFileTagsWithDetails(fileID string) ([]models.FileGlobalTagRelation, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var relations []models.FileGlobalTagRelation
	err := s.db.Where("file_id = ?", fileID).
		Preload("Tag").
		Order("confidence DESC, created_at DESC").
		Find(&relations).Error

	if err != nil {
		return nil, fmt.Errorf("获取文件标签详情失败: %v", err)
	}

	return relations, nil
}

/* ReplaceFileTags 替换文件的所有标签 */
func (s *FileGlobalTagService) ReplaceFileTags(fileID string, newTagIDs []uint, source string, confidence float64) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		return s.ReplaceFileTagsTx(tx, fileID, newTagIDs, source, confidence)
	})
}

/* ReplaceFileTagsTx 使用外部事务替换文件的所有标签 */
func (s *FileGlobalTagService) ReplaceFileTagsTx(tx *gorm.DB, fileID string, newTagIDs []uint, source string, confidence float64) error {
	if tx == nil {
		return fmt.Errorf("事务不能为空")
	}

	// 1. 获取文件信息
	var file models.File
	if err := tx.Where("id = ?", fileID).First(&file).Error; err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	// 2. 获取当前已有的标签
	var existingRelations []models.FileGlobalTagRelation
	if err := tx.Where("file_id = ?", fileID).Find(&existingRelations).Error; err != nil {
		return fmt.Errorf("获取现有标签关联失败: %v", err)
	}

	// 3. 构建旧标签集合和新标签集合
	existingMap := make(map[uint]bool)
	for _, rel := range existingRelations {
		existingMap[rel.TagID] = true
	}

	newMap := make(map[uint]bool)
	for _, tagID := range newTagIDs {
		newMap[tagID] = true
	}

	// 4. 计算需要删除的标签（在旧集合中但不在新集合中）
	var tagsToRemove []uint
	for tagID := range existingMap {
		if !newMap[tagID] {
			tagsToRemove = append(tagsToRemove, tagID)
		}
	}

	// 5. 计算需要添加的标签（在新集合中但不在旧集合中）
	var tagsToAdd []uint
	for tagID := range newMap {
		if !existingMap[tagID] {
			tagsToAdd = append(tagsToAdd, tagID)
		}
	}

	// 6. 删除不再需要的标签关联
	if len(tagsToRemove) > 0 {
		if err := tx.Where("file_id = ? AND tag_id IN ?", fileID, tagsToRemove).
			Delete(&models.FileGlobalTagRelation{}).Error; err != nil {
			return fmt.Errorf("删除旧标签关联失败: %v", err)
		}
	}

	// 7. 添加新的标签关联
	if len(tagsToAdd) > 0 {
		newRelations := make([]models.FileGlobalTagRelation, 0, len(tagsToAdd))
		for _, tagID := range tagsToAdd {
			newRelations = append(newRelations, models.FileGlobalTagRelation{
				FileID:      fileID,
				TagID:       tagID,
				UserID:      file.UserID,
				AccessLevel: file.AccessLevel,
				Source:      source,
				Confidence:  confidence,
			})
		}

		if err := tx.Create(&newRelations).Error; err != nil {
			return fmt.Errorf("创建新标签关联失败: %v", err)
		}
	}

	// 8. 手动批量更新受影响标签的usage_count
	affectedTagIDs := make(map[uint]bool)
	for _, tagID := range tagsToRemove {
		affectedTagIDs[tagID] = true
	}
	for _, tagID := range tagsToAdd {
		affectedTagIDs[tagID] = true
	}

	if len(affectedTagIDs) > 0 {
		tagIDList := make([]uint, 0, len(affectedTagIDs))
		for tagID := range affectedTagIDs {
			tagIDList = append(tagIDList, tagID)
		}

		// 批量更新每个受影响标签的usage_count
		for _, tagID := range tagIDList {
			var count int64
			if err := tx.Model(&models.FileGlobalTagRelation{}).
				Where("tag_id = ?", tagID).
				Count(&count).Error; err != nil {
				continue
			}

			if err := tx.Model(&models.GlobalTag{}).
				Where("id = ?", tagID).
				Update("usage_count", count).Error; err != nil {
			}
		}
	}

	return nil
}

/* BuildTagsForAI 为AI构建标签列表，支持游客场景（userID可为0） */
func (s *FileGlobalTagService) BuildTagsForAI(userID uint, categoryID *uint) ([]ai.TagInfo, error) {
	return s.getTagsForAI(userID, categoryID)
}

/* ProcessAITaggingResult 处理AI标注结果 */
func (s *FileGlobalTagService) ProcessAITaggingResult(fileID string, aiResult *ai.FileAnalysisResponse) error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	if aiResult == nil || len(aiResult.Tags) == 0 {
		return nil
	}

	var file models.File
	err := s.db.Where("id = ?", fileID).First(&file).Error
	if err != nil {
		return fmt.Errorf("获取文件信息失败: %v", err)
	}

	globalTagService := NewGlobalTagService()

	var tagIDs []uint
	for _, tagName := range aiResult.Tags {
		if tagName == "" {
			continue
		}

		tag, err := globalTagService.CreateOrGetGlobalTag(tagName, "", file.UserID, false)
		if err != nil {
			logger.Warn("处理AI标签失败 [%s]: %v", tagName, err)
			continue
		}

		err = globalTagService.AddUserTagReference(file.UserID, tag.ID, "ai")
		if err != nil {
			logger.Warn("添加AI标签用户引用失败 [%s]: %v", tagName, err)
		}

		tagIDs = append(tagIDs, tag.ID)
	}

	if len(tagIDs) == 0 {
		return nil
	}

	err = s.AddTagsToFile(fileID, tagIDs, "ai", common.AIDefaultConfidence) // AI标注置信度
	if err != nil {
		return fmt.Errorf("添加AI标签到文件失败: %v", err)
	}

	return nil
}

/* GetTaggedImages 获取包含指定标签的文件 */
func (s *FileGlobalTagService) GetTaggedFiles(tagIDs []uint, userID uint, accessLevel string, limit, offset int) ([]models.File, int64, error) {
	if s.db == nil {
		return nil, 0, fmt.Errorf("数据库连接失败")
	}

	if len(tagIDs) == 0 {
		return []models.File{}, 0, nil
	}

	query := s.db.Table("file").
		Joins("JOIN file_global_tag_relation ON file.id = file_global_tag_relation.file_id").
		Where("file_global_tag_relation.tag_id IN ?", tagIDs)

	if userID > 0 {
		query = query.Where("file.user_id = ?", userID)
	}
	if accessLevel != "" {
		query = query.Where("file.access_level = ?", accessLevel)
	}

	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("获取标签文件数量失败: %v", err)
	}

	var images []models.File
	err = query.Distinct("id").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&images).Error

	if err != nil {
		return nil, 0, fmt.Errorf("获取标签文件失败: %v", err)
	}

	return images, total, nil
}

/* TagResultDB 统一查询标签结果结构 */
type TagResultDB struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Source      string `json:"source"` // user_custom, category_system, system_popular
	UsageCount  int    `json:"usage_count"`
}

func (s *FileGlobalTagService) getTagsForAI(userID uint, categoryID *uint) ([]ai.TagInfo, error) {
	if s.db == nil {
		return nil, fmt.Errorf("数据库连接失败")
	}

	var allTags []ai.TagInfo
	tagIDMap := make(map[uint]bool) // 用于去重

	var userTagsAdded, categoryTagsAdded, systemTagsAdded int

	if userID > 0 {
		var userTags []models.GlobalTag
		err := s.db.Table("global_tag").
			Joins("JOIN user_tag_reference ON global_tag.id = user_tag_reference.tag_id").
			Where("user_tag_reference.user_id = ?", userID).
			Order("global_tag.usage_count DESC").
			Find(&userTags).Error

		if err != nil {
			logger.Warn("获取用户标签失败: %v", err)
		} else {
			for _, tag := range userTags {
				if !tagIDMap[tag.ID] && len(allTags) < common.AIMaxTagsLimit {
					allTags = append(allTags, ai.TagInfo{
						ID:          tag.ID,
						Name:        tag.Name,
						Description: tag.Description,
						Source:      "user_custom",
						UsageCount:  tag.UsageCount,
					})
					tagIDMap[tag.ID] = true
					userTagsAdded++
				}
			}
		}
	}

	if categoryID != nil && len(allTags) < common.AIMaxTagsLimit {
		var categoryTags []models.GlobalTag
		err := s.db.Table("global_tag").
			Joins("JOIN file_global_tag_relation ON global_tag.id = file_global_tag_relation.tag_id").
			Joins("JOIN file_category_relation ON file_global_tag_relation.file_id = file_category_relation.file_id").
			Where("file_category_relation.category_id = ? AND global_tag.usage_count > 0", *categoryID).
			Group("global_tag.id").
			Order("global_tag.usage_count DESC").
			Find(&categoryTags).Error

		if err != nil {
			logger.Warn("获取分类标签失败: %v", err)
		} else {
			for _, tag := range categoryTags {
				if !tagIDMap[tag.ID] && len(allTags) < common.AIMaxTagsLimit {
					allTags = append(allTags, ai.TagInfo{
						ID:          tag.ID,
						Name:        tag.Name,
						Description: tag.Description,
						Source:      "category_system",
						UsageCount:  tag.UsageCount,
					})
					tagIDMap[tag.ID] = true
					categoryTagsAdded++
				}
			}
		}
	}

	if len(allTags) < 500 {
		var systemTags []models.GlobalTag
		limit := 500 - len(allTags)
		if limit > 500 {
			limit = 500
		}

		err := s.db.Where("usage_count > 0").
			Order("usage_count DESC").
			Limit(limit).
			Find(&systemTags).Error

		if err != nil {
			logger.Warn("获取系统标签失败: %v", err)
		} else {
			for _, tag := range systemTags {
				if !tagIDMap[tag.ID] && len(allTags) < common.AIMaxTagsLimit {
					allTags = append(allTags, ai.TagInfo{
						ID:          tag.ID,
						Name:        tag.Name,
						Description: tag.Description,
						Source:      "system_popular",
						UsageCount:  tag.UsageCount,
					})
					tagIDMap[tag.ID] = true
					systemTagsAdded++
				}
			}
		}
	}

	return allTags, nil
}

/* CalibrateAllTagUsageCount 校准所有标签的usage_count，用于定时任务 */
func (s *FileGlobalTagService) CalibrateAllTagUsageCount() error {
	if s.db == nil {
		return fmt.Errorf("数据库连接失败")
	}

	// 获取所有标签
	var tags []models.GlobalTag
	if err := s.db.Find(&tags).Error; err != nil {
		return fmt.Errorf("获取标签列表失败: %v", err)
	}


	updated := 0
	errors := 0

	// 批量校准每个标签的usage_count
	for _, tag := range tags {
		var count int64
		if err := s.db.Model(&models.FileGlobalTagRelation{}).
			Where("tag_id = ?", tag.ID).
			Count(&count).Error; err != nil {
			errors++
			continue
		}

		// 只更新不一致的记录
		if int(count) != tag.UsageCount {
			if err := s.db.Model(&models.GlobalTag{}).
				Where("id = ?", tag.ID).
				Update("usage_count", count).Error; err != nil {
				errors++
			} else {
				updated++
			}
		}
	}


	if errors > 0 {
		return fmt.Errorf("校准过程中发生 %d 个错误", errors)
	}

	return nil
}
