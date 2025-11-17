package user_category

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strings"

	"gorm.io/gorm"
)

/* CreateCategory 创建用户分类 */
func CreateCategory(userID uint, name, description string, sortOrder int) (*models.FileCategory, error) {
	db := database.GetDB()

	var existingCategory models.FileCategory
	if err := db.Where("user_id = ? AND name = ? AND status = ?", userID, name, "active").First(&existingCategory).Error; err == nil {
		return nil, errors.New(errors.CodeDBDuplicate, "分类名称已存在")
	}

	category := &models.FileCategory{
		Name:        name,
		Description: description,
		UserID:      userID,
		SortOrder:   sortOrder,
		Source:      "user",
		Status:      "active",
		FileCount:   0,
	}

	if err := db.Create(category).Error; err != nil {
		logger.Error("创建分类失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建分类失败")
	}

	return category, nil
}

/* GetCategoryList 获取用户分类列表（分页） */
func GetCategoryList(userID uint, keyword, status, sortBy, sortOrder string, page, pageSize int) ([]models.FileCategory, int64, error) {
	db := database.GetDB()

	query := db.Model(&models.FileCategory{}).Where("user_id = ?", userID)

	if keyword != "" {
		keyword = strings.TrimSpace(keyword)
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	if status != "" && (status == "active" || status == "archived") {
		query = query.Where("status = ?", status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logger.Error("获取分类总数失败: %v", err)
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类列表失败")
	}

	validSortBy := map[string]bool{
		"name":       true,
		"sort_order": true,
		"file_count": true,
		"created_at": true,
		"updated_at": true,
	}
	if !validSortBy[sortBy] {
		sortBy = "sort_order"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)
	query = query.Order(orderClause)

	offset := (page - 1) * pageSize
	query = query.Offset(offset).Limit(pageSize)

	var categories []models.FileCategory
	if err := query.Find(&categories).Error; err != nil {
		logger.Error("获取分类列表失败: %v", err)
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类列表失败")
	}

	return categories, total, nil
}

/* GetAllCategories 获取用户所有激活的分类（不分页） */
func GetAllCategories(userID uint) ([]models.FileCategory, error) {
	db := database.GetDB()

	var categories []models.FileCategory
	if err := db.Where("user_id = ? AND status = ?", userID, "active").Order("sort_order ASC, created_at DESC").Find(&categories).Error; err != nil {
		logger.Error("获取所有分类失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取所有分类失败")
	}

	return categories, nil
}

/* GetCategoryDetail 获取分类详情 */
func GetCategoryDetail(userID, categoryID uint) (*models.FileCategory, error) {
	db := database.GetDB()

	var category models.FileCategory
	if err := db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "分类不存在")
		}
		logger.Error("获取分类详情失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类详情失败")
	}

	return &category, nil
}

/* UpdateCategory 更新分类 */
func UpdateCategory(userID, categoryID uint, name, description string, sortOrder int) (*models.FileCategory, error) {
	db := database.GetDB()

	var category models.FileCategory
	if err := db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "分类不存在")
		}
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类失败")
	}

	// 如果名称改变，检查新名称是否已存在
	if name != category.Name {
		var existingCategory models.FileCategory
		if err := db.Where("user_id = ? AND name = ? AND status = ? AND id != ?", userID, name, "active", categoryID).First(&existingCategory).Error; err == nil {
			return nil, errors.New(errors.CodeDBDuplicate, "分类名称已存在")
		}
	}

	updates := map[string]interface{}{
		"name":        name,
		"description": description,
		"sort_order":  sortOrder,
	}

	if err := db.Model(&category).Updates(updates).Error; err != nil {
		logger.Error("更新分类失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新分类失败")
	}

	if err := db.Where("id = ?", categoryID).First(&category).Error; err != nil {
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取更新后的分类失败")
	}

	return &category, nil
}

/* DeleteCategory 删除分类 */
func DeleteCategory(userID, categoryID uint) error {
	db := database.GetDB()

	var category models.FileCategory
	if err := db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分类不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类失败")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		// 删除该分类的所有文件关联（将分类ID设置为null）
		if err := tx.Model(&models.File{}).Where("category_id = ?", categoryID).Update("category_id", nil).Error; err != nil {
			logger.Error("解除文件关联失败: %v", err)
			return errors.Wrap(err, errors.CodeDBUpdateFailed, "解除文件关联失败")
		}

		if err := tx.Delete(&category).Error; err != nil {
			logger.Error("删除分类失败: %v", err)
			return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分类失败")
		}

		return nil
	})
}

/* BatchDeleteCategories 批量删除分类 */
func BatchDeleteCategories(userID uint, categoryIDs []uint) (int, error) {
	db := database.GetDB()

	deletedCount := 0

	return deletedCount, db.Transaction(func(tx *gorm.DB) error {
		var categories []models.FileCategory
		if err := tx.Where("id IN ? AND user_id = ?", categoryIDs, userID).Find(&categories).Error; err != nil {
			return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类失败")
		}

		if len(categories) == 0 {
			return errors.New(errors.CodeNotFound, "未找到可删除的分类")
		}

		validIDs := make([]uint, len(categories))
		for i, cat := range categories {
			validIDs[i] = cat.ID
		}

		if err := tx.Model(&models.File{}).Where("category_id IN ?", validIDs).Update("category_id", nil).Error; err != nil {
			logger.Error("批量解除文件关联失败: %v", err)
			return errors.Wrap(err, errors.CodeDBUpdateFailed, "解除文件关联失败")
		}

		result := tx.Where("id IN ?", validIDs).Delete(&models.FileCategory{})
		if result.Error != nil {
			logger.Error("批量删除分类失败: %v", result.Error)
			return errors.Wrap(result.Error, errors.CodeDBDeleteFailed, "批量删除分类失败")
		}

		deletedCount = int(result.RowsAffected)
		return nil
	})
}

/* UpdateCategoryStatus 更新分类状态 */
func UpdateCategoryStatus(userID, categoryID uint, status string) error {
	db := database.GetDB()

	var category models.FileCategory
	if err := db.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分类不存在")
		}
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类失败")
	}

	if err := db.Model(&category).Update("status", status).Error; err != nil {
		logger.Error("更新分类状态失败: %v", err)
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新分类状态失败")
	}

	return nil
}

/* BatchUpdateSortOrder 批量更新排序 */
func BatchUpdateSortOrder(userID uint, updates map[uint]int) error {
	db := database.GetDB()

	categoryIDs := make([]uint, 0, len(updates))
	for id := range updates {
		categoryIDs = append(categoryIDs, id)
	}

	var count int64
	if err := db.Model(&models.FileCategory{}).Where("id IN ? AND user_id = ?", categoryIDs, userID).Count(&count).Error; err != nil {
		return errors.Wrap(err, errors.CodeDBQueryFailed, "验证分类失败")
	}

	if int(count) != len(categoryIDs) {
		return errors.New(errors.CodeNotFound, "部分分类不存在或无权限")
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for id, sortOrder := range updates {
			if err := tx.Model(&models.FileCategory{}).Where("id = ? AND user_id = ?", id, userID).Update("sort_order", sortOrder).Error; err != nil {
				logger.Error("更新分类排序失败: id=%d, err=%v", id, err)
				return errors.Wrap(err, errors.CodeDBUpdateFailed, "更新排序失败")
			}
		}
		return nil
	})
}
