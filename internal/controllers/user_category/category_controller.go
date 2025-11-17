package user_category

import (
	"pixelpunk/internal/controllers/user_category/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/user_category"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.CreateCategoryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	category, err := user_category.CreateCategory(userID, req.Name, req.Description, req.SortOrder)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, category, "创建分类成功")
}

func GetCategoryList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	keyword := c.Query("keyword")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "sort_order")
	sortOrder := c.DefaultQuery("sort_order_dir", "asc")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	categories, total, err := user_category.GetCategoryList(userID, keyword, status, sortBy, sortOrder, page, pageSize)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"categories": categories,
		"total":      total,
		"page":       page,
		"size":       pageSize,
	}, "获取分类列表成功")
}

func GetAllCategories(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	categories, err := user_category.GetAllCategories(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"categories": categories,
		"count":      len(categories),
	}, "获取分类列表成功")
}

func GetCategoryDetail(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的分类ID"))
		return
	}

	category, err := user_category.GetCategoryDetail(userID, uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, category, "获取分类详情成功")
}

func UpdateCategory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.UpdateCategoryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	category, err := user_category.UpdateCategory(userID, req.ID, req.Name, req.Description, req.SortOrder)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, category, "更新分类成功")
}

func DeleteCategory(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.DeleteCategoryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user_category.DeleteCategory(userID, req.ID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除分类成功")
}

func BatchDelete(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.BatchDeleteDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	deletedCount, err := user_category.BatchDeleteCategories(userID, req.IDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"deleted_count": deletedCount,
	}, "批量删除分类成功")
}

func UpdateStatus(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.UpdateStatusDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := user_category.UpdateCategoryStatus(userID, req.ID, req.Status); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新分类状态成功")
}

func BatchUpdateSortOrder(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.BatchSortOrderDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 转换为Service层所需的格式
	updates := make(map[uint]int)
	for _, item := range req.SortOrders {
		updates[item.ID] = item.SortOrder
	}

	if err := user_category.BatchUpdateSortOrder(userID, updates); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新排序成功")
}
