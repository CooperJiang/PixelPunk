package admin

import (
	"encoding/json"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/ai"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/internal/services/review"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReviewQueueQueryDTO struct {
	Page     int    `form:"page,default=1" binding:"min=1"`
	Size     int    `form:"size,default=20" binding:"min=1,max=100"`
	Sort     string `form:"sort,default=newest"`
	Keyword  string `form:"keyword"`
	NSFWOnly bool   `form:"nsfw_only"`
}

type ReviewActionDTO struct {
	FileID     string `json:"file_id" binding:"required"`
	Action     string `json:"action" binding:"required,oneof=approve reject"`
	Reason     string `json:"reason"`
	HardDelete bool   `json:"hard_delete"` // 是否硬删除（仅reject时有效）
}

type BatchReviewActionDTO struct {
	FileIDs    []string `json:"file_ids" binding:"required,min=1"`
	Action     string   `json:"action" binding:"required,oneof=approve reject"`
	Reason     string   `json:"reason"`
	HardDelete bool     `json:"hard_delete"` // 是否硬删除（仅reject时有效）
}

type ReviewLogQueryDTO struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Size      int    `form:"size,default=20" binding:"min=1,max=100"`
	Action    string `form:"action"`     // approve, reject, 空为全部
	AuditorID uint   `form:"auditor_id"` // 审核员ID，0为全部
	Keyword   string `form:"keyword"`    // 搜索关键词
	DateFrom  string `form:"date_from"`  // 开始日期 YYYY-MM-DD
	DateTo    string `form:"date_to"`    // 结束日期 YYYY-MM-DD
}

func GetReviewQueue(c *gin.Context) {
	page := 1
	size := 20
	keyword := ""

	if pageStr := c.Query("params[page]"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	} else if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil {
			page = p
		}
	}

	if sizeStr := c.Query("params[size]"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			size = s
		}
	} else if sizeStr := c.Query("size"); sizeStr != "" {
		if s, err := strconv.Atoi(sizeStr); err == nil {
			size = s
		}
	}

	if k := c.Query("params[keyword]"); k != "" {
		keyword = k
	} else if k := c.Query("keyword"); k != "" {
		keyword = k
	}

	var reviewFiles []models.File
	db := database.GetDB().Model(&models.File{}).Preload("AIInfo").Preload("User").Where("status = ?", "pending_review")

	if keyword != "" {
		db = db.Where("original_name LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	var total int64
	countErr := db.Count(&total).Error
	if countErr != nil {
		errors.HandleError(c, errors.Wrap(countErr, errors.CodeDBQueryFailed, "统计审核队列总数失败"))
		return
	}

	offset := (page - 1) * size

	if err := db.Offset(offset).Limit(size).Order("created_at DESC").Find(&reviewFiles).Error; err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "查询审核队列失败"))
		return
	}

	signer := utils.GetURLSigner()

	var responseFiles []map[string]interface{}
	for _, file := range reviewFiles {
		var fullURL, fullThumbURL string
		if file.ID != "" {
			fullURL = signer.SignFileURL(file.ID, utils.SIGNATURE_DURATION)
			fullThumbURL = signer.SignThumbURL(file.ID, utils.SIGNATURE_DURATION)
		}

		var aiInfo map[string]interface{}
		if file.AIInfo != nil {
			var tags []string
			var semanticKeywords []string
			var colorPalette []string
			var nsfwCategories map[string]interface{}

			if file.AIInfo.Tags != nil {
				json.Unmarshal(file.AIInfo.Tags, &tags)
			}
			if file.AIInfo.SemanticKeywords != nil {
				json.Unmarshal(file.AIInfo.SemanticKeywords, &semanticKeywords)
			}
			if file.AIInfo.ColorPalette != nil {
				json.Unmarshal(file.AIInfo.ColorPalette, &colorPalette)
			}
			if file.AIInfo.NSFWCategories != nil {
				json.Unmarshal(file.AIInfo.NSFWCategories, &nsfwCategories)
			}

			aiInfo = map[string]interface{}{
				"description":     file.AIInfo.Description,
				"tags":            tags,
				"dominant_color":  file.AIInfo.DominantColor,
				"resolution":      file.AIInfo.Resolution,
				"is_nsfw":         file.AIInfo.IsNSFW,
				"nsfw_score":      file.AIInfo.NSFWScore,
				"nsfw_evaluation": file.AIInfo.NSFWEvaluation,
				"color_palette":   colorPalette,
				"aspect_ratio":    file.AIInfo.AspectRatio,
				"composition":     file.AIInfo.Composition,
				"objects_count":   file.AIInfo.ObjectsCount,
				"nsfw_categories": nsfwCategories,
			}
		}

		var uploaderInfo map[string]interface{}
		if file.User != nil {
			uploaderInfo = map[string]interface{}{
				"id":       file.User.ID,
				"username": file.User.Username,
				"email":    file.User.Email,
			}
		}

		responseFiles = append(responseFiles, map[string]interface{}{
			"id":             file.ID,
			"original_name":  file.OriginalName,
			"display_name":   file.DisplayName,
			"url":            file.URL,      // 保持原始文件名
			"thumb_url":      file.ThumbURL, // 保持原始缩略图文件名
			"full_url":       fullURL,       // 新增：完整签名URL
			"full_thumb_url": fullThumbURL,  // 新增：完整缩略图签名URL
			"size_formatted": file.SizeFormatted,
			"width":          file.Width,
			"height":         file.Height,
			"format":         file.Format,
			"nsfw":           file.NSFW,
			"nsfw_score": func() interface{} {
				if file.AIInfo != nil {
					return file.AIInfo.NSFWScore
				}
				return nil
			}(),
			"created_at": file.CreatedAt,
			"user_id":    file.UserID,
			"uploader":   uploaderInfo, // 新增：上传者信息
			"ai_info":    aiInfo,       // AI信息
		})
	}

	pageSize := size
	totalPages := (int(total) + pageSize - 1) / pageSize

	result := map[string]interface{}{
		"data": responseFiles,
		"pagination": map[string]interface{}{
			"page":       page,
			"page_size":  pageSize,
			"total":      int(total),
			"total_page": totalPages,
		},
	}

	errors.ResponseSuccess(c, result, "获取审核队列成功")
}

func GetReviewStats(c *gin.Context) {
	stats, err := ai.GetReviewQueueStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, stats, "获取审核统计成功")
}

func ReviewFile(c *gin.Context) {
	req, err := common.ValidateRequest[ReviewActionDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	auditorID := middleware.GetCurrentUserID(c)
	if auditorID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未找到当前用户信息"))
		return
	}

	var reviewErr error
	switch req.Action {
	case "approve":
		reviewErr = review.ApproveFileWithLog(req.FileID, auditorID, req.Reason)
	case "reject":
		reviewErr = review.RejectFileWithLog(req.FileID, auditorID, req.Reason, req.HardDelete)
	default:
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的审核操作"))
		return
	}

	if reviewErr != nil {
		errors.HandleError(c, errors.Wrap(reviewErr, errors.CodeInternal, "审核操作失败"))
		return
	}

	deleteType := "软删除"
	if req.HardDelete {
		deleteType = "硬删除"
	}
	message := "审核操作成功"
	if req.Action == "reject" {
		message = "审核拒绝成功 (" + deleteType + ")"
	}

	errors.ResponseSuccess(c, nil, message)
}

func BatchReview(c *gin.Context) {
	req, err := common.ValidateRequest[BatchReviewActionDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(req.FileIDs) > 100 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "批量操作最多支持100个文件"))
		return
	}

	auditorID := middleware.GetCurrentUserID(c)
	if auditorID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未找到当前用户信息"))
		return
	}

	results, err := review.BatchReviewFilesWithLog(req.FileIDs, req.Action, auditorID, req.Reason, req.HardDelete)
	if err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeInternal, "批量审核操作失败"))
		return
	}

	successCount := 0
	failCount := 0
	for _, result := range results {
		if result == "success" {
			successCount++
		} else {
			failCount++
		}
	}

	deleteType := "软删除"
	if req.HardDelete {
		deleteType = "硬删除"
	}
	message := "批量审核操作完成"
	if req.Action == "reject" {
		message = "批量审核操作完成 (" + deleteType + ")"
	}

	response := map[string]interface{}{
		"success_count": successCount,
		"fail_count":    failCount,
		"results":       results,
		"delete_type":   deleteType,
	}

	errors.ResponseSuccess(c, response, message)
}

func GetFileDetail(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	var file models.File
	if err := database.GetDB().Where("id = ?", fileID).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件不存在"))
		} else {
			errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件失败"))
		}
		return
	}

	if !file.NSFW {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "该文件不在审核队列中"))
		return
	}

	detail, err := filesvc.GetFileDetail(file.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, detail, "获取文件详情成功")
}

func GetReviewLogs(c *gin.Context) {
	req, err := common.ValidateRequest[ReviewLogQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	db := database.GetDB()

	query := db.Model(&models.ReviewLog{}).
		Preload("File").
		Preload("Auditor").
		Preload("Uploader")

	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}

	if req.AuditorID > 0 {
		query = query.Where("auditor_id = ?", req.AuditorID)
	}

	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where(
			"reason LIKE ? OR EXISTS (SELECT 1 FROM user WHERE user.id = review_log.auditor_id AND user.username LIKE ?) OR EXISTS (SELECT 1 FROM user WHERE user.id = review_log.uploader_id AND user.username LIKE ?) OR EXISTS (SELECT 1 FROM file WHERE file.id = review_log.file_id AND (file.original_name LIKE ? OR file.display_name LIKE ?))",
			keyword, keyword, keyword, keyword, keyword,
		)
	}

	if req.DateFrom != "" {
		query = query.Where("DATE(created_at) >= ?", req.DateFrom)
	}
	if req.DateTo != "" {
		query = query.Where("DATE(created_at) <= ?", req.DateTo)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "统计审核记录总数失败"))
		return
	}

	offset := (req.Page - 1) * req.Size
	var reviewLogs []models.ReviewLog
	if err := query.Offset(offset).Limit(req.Size).Order("created_at DESC").Find(&reviewLogs).Error; err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeDBQueryFailed, "查询审核记录失败"))
		return
	}

	var responseData []map[string]interface{}
	for _, log := range reviewLogs {
		logData := map[string]interface{}{
			"id":             log.ID,
			"file_id":        log.FileID,
			"action":         log.Action,
			"delete_type":    log.DeleteType,
			"reason":         log.Reason,
			"nsfw_score":     log.NSFWScore,
			"nsfw_threshold": log.NSFWThreshold,
			"is_nsfw":        log.IsNSFW,
			"created_at":     log.CreatedAt,
			"updated_at":     log.UpdatedAt,
		}

		if log.Auditor != nil {
			logData["auditor"] = map[string]interface{}{
				"id":       log.Auditor.ID,
				"username": log.Auditor.Username,
				"email":    log.Auditor.Email,
			}
		}

		if log.Uploader != nil {
			logData["uploader"] = map[string]interface{}{
				"id":       log.Uploader.ID,
				"username": log.Uploader.Username,
				"email":    log.Uploader.Email,
			}
		}

		if log.File != nil {
			signer := utils.GetURLSigner()
			var fullURL, fullThumbURL string
			if log.File.ID != "" {
				fullURL = signer.SignFileURL(log.File.ID, utils.SIGNATURE_DURATION)
				fullThumbURL = signer.SignThumbURL(log.File.ID, utils.SIGNATURE_DURATION)
			}

			logData["file"] = map[string]interface{}{
				"id":             log.File.ID,
				"original_name":  log.File.OriginalName,
				"display_name":   log.File.DisplayName,
				"url":            log.File.URL,
				"thumb_url":      log.File.ThumbURL,
				"full_url":       fullURL,
				"full_thumb_url": fullThumbURL,
				"size_formatted": log.File.SizeFormatted,
				"width":          log.File.Width,
				"height":         log.File.Height,
				"format":         log.File.Format,
				"status":         log.File.Status,
				"created_at":     log.File.CreatedAt,
			}
		} else {
			logData["file"] = map[string]interface{}{
				"id":            log.FileID,
				"original_name": "文件已删除",
				"status":        "deleted",
			}
		}

		responseData = append(responseData, logData)
	}

	totalPages := (int(total) + req.Size - 1) / req.Size

	result := map[string]interface{}{
		"data": responseData,
		"pagination": map[string]interface{}{
			"page":       req.Page,
			"page_size":  req.Size,
			"total":      int(total),
			"total_page": totalPages,
		},
	}

	errors.ResponseSuccess(c, result, "获取审核记录成功")
}

func HardDeleteReviewedFile(c *gin.Context) {
	fileID := c.Param("fileId")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	operatorID := middleware.GetCurrentUserID(c)
	if operatorID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未找到当前用户信息"))
		return
	}

	if err := review.HardDeleteSoftDeletedFile(fileID, operatorID); err != nil {
		errors.HandleError(c, errors.Wrap(err, errors.CodeInternal, "硬删除失败"))
		return
	}
	errors.ResponseSuccess(c, nil, "文件已彻底删除")
}
