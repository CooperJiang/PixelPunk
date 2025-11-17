package tag

import (
	"fmt"
	"pixelpunk/internal/controllers/tag/dto"
	"pixelpunk/internal/middleware"
	tagService "pixelpunk/internal/services/tag"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

type GlobalTagController struct {
	globalTagService *tagService.GlobalTagService
	fileTagService   *tagService.FileGlobalTagService
}

func NewGlobalTagController() *GlobalTagController {
	return &GlobalTagController{
		globalTagService: tagService.NewGlobalTagService(),
		fileTagService:   tagService.NewFileGlobalTagService(),
	}
}

func (gc *GlobalTagController) DeleteTag(c *gin.Context) {
	tagIDStr := c.Param("tag_id")
	if tagIDStr == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "标签ID不能为空"))
		return
	}

	tagID, err := strconv.ParseUint(tagIDStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的标签ID"))
		return
	}

	if err := gc.globalTagService.DeleteGlobalTag(uint(tagID)); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("删除标签失败: %v", err)))
		return
	}

	errors.ResponseSuccess(c, nil, "标签删除成功")
}

func (gc *GlobalTagController) CreateTag(c *gin.Context) {
	user := middleware.GetCurrentUser(c)
	if user == nil {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	var req dto.CreateTagDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	tag, err := gc.globalTagService.CreateOrGetGlobalTag(req.Name, req.Description, user.UserID, false)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("创建标签失败: %v", err)))
		return
	}

	if req.SortOrder != nil {
		if upd, updErr := gc.globalTagService.UpdateGlobalTag(tag.ID, nil, nil, req.SortOrder); updErr == nil {
			tag = upd
		} else {
			logger.Warn("更新标签排序失败: %v", updErr)
		}
	}

	respItem, _ := gc.globalTagService.ConvertToResponseItem(*tag)
	errors.ResponseSuccess(c, respItem, "标签创建成功")
}

func (gc *GlobalTagController) UpdateTag(c *gin.Context) {
	_ = middleware.GetCurrentUserID(c) // 保持和其它接口一致（已由中间件校验）

	var req dto.UpdateTagDTO
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	var namePtr *string
	var descPtr *string
	if req.Name != "" {
		namePtr = &req.Name
	}
	if req.Description != "" {
		descPtr = &req.Description
	}

	tag, err := gc.globalTagService.UpdateGlobalTag(req.ID, namePtr, descPtr, req.SortOrder)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("更新标签失败: %v", err)))
		return
	}

	respItem, _ := gc.globalTagService.ConvertToResponseItem(*tag)
	errors.ResponseSuccess(c, respItem, "标签更新成功")
}

func (gc *GlobalTagController) ListTags(c *gin.Context) {
	userInfo := middleware.GetCurrentUser(c)
	if userInfo == nil {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "50")
	keyword := c.Query("keyword")
	sortBy := c.DefaultQuery("sort_by", "usage_count")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	size, err := strconv.Atoi(sizeStr)
	if err != nil || size <= 0 {
		size = 50
	}
	if size > 1000 {
		size = 1000
	}

	result, err := gc.globalTagService.GetTagsWithPagination(userInfo.UserID, keyword, sortBy, sortOrder, page, size)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取标签列表失败: %v", err)))
		return
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"items": result.Items,
		"total": result.Total,
		"pagination": map[string]interface{}{
			"page": page,
			"size": size,
		},
	}, "获取标签列表成功")
}

func (gc *GlobalTagController) GetFileTags(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	tags, err := gc.fileTagService.GetFileTags(fileID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取文件标签失败: %v", err)))
		return
	}

	var responseTags []tagService.GlobalTagResponseItem
	for _, tag := range tags {
		responseItem, err := gc.globalTagService.ConvertToResponseItem(tag)
		if err != nil {
			logger.Warn("转换标签响应项失败: %v", err)
			responseItem = tagService.GlobalTagResponseItem{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				IsSystem:    tag.IsSystem,
				UsageCount:  tag.UsageCount,
				Count:       0,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			}
		}
		responseTags = append(responseTags, responseItem)
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"tags": responseTags,
	}, "获取文件标签成功")
}

func (gc *GlobalTagController) AddTagsToFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	userInfo := middleware.GetCurrentUser(c)
	if userInfo == nil {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "用户未登录"))
		return
	}

	var req dto.AddTagsToFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	if len(req.TagNames) == 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "标签名称不能为空"))
		return
	}

	globalTags, err := gc.globalTagService.CreateTagsFromNames(req.TagNames, userInfo.UserID, "manual")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("创建标签失败: %v", err)))
		return
	}

	var tagIDs []uint
	for _, tag := range globalTags {
		tagIDs = append(tagIDs, tag.ID)
	}

	err = gc.fileTagService.AddTagsToFile(fileID, tagIDs, "manual", 1.0)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("添加标签到文件失败: %v", err)))
		return
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"added_count": len(tagIDs),
	}, "添加标签成功")
}

func (gc *GlobalTagController) RemoveTagsFromFile(c *gin.Context) {
	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	var req dto.RemoveTagsFromFileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, fmt.Sprintf("请求参数错误: %v", err)))
		return
	}

	if len(req.TagIDs) == 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "标签ID不能为空"))
		return
	}

	err := gc.fileTagService.RemoveTagsFromFile(fileID, req.TagIDs)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("移除文件标签失败: %v", err)))
		return
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"removed_count": len(req.TagIDs),
	}, "移除标签成功")
}

func (gc *GlobalTagController) SearchTags(c *gin.Context) {
	keyword := c.Query("keyword")
	if keyword == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "搜索关键词不能为空"))
		return
	}

	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	tags, err := gc.globalTagService.SearchTags(keyword, limit)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索标签失败: %v", err)))
		return
	}

	var responseTags []tagService.GlobalTagResponseItem
	for _, tag := range tags {
		responseItem, err := gc.globalTagService.ConvertToResponseItem(tag)
		if err != nil {
			logger.Warn("转换标签响应项失败: %v", err)
			responseItem = tagService.GlobalTagResponseItem{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				IsSystem:    tag.IsSystem,
				UsageCount:  tag.UsageCount,
				Count:       0,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			}
		}
		responseTags = append(responseTags, responseItem)
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"items": responseTags,
		"total": len(responseTags),
		"pagination": map[string]interface{}{
			"page": 1,
			"size": limit,
		},
	}, "搜索标签成功")
}

func (gc *GlobalTagController) GetTagStats(c *gin.Context) {
	popularTags, err := gc.globalTagService.GetPopularTags(10)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取热门标签失败: %v", err)))
		return
	}

	var responsePopularTags []tagService.GlobalTagResponseItem
	for _, tag := range popularTags {
		responseItem, err := gc.globalTagService.ConvertToResponseItem(tag)
		if err != nil {
			logger.Warn("转换标签响应项失败: %v", err)
			responseItem = tagService.GlobalTagResponseItem{
				ID:          tag.ID,
				Name:        tag.Name,
				Description: tag.Description,
				IsSystem:    tag.IsSystem,
				UsageCount:  tag.UsageCount,
				Count:       0,
				CreatedAt:   tag.CreatedAt,
				UpdatedAt:   tag.UpdatedAt,
			}
		}
		responsePopularTags = append(responsePopularTags, responseItem)
	}

	errors.ResponseSuccess(c, map[string]interface{}{
		"popular_tags": responsePopularTags,
		"stats": map[string]interface{}{
			"total_tags": len(responsePopularTags),
		},
	}, "获取标签统计成功")
}
