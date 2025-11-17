package tag

import (
	"fmt"
	"pixelpunk/internal/controllers/tag/dto"
	"pixelpunk/internal/middleware"
	tagService "pixelpunk/internal/services/tag"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

type TagController struct{ globalController *GlobalTagController }

func NewTagController() *TagController {
	return &TagController{globalController: NewGlobalTagController()}
}

var Controller *TagController

func getController() *TagController {
	if Controller == nil {
		Controller = NewTagController()
	}
	return Controller
}

func ListTags(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.ListTags(c)
}

func AddTagsToFile(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.AddTagsToFile(c)
}

func RemoveTagsFromFile(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.RemoveTagsFromFile(c)
}

func GetFileTags(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.GetFileTags(c)
}

func GetTagStats(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.GetTagStats(c)
}

func SearchTags(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.SearchTags(c)
}

func CreateTag(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.CreateTag(c)
}

func UpdateTag(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.UpdateTag(c)
}

func DeleteTag(c *gin.Context) {
	globalController := NewGlobalTagController()
	globalController.DeleteTag(c)
}

func BatchOperateTags(c *gin.Context) {
	getController().BatchOperateTags(c)
}

func GetDetailedTagStats(c *gin.Context) {
	getController().GetDetailedTagStats(c)
}

func GetTagAnalytics(c *gin.Context) {
	getController().GetTagAnalytics(c)
}

func (ctrl *TagController) BatchOperateTags(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.BatchOperateTagsDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	switch req.Action {
	case "delete":
		if ctrl.globalController == nil {
			ctrl.globalController = NewGlobalTagController()
		}

		for _, id := range req.TagIDs {
			if err := ctrl.globalController.globalTagService.DeleteGlobalTag(id); err != nil {
				errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("删除标签 %d 失败: %v", id, err)))
				return
			}
		}
		errors.ResponseSuccess(c, nil, "批量删除标签成功")

	case "merge":
		if req.TargetID == 0 {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "合并操作需要指定目标标签ID"))
			return
		}
		gs := tagService.NewGlobalTagService()
		if err := gs.MergeGlobalTags(userID, req.TagIDs, req.TargetID); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("合并标签失败: %v", err)))
			return
		}
		errors.ResponseSuccess(c, nil, "批量合并标签成功")

	case "transfer":
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "批量转移标签功能正在开发中"))
	default:
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "不支持的操作类型"))
	}
}

func (ctrl *TagController) GetDetailedTagStats(c *gin.Context) {
	if _, err := common.ValidateRequest[dto.GetTagStatsQueryDTO](c); err != nil {
		errors.HandleError(c, err)
		return
	}

	gs := tagService.NewGlobalTagService()
	pg, err := gs.GetTagsWithPagination(0, "", "usage_count", "desc", 1, 50)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取标签统计失败: %v", err)))
		return
	}

	var totalUsage int64
	var systemTagsCount int64
	var userTagsCount int64
	for _, t := range pg.Items {
		totalUsage += t.Count
		if t.IsSystem {
			systemTagsCount++
		} else {
			userTagsCount++
		}
	}

	var avgUsage float64
	if pg.Total > 0 {
		avgUsage = float64(totalUsage) / float64(pg.Total)
	}

	stats := map[string]interface{}{
		"summary": map[string]interface{}{
			"total_tags":    pg.Total,
			"system_tags":   systemTagsCount,
			"user_tags":     userTagsCount,
			"total_usage":   totalUsage,
			"average_usage": avgUsage,
		},
		"detailed_stats": pg.Items,
		"generated_at":   "now",
	}
	errors.ResponseSuccess(c, stats, "获取详细统计成功")
}

func (ctrl *TagController) GetTagAnalytics(c *gin.Context) {
	if _, err := common.ValidateRequest[dto.GetTagStatsQueryDTO](c); err != nil {
		errors.HandleError(c, err)
		return
	}

	gs := tagService.NewGlobalTagService()
	topPg, err := gs.GetTagsWithPagination(0, "", "usage_count", "desc", 1, 10)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("获取标签统计失败: %v", err)))
		return
	}
	allPg, err := gs.GetTagsWithPagination(0, "", "count", "asc", 1, 1000)
	if err != nil {
		allPg = &tagService.TagsPaginationResult{Items: []tagService.GlobalTagResponseItem{}, Total: 0}
	}

	var totalUsage int64
	for _, t := range topPg.Items {
		totalUsage += t.Count
	}
	var unused int64
	for _, t := range allPg.Items {
		if t.Count == 0 {
			unused++
		}
	}

	analytics := map[string]interface{}{
		"overview": map[string]interface{}{
			"total_tags":     allPg.Total,
			"total_usage":    totalUsage,
			"unused_tags":    unused,
			"top_tags_count": len(topPg.Items),
		},
		"top_tags": topPg.Items,
		"insights": []string{
			fmt.Sprintf("共有 %d 个标签", allPg.Total),
			fmt.Sprintf("Top10 使用总次数: %d", totalUsage),
			fmt.Sprintf("未使用的标签: %d 个", unused),
		},
		"generated_at": "now",
	}
	errors.ResponseSuccess(c, analytics, "获取分析数据成功")
}
