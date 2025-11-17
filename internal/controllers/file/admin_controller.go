package file

import (
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	filesvc "pixelpunk/internal/services/file"
	messageService "pixelpunk/internal/services/message"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strconv"
	"strings"

	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

// AdminGetFileListParams 管理员查询文件参数
type AdminGetFileListParams struct {
	Page          int     `form:"page" json:"page"`
	Size          int     `form:"size" json:"size"`
	Sort          string  `form:"sort" json:"sort"`
	Keyword       string  `form:"keyword" json:"keyword"`
	Tags          string  `form:"tags" json:"tags"`                     // 逗号分隔的标签字符串
	CategoryID    string  `form:"categoryId" json:"categoryId"`         // 逗号分隔的分类ID字符串
	DominantColor string  `form:"dominant_color" json:"dominant_color"` // 逗号分隔的颜色字符串
	Resolution    string  `form:"resolution" json:"resolution"`
	NSFWMinScore  float64 `form:"nsfw_min_score" json:"nsfw_min_score"`
	NSFWMaxScore  float64 `form:"nsfw_max_score" json:"nsfw_max_score"`
	IsNSFW        *bool   `form:"is_nsfw" json:"is_nsfw"`
	StorageType   string  `form:"storage_type" json:"storage_type"`
	MinWidth      int     `form:"min_width" json:"min_width"`
	MaxWidth      int     `form:"max_width" json:"max_width"`
	MinHeight     int     `form:"min_height" json:"min_height"`
	MaxHeight     int     `form:"max_height" json:"max_height"`
	UserID        uint    `form:"user_id" json:"user_id"`
	IsRecommended *bool   `form:"is_recommended" json:"is_recommended"`
}

// AdminGetFileList 管理员获取文件列表
func AdminGetFileList(c *gin.Context) {

	var params AdminGetFileListParams
	if err := c.ShouldBindQuery(&params); err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeInvalidParameter, "参数无效"))
		return
	}

	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}
	// 移除最大限制，允许用户无限加载
	// 为了防止过度请求，设置一个合理的上限
	if params.Size > 1000 {
		params.Size = 1000 // 限制单次请求最大1000条
	}

	var tagsArray []string
	if params.Tags != "" {
		tagsArray = strings.Split(params.Tags, ",")
	}

	var dominantColorsArray []string
	if params.DominantColor != "" {
		dominantColorsArray = strings.Split(params.DominantColor, ",")
	}

	var categoryIDsArray []string
	if params.CategoryID != "" {
		categoryIDsArray = strings.Split(params.CategoryID, ",")
	}

	searchParams := filesvc.AdminFileSearchParams{
		Page:          params.Page,
		Size:          params.Size,
		Sort:          params.Sort,
		Keyword:       params.Keyword,
		Tags:          tagsArray,
		CategoryIDs:   categoryIDsArray,
		DominantColor: dominantColorsArray,
		Resolution:    params.Resolution,
		NSFWMinScore:  params.NSFWMinScore,
		NSFWMaxScore:  params.NSFWMaxScore,
		IsNSFW:        params.IsNSFW,
		StorageType:   params.StorageType,
		MinWidth:      params.MinWidth,
		MaxWidth:      params.MaxWidth,
		MinHeight:     params.MinHeight,
		MaxHeight:     params.MaxHeight,
		UserID:        params.UserID,
		IsRecommended: params.IsRecommended,
	}

	files, total, err := filesvc.AdminGetFileList(searchParams)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"items": files,
		"pagination": gin.H{
			"total":        total,
			"size":         params.Size,
			"current_page": params.Page,
			"last_page":    (total + int64(params.Size) - 1) / int64(params.Size),
		},
	}

	errors.ResponseSuccess(c, data, "获取管理员文件列表成功")
}

// AdminGetTagList 管理员获取标签列表
func AdminGetTagList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	tags, total, err := filesvc.AdminGetTagList(page, size)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": tags,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回标签列表
		data = gin.H{
			"items": tags,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取标签列表成功")
}

// AdminGetColorList 管理员获取颜色值列表
func AdminGetColorList(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("size"))

	colors, total, err := filesvc.AdminGetColorList(page, size)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var data gin.H
	if page > 0 && size > 0 {
		data = gin.H{
			"items": colors,
			"pagination": gin.H{
				"total":        total,
				"size":         size,
				"current_page": page,
				"last_page":    (total + int64(size) - 1) / int64(size),
			},
		}
	} else {
		// 不包含分页信息，直接返回颜色列表
		data = gin.H{
			"items": colors,
			"total": total,
		}
	}

	errors.ResponseSuccess(c, data, "获取颜色列表成功")
}

// AdminRecommendFile 管理员推荐/取消推荐文件
func AdminRecommendFile(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminRecommendFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	fileInfo, err := filesvc.ToggleFileRecommendStatus(req.FileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var message string
	if fileInfo.IsRecommended {
		message = "文件已设为推荐"
	} else {
		message = "文件已取消推荐"
	}

	errors.ResponseSuccess(c, fileInfo, message)
}

// AdminDeleteFile 管理员删除文件
func AdminDeleteFile(c *gin.Context) {
	var req struct {
		FileID string `json:"file_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误，缺少 file_id"))
		return
	}

	if req.FileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	var fileRecord models.File
	if err := database.DB.Where("id = ?", req.FileID).First(&fileRecord).Error; err != nil {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件不存在"))
		return
	}

	err := filesvc.DeleteFile(fileRecord.UserID, req.FileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	adminID := middleware.GetCurrentUserID(c)

	// 记录管理员删除活动日志
	activity.LogAdminDelete(fileRecord.UserID, req.FileID, adminID)

	// 异步发送管理员删除文件通知
	go sendAdminDeleteNotification(fileRecord.UserID, fileRecord.ID, fileRecord.OriginalName)

	errors.ResponseSuccess(c, gin.H{"id": req.FileID}, "管理员删除文件成功")
}

// AdminBatchRecommendFiles 管理员批量推荐/取消推荐文件
func AdminBatchRecommendFiles(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminBatchRecommendFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(req.FileIDs) > 100 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "单次批量操作最多支持100个文件"))
		return
	}

	result, err := filesvc.BatchSetFileRecommendation(req.FileIDs, req.IsRecommended)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var message string
	if req.IsRecommended {
		message = fmt.Sprintf("批量推荐操作完成，成功 %d 个，失败 %d 个", result.SuccessCount, result.FailCount)
	} else {
		message = fmt.Sprintf("批量取消推荐操作完成，成功 %d 个，失败 %d 个", result.SuccessCount, result.FailCount)
	}

	errors.ResponseSuccess(c, result, message)
}

// AdminBatchDeleteFiles 管理员批量删除文件
func AdminBatchDeleteFiles(c *gin.Context) {
	req, err := common.ValidateRequest[dto.AdminBatchDeleteFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(req.FileIDs) > 100 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "单次批量操作最多支持100个文件"))
		return
	}

	result, err := filesvc.BatchDeleteFiles(req.FileIDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	message := fmt.Sprintf("批量删除操作完成，成功 %d 个，失败 %d 个", result.SuccessCount, result.FailCount)

	errors.ResponseSuccess(c, result, message)
}

// sendAdminDeleteNotification 发送管理员删除文件通知
func sendAdminDeleteNotification(userID uint, fileID, fileName string) {
	variables := map[string]interface{}{
		"file_id":      fileID,
		"file_name":    fileName,
		"reason":       "管理员删除",
		"related_type": "file",
		"related_id":   fileID,
	}

	msgService := messageService.GetMessageService()
	if err := msgService.SendTemplateMessage(userID, common.MessageTypeFileDeletedByAdmin, variables); err != nil {
		logger.Warn("发送管理员删除文件消息失败: userID=%d, fileID=%s, error=%v", userID, fileID, err)
	} else {
	}
}
