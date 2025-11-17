package category

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
)

/* TemplateService 分类模板服务 */
type TemplateService struct {
	db *gorm.DB
}

func (s *TemplateService) getDB() (*gorm.DB, error) {
	if s.db == nil {
		s.db = database.DB
	}
	if s.db == nil {
		logger.Error("分类模板服务尝试访问未初始化的数据库实例")
		return nil, errors.New(errors.CodeInternal, "数据库未初始化")
	}
	return s.db, nil
}

/* NewTemplateService 创建分类模板服务实例 */
func NewTemplateService() *TemplateService {
	return &TemplateService{
		db: database.DB,
	}
}

/* CreateTemplateRequest 创建分类模板请求 */
type CreateTemplateRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	IsPopular   bool   `json:"is_popular"`
	SortOrder   int    `json:"sort_order"`
}

/* UpdateTemplateRequest 更新分类模板请求 */
type UpdateTemplateRequest struct {
	Name        string `json:"name" binding:"omitempty,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=500"`
	Icon        string `json:"icon" binding:"omitempty,max=50"`
	IsPopular   *bool  `json:"is_popular"`
	SortOrder   *int   `json:"sort_order"`
}

/* TemplateListQuery 模板列表查询参数 */
type TemplateListQuery struct {
	Keyword   string `form:"keyword"`
	IsPopular *bool  `form:"is_popular"`
	Page      int    `form:"page" binding:"omitempty,min=1"`
	Size      int    `form:"size" binding:"omitempty,min=1,max=100"`
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=name sort_order usage_count created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

/* TemplateListResponse 模板列表响应 */
type TemplateListResponse struct {
	Templates []models.CategoryTemplate `json:"templates"`
	Total     int64                     `json:"total"`
	Page      int                       `json:"page"`
	Size      int                       `json:"size"`
}

/* CreateTemplate 创建系统分类模板 */
func (s *TemplateService) CreateTemplate(req CreateTemplateRequest) (*models.CategoryTemplate, error) {

	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	var existingTemplate models.CategoryTemplate
	err = db.Where("name = ?", req.Name).First(&existingTemplate).Error
	if err == nil {
		return nil, errors.New(errors.CodeDBDuplicate, "分类模板名称已存在")
	}
	if err != gorm.ErrRecordNotFound {
		logger.Error("检查分类模板名称失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "检查分类模板名称失败")
	}

	template := &models.CategoryTemplate{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		IsPopular:   req.IsPopular,
		SortOrder:   req.SortOrder,
		UsageCount:  0,
	}

	if req.SortOrder == 0 {
		var maxOrder int
		db.Model(&models.CategoryTemplate{}).
			Select("COALESCE(MAX(sort_order), 0)").
			Scan(&maxOrder)
		template.SortOrder = maxOrder + 1
	}

	if err := db.Create(template).Error; err != nil {
		logger.Error("创建分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBCreateFailed, "创建分类模板失败")
	}

	return template, nil
}

/* GetTemplate 获取单个分类模板 */
func (s *TemplateService) GetTemplate(id uint) (*models.CategoryTemplate, error) {

	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	var template models.CategoryTemplate
	err = db.First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "分类模板不存在")
		}
		logger.Error("获取分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类模板失败")
	}

	return &template, nil
}

/* UpdateTemplate 更新分类模板 */
func (s *TemplateService) UpdateTemplate(id uint, req UpdateTemplateRequest) (*models.CategoryTemplate, error) {

	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	var template models.CategoryTemplate
	err = db.First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.CodeNotFound, "分类模板不存在")
		}
		logger.Error("获取分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类模板失败")
	}

	if req.Name != "" && req.Name != template.Name {
		var existingTemplate models.CategoryTemplate
		err = db.Where("name = ? AND id != ?", req.Name, id).First(&existingTemplate).Error
		if err == nil {
			return nil, errors.New(errors.CodeDBDuplicate, "分类模板名称已存在")
		}
		if err != gorm.ErrRecordNotFound {
			logger.Error("检查分类模板名称失败: %v", err)
			return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "检查分类模板名称失败")
		}
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.IsPopular != nil {
		updates["is_popular"] = *req.IsPopular
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) > 0 {
		if err := db.Model(&template).Updates(updates).Error; err != nil {
			logger.Error("更新分类模板失败: %v", err)
			return nil, errors.Wrap(err, errors.CodeDBUpdateFailed, "更新分类模板失败")
		}
	}

	err = db.First(&template, id).Error
	if err != nil {
		logger.Error("获取更新后的分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取更新后的分类模板失败")
	}

	return &template, nil
}

/* DeleteTemplate 删除分类模板 */
func (s *TemplateService) DeleteTemplate(id uint) error {

	db, err := s.getDB()
	if err != nil {
		return err
	}

	var template models.CategoryTemplate
	err = db.First(&template, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.CodeNotFound, "分类模板不存在")
		}
		logger.Error("获取分类模板失败: %v", err)
		return errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类模板失败")
	}

	var userCategoryCount int64
	err = db.Model(&models.FileCategory{}).
		Where("template_id = ?", id).
		Count(&userCategoryCount).Error
	if err != nil {
		logger.Error("检查模板使用情况失败: %v", err)
		return errors.Wrap(err, errors.CodeDBQueryFailed, "检查模板使用情况失败")
	}

	if userCategoryCount > 0 {
		return errors.New(errors.CodeConflict, fmt.Sprintf("该模板正在被 %d 个用户分类使用，无法删除", userCategoryCount))
	}

	if err := db.Delete(&template).Error; err != nil {
		logger.Error("删除分类模板失败: %v", err)
		return errors.Wrap(err, errors.CodeDBDeleteFailed, "删除分类模板失败")
	}

	return nil
}

/* ListTemplates 获取分类模板列表 */
func (s *TemplateService) ListTemplates(query TemplateListQuery) (*TemplateListResponse, error) {
	db, err := s.getDB()
	if err != nil {
		logger.Warn("数据库未初始化，返回空的分类模板列表")
		return &TemplateListResponse{
			Templates: []models.CategoryTemplate{},
			Total:     0,
			Page:      query.Page,
			Size:      query.Size,
		}, nil
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	if query.SortBy == "" {
		query.SortBy = "sort_order"
	}
	if query.SortOrder == "" {
		query.SortOrder = "asc"
	}

	db = db.Model(&models.CategoryTemplate{})

	if query.Keyword != "" {
		keyword := "%" + query.Keyword + "%"
		db = db.Where("name LIKE ? OR description LIKE ?", keyword, keyword)
	}

	if query.IsPopular != nil {
		db = db.Where("is_popular = ?", *query.IsPopular)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		logger.Error("统计分类模板数量失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "统计分类模板数量失败")
	}

	orderClause := fmt.Sprintf("%s %s", query.SortBy, query.SortOrder)
	var templates []models.CategoryTemplate
	err = db.Order(orderClause).
		Offset((query.Page - 1) * query.Size).
		Limit(query.Size).
		Find(&templates).Error

	if err != nil {
		logger.Error("获取分类模板列表失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类模板列表失败")
	}

	response := &TemplateListResponse{
		Templates: templates,
		Total:     total,
		Page:      query.Page,
		Size:      query.Size,
	}

	return response, nil
}

/* GetPopularTemplates 获取热门分类模板（用于AI提示） */
func (s *TemplateService) GetPopularTemplates(limit int) ([]models.CategoryTemplate, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 20
	}

	var templates []models.CategoryTemplate
	err = db.Where("is_popular = ?", true).
		Order("usage_count DESC, sort_order ASC").
		Limit(limit).
		Find(&templates).Error

	if err != nil {
		logger.Error("获取热门分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取热门分类模板失败")
	}

	return templates, nil
}

/* GetAllTemplatesForAI 获取所有模板用于AI提示（按使用频率排序） */
func (s *TemplateService) GetAllTemplatesForAI() ([]models.CategoryTemplate, error) {
	db, err := s.getDB()
	if err != nil {
		return nil, err
	}

	var templates []models.CategoryTemplate
	err = db.Order("usage_count DESC, is_popular DESC, sort_order ASC").
		Find(&templates).Error

	if err != nil {
		logger.Error("获取分类模板失败: %v", err)
		return nil, errors.Wrap(err, errors.CodeDBQueryFailed, "获取分类模板失败")
	}

	return templates, nil
}

/* IncrementUsageCount 增加模板使用次数（当用户采纳模板时调用） */
func (s *TemplateService) IncrementUsageCount(templateID uint) error {
	db, err := s.getDB()
	if err != nil {
		return err
	}

	err = db.Model(&models.CategoryTemplate{}).
		Where("id = ?", templateID).
		UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error

	if err != nil {
		logger.Error("增加模板使用次数失败: %v", err)
		return errors.Wrap(err, errors.CodeDBUpdateFailed, "增加模板使用次数失败")
	}

	return nil
}

/* BatchUpdateSortOrder 批量更新排序 */
func (s *TemplateService) BatchUpdateSortOrder(sortOrders map[uint]int) error {
	db, err := s.getDB()
	if err != nil {
		return err
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for templateID, sortOrder := range sortOrders {
			err := tx.Model(&models.CategoryTemplate{}).
				Where("id = ?", templateID).
				Update("sort_order", sortOrder).Error
			if err != nil {
				logger.Error("更新模板排序失败: templateID=%d, sortOrder=%d, err=%v", templateID, sortOrder, err)
				return err
			}
		}
		return nil
	})
}
