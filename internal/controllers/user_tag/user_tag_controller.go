package user_tag

import (
	"pixelpunk/internal/controllers/user_tag/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/user_tag"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetUserTagList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.GetUserTagListQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 50
	}
	if req.SortBy == "" {
		req.SortBy = "file_count"
	}
	if req.SortOrder == "" {
		req.SortOrder = "desc"
	}

	tags, total, err := user_tag.GetUserTagList(userID, req.Keyword, req.SortBy, req.SortOrder, req.Page, req.Size)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"items": tags,
		"total": total,
		"pagination": gin.H{
			"page": req.Page,
			"size": req.Size,
		},
	}, "获取标签列表成功")
}

func GetAllUserTags(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	tags, err := user_tag.GetAllUserTags(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, tags, "获取标签列表成功")
}

func CreateUserTag(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.CreateUserTagDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	tag, err := user_tag.CreateUserTag(userID, req.Name)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, tag, "创建标签成功")
}

func UpdateUserTag(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.UpdateUserTagDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	tag, err := user_tag.UpdateUserTag(userID, req.ID, req.Name)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, tag, "更新标签成功")
}

func DeleteUserTag(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.DeleteUserTagDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = user_tag.DeleteUserTag(userID, req.ID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除标签成功")
}

func BatchDeleteUserTags(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.BatchDeleteUserTagsDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	count, err := user_tag.BatchDeleteUserTags(userID, req.IDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"deleted_count": count,
	}, "批量删除标签成功")
}

func MergeUserTags(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.MergeUserTagsDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = user_tag.MergeUserTags(userID, req.SourceIDs, req.TargetID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "合并标签成功")
}

func GetUserTagStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	stats, err := user_tag.GetUserTagStats(userID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, stats, "获取统计成功")
}
