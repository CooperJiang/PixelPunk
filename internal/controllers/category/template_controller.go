package category

import (
	"pixelpunk/internal/controllers/category/dto"
	categoryService "pixelpunk/internal/services/category"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TemplateController struct {
	templateService *categoryService.TemplateService
}

func NewTemplateController() *TemplateController {
	return &TemplateController{
		templateService: categoryService.NewTemplateService(),
	}
}

func (c *TemplateController) CreateTemplate(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.CreateTemplateDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	serviceReq := categoryService.CreateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		IsPopular:   req.IsPopular,
		SortOrder:   req.SortOrder,
	}

	template, err := c.templateService.CreateTemplate(serviceReq)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, template, "分类模板创建成功")
}

func (c *TemplateController) GetTemplate(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "模板ID不能为空"))
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(ctx, errors.New(errors.CodeInvalidParameter, "无效的模板ID"))
		return
	}

	template, err := c.templateService.GetTemplate(uint(id))
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, template, "获取成功")
}

func (c *TemplateController) UpdateTemplate(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.UpdateTemplateDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	serviceReq := categoryService.UpdateTemplateRequest{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		IsPopular:   req.IsPopular,
		SortOrder:   req.SortOrder,
	}

	template, err := c.templateService.UpdateTemplate(req.ID, serviceReq)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, template, "分类模板更新成功")
}

func (c *TemplateController) DeleteTemplate(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.DeleteTemplateDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	err = c.templateService.DeleteTemplate(req.ID)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, nil, "分类模板删除成功")
}

func (c *TemplateController) ListTemplates(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.TemplateListQueryDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	serviceQuery := categoryService.TemplateListQuery{
		Keyword:   req.Keyword,
		IsPopular: req.IsPopular,
		Page:      req.Page,
		Size:      req.Size,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	response, err := c.templateService.ListTemplates(serviceQuery)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, response, "获取成功")
}

func (c *TemplateController) GetPopularTemplates(ctx *gin.Context) {
	limitStr := ctx.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	templates, err := c.templateService.GetPopularTemplates(limit)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	response := map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
	}

	errors.ResponseSuccess(ctx, response, "获取成功")
}

func (c *TemplateController) GetAllTemplates(ctx *gin.Context) {
	templates, err := c.templateService.GetAllTemplatesForAI()
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	response := map[string]interface{}{
		"templates": templates,
		"count":     len(templates),
	}

	errors.ResponseSuccess(ctx, response, "获取成功")
}

func (c *TemplateController) BatchUpdateSortOrder(ctx *gin.Context) {
	req, err := common.ValidateRequest[dto.BatchSortOrderDTO](ctx)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	sortOrders := make(map[uint]int)
	for _, item := range req.SortOrders {
		sortOrders[item.ID] = item.SortOrder
	}

	err = c.templateService.BatchUpdateSortOrder(sortOrders)
	if err != nil {
		errors.HandleError(ctx, err)
		return
	}

	errors.ResponseSuccess(ctx, nil, "排序更新成功")
}
