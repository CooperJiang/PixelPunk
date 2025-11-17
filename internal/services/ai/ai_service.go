package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/setting"
	tagService "pixelpunk/internal/services/tag"
	"pixelpunk/pkg/ai"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/logger"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AITaggingResult 是AI标记的结构化结果
type AITaggingResult struct {
	BasicInfo struct {
		AspectRatio   float64 `json:"aspect_ratio"`
		EstimatedSize string  `json:"estimated_size"`
		Height        int     `json:"height"`
		ImageType     string  `json:"image_type"`
		Resolution    string  `json:"resolution"`
		Width         int     `json:"width"`
	} `json:"basic_info"`
	ContentSafety struct {
		Categories struct {
			AlcoholTobacco float64 `json:"alcohol_tobacco"`
			Gambling       float64 `json:"gambling"`
			HateSpeech     float64 `json:"hate_speech"`
			Nudity         float64 `json:"nudity"`
			Violence       float64 `json:"violence"`
		} `json:"categories"`
		EvaluationResult string  `json:"evaluation_result"`
		IsNSFW           bool    `json:"is_nsfw"`
		NSFWScore        float64 `json:"nsfw_score"`
		NSFWReason       string  `json:"nsfw_reason"` // NSFW判断的原因说明
	} `json:"content_safety"`
	Description      string   `json:"description"`
	SearchContent    string   `json:"search_content"`    // 专为语义搜索优化的详细描述
	SemanticKeywords []string `json:"semantic_keywords"` // 语义关键词数组
	IsRecommended    bool     `json:"is_recommended"`
	Tags             []string `json:"tags"`
	VisualElements   struct {
		ColorPalette  []string `json:"color_palette"`
		Composition   string   `json:"composition"`
		DominantColor string   `json:"dominant_color"`
		ObjectsCount  int      `json:"objects_count"`
	} `json:"visual_elements"`
}

func AiImageTaggingAndSaveWithBase64(file models.File, base64Data, imageFormat string) error {
	db := GetDBFromContext()
	if db == nil {
		logger.Error("无法获取数据库连接")
		return fmt.Errorf("无法获取数据库连接")
	}

	contentDetectionEnabled := setting.GetBool("upload", "content_detection_enabled", true)
	sensitiveContentHandling := setting.GetString("upload", "sensitive_content_handling", "mark_only")
	aiAnalysisEnabled := setting.GetBool("upload", "ai_analysis_enabled", true)

	if !aiAnalysisEnabled {
		return db.Model(&models.File{}).Where("id = ?", file.ID).Update("ai_tagging_status", common.AITaggingStatusSkipped).Error
	}

	categoryResult, err := performAIImageCategorizationOutsideTx(file, base64Data, imageFormat)
	if err != nil {
		logger.Warn("AI分类失败，跳过AI打标: %v", err)
		db.Model(&models.File{}).Where("id = ?", file.ID).Update("ai_tagging_status", common.AITaggingStatusSkipped)
		return nil
	}

	categoryName, categoryDescription, categoryID := buildTaggingContext(categoryResult)

	aiResponse, err := performAITagging(file, base64Data, imageFormat, categoryName, categoryDescription, categoryID)
	if err != nil {
		logger.Warn("AI标签识别失败，跳过AI打标: %v", err)
		db.Model(&models.File{}).Where("id = ?", file.ID).Update("ai_tagging_status", common.AITaggingStatusSkipped)
		return nil
	}

	var fileCheck models.File
	if err := db.Where("id = ?", file.ID).Select("id").Take(&fileCheck).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errFileDeleted
		}
		return err
	}

	if err := saveCategoryResultIfNeeded(db, file.ID, categoryResult); err != nil {
		logger.Warn("保存分类结果失败，继续执行: %v", err)
	}

	if err := processAIResponse(db, file, aiResponse, contentDetectionEnabled, sensitiveContentHandling, base64Data, imageFormat); err != nil {
		logger.Warn("保存标签结果失败，继续执行: %v", err)
	}

	if err := db.Model(&models.File{}).Where("id = ?", file.ID).Update("ai_tagging_status", common.AITaggingStatusDone).Error; err != nil {
		logger.Error("更新打标完成状态失败: %v", err)
		return err
	}

	return nil
}

var errFileDeleted = errors.New("ai:file_deleted")

func fileExists(tx *gorm.DB, fileID string) bool {
	var n int64
	if err := tx.Model(&models.File{}).Where("id = ?", fileID).Count(&n).Error; err != nil {
		return false
	}
	return n > 0
}

func isDeadlockError(err error) bool {
	if err == nil {
		return false
	}
	s := strings.ToLower(err.Error())
	return strings.Contains(s, "deadlock found") || strings.Contains(s, "error 1213")
}

// processAIResponse 处理AI响应并保存结果
func processAIResponse(tx *gorm.DB, file models.File, aiResp *AIFileResponse, contentDetectionEnabled bool, sensitiveContentHandling string, base64Data string, imageFormat string) error {
	if aiResp == nil || !aiResp.Success {
		errMsg := "AI分析返回无效结果"
		if aiResp != nil && aiResp.ErrMsg != "" {
			errMsg = aiResp.ErrMsg
		}

		if contentDetectionEnabled && aiResp != nil && isNSFWRejection(aiResp) {

			// 根据敏感内容处理方式处理
			if sensitiveContentHandling == "auto_delete" {
				// 违规文件根据配置被标记为删除
				if err := markFileForDeletion(tx, file.ID); err != nil {
				}
			} else {
				if err := updateFileNSFWStatus(tx, file.ID, true); err != nil {
				}

				// 如果是等待审核，更新状态
				if sensitiveContentHandling == "pending_review" {
					nsfwReason := "AI检测到违规内容，详细原因不可用"
					if aiResp.ErrMsg != "" {
						nsfwReason = aiResp.ErrMsg
					}
					if err := markFileForReview(tx, file.ID, nsfwReason); err != nil {
						logger.Error("标记文件为待审核失败: %v", err)
					}
				}
			}

			updateFileStatus(tx, file.ID, common.AITaggingStatusDone)
			return fmt.Errorf("文件 %s 疑似违规内容，AI拒绝处理", file.ID)
		}

		logger.Error(errMsg)
		updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
		return errors.New(errMsg)
	}

	result, err := parseAITaggingResult(aiResp.Data)
	if err != nil {
		logger.Error("解析AI返回数据失败: %v", err)
		updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
		return err
	}

	// 如果原始文件已有分辨率信息，优先使用它而不是AI识别的分辨率
	if file.Resolution != "" {
		result.BasicInfo.Resolution = file.Resolution
	}

	// 使用 UPSERT 保存AI信息，自动处理新建或更新
	_, err = saveFileAIInfo(tx, file.ID, result, aiResp.Usage)
	if err != nil {
		if isDeadlockError(err) && !fileExists(tx, file.ID) {
			return errFileDeleted
		}
		updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
		logger.Error("保存AI标记结果失败: %v", err)
		return err
	}

	// 处理标签 - 根据配置决定是否为敏感内容生成标签
	if !contentDetectionEnabled || !result.ContentSafety.IsNSFW {
		// 如果没有启用内容检测或不是违规内容，正常处理标签
		if len(result.Tags) == 0 {
		}
		if err := processAndSaveTags(tx, file, result.Tags); err != nil {
			if isDeadlockError(err) && !fileExists(tx, file.ID) {
				return errFileDeleted
			}
			updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
			logger.Error("保存AI标记结果失败: %v", err)
			return err
		}
	} else {
		// 根据敏感内容处理方式处理
		if sensitiveContentHandling == "auto_delete" {
			// 违规文件根据配置被标记为删除
			if err := markFileForDeletion(tx, file.ID); err != nil {
			}
		} else if sensitiveContentHandling == "pending_review" {
			if err := markFileForReview(tx, file.ID, result.ContentSafety.NSFWReason); err != nil {
				logger.Error("标记文件为待审核失败: %v", err)
			}
			// 仍然保存标签，以便管理员审核时查看
			if err := processAndSaveTags(tx, file, result.Tags); err != nil {
				if isDeadlockError(err) && !fileExists(tx, file.ID) {
					return errFileDeleted
				}
				updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
				logger.Error("保存AI标记结果失败: %v", err)
				return err
			}
		} else {
			// 仅标记模式（默认），仍然保存标签
			if err := processAndSaveTags(tx, file, result.Tags); err != nil {
				if isDeadlockError(err) && !fileExists(tx, file.ID) {
					return errFileDeleted
				}
				updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
				logger.Error("保存AI标记结果失败: %v", err)
				return err
			}
		}
	}

	// 更新文件 NSFW 状态（如果启用内容检测且AI检测出不适内容）
	if contentDetectionEnabled && result.ContentSafety.IsNSFW {
		if err := updateFileNSFWStatus(tx, file.ID, true); err != nil {
			if isDeadlockError(err) && !fileExists(tx, file.ID) {
				return errFileDeleted
			}
			updateFileStatus(tx, file.ID, common.AITaggingStatusFailed)
			logger.Error("保存AI标记结果失败: %v", err)
			return err
		}
	}

	updateFileStatus(tx, file.ID, common.AITaggingStatusDone)

	if result != nil && result.Description != "" {
		if err := createPendingVectorRecord(file.ID, result.Description); err != nil {
			logger.Error("创建向量记录失败: %v", err)
			// 不返回错误，因为AI识别本身是成功的
		}
	}

	// 记录AI用量到打标日志（便于后续成本观测）
	if aiResp.Usage != nil {
		logData := map[string]interface{}{
			"prompt_tokens":     aiResp.Usage.PromptTokens,
			"completion_tokens": aiResp.Usage.CompletionTokens,
			"total_tokens":      aiResp.Usage.TotalTokens,
		}
		data, _ := json.Marshal(logData)
		logEntry := models.FileTaggingLog{
			FileID: file.ID,
			Status: common.AITaggingStatusDone,
			Action: common.TaggingActionAuto,
			Type:   "tagging.ai_completed",
			Data:   string(data),
		}
		if err := tx.Create(&logEntry).Error; err != nil {
			logger.Warn("写入AI用量日志失败: %v", err)
		}
	}

	go propagateAIToDuplicates(file.ID)

	return nil
}

// propagateAIToDuplicates 将原图的AI信息/标签/分类传播给其重复文件
func propagateAIToDuplicates(originalID string) {
	db := database.GetDB()
	if db == nil {
		return
	}

	var orig models.File
	if err := db.Where("id = ?", originalID).First(&orig).Error; err != nil {
		return
	}
	var origAI models.FileAIInfo
	_ = db.Where("file_id = ?", originalID).First(&origAI).Error

	var dups []models.File
	if err := db.Where("original_file_id = ?", originalID).Where("status <> ?", "pending_deletion").Find(&dups).Error; err != nil {
		return
	}
	if len(dups) == 0 {
		return
	}

	for _, dup := range dups {
		// 复制AI信息到重复文件（使用UPSERT避免重复插入）
		if origAI.FileID != "" {
			clone := origAI
			clone.ID = 0
			clone.FileID = dup.ID
			// 使用 UPSERT 操作，避免并发时的重复插入问题
			_ = db.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "file_id"}},
				UpdateAll: true,
			}).Create(&clone).Error
			if clone.Description != "" {
				_ = db.Model(&models.File{}).Where("id = ?", dup.ID).Update("description", clone.Description).Error
			}
		}

		// 标签：若重复图无任何标签关联则复制
		var count int64
		if err := db.Model(&models.FileGlobalTagRelation{}).Where("file_id = ?", dup.ID).Count(&count).Error; err == nil && count == 0 {
			var rels []models.FileGlobalTagRelation
			if err := db.Where("file_id = ?", originalID).Find(&rels).Error; err == nil {
				for _, r := range rels {
					nr := models.FileGlobalTagRelation{FileID: dup.ID, TagID: r.TagID, UserID: r.UserID, AccessLevel: r.AccessLevel, Source: r.Source, Confidence: r.Confidence}
					_ = db.Create(&nr).Error
				}
			}
		}

		updates := map[string]interface{}{"ai_tagging_status": common.AITaggingStatusDone}
		if orig.CategoryID != nil {
			updates["category_id"] = *orig.CategoryID
			updates["category_source"] = orig.CategorySource
		}
		_ = db.Model(&models.File{}).Where("id = ?", dup.ID).Updates(updates).Error
	}
}

func GetReviewQueueStats() (map[string]interface{}, error) {
	var pending int64
	var approved int64
	var rejected int64

	database.DB.Model(&models.File{}).Where("status = ?", "pending_review").Count(&pending)

	today := time.Now().Format("2006-01-02")
	database.DB.Model(&models.File{}).Where("DATE(created_at) = ? AND status = ?", today, "active").Count(&approved)

	return map[string]interface{}{
		"pending_count":  pending,
		"approved_today": approved,
		"rejected_today": rejected,
	}, nil
}

// isNSFWRejection 检测AI回复是否表示拒绝处理（可能是由于文件内容违规）
func isNSFWRejection(aiResp *AIFileResponse) bool {
	if aiResp.RawResponse == "" {
		return false
	}

	// 尝试解析为OpenAI标准回复格式
	var openaiResp struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	err := json.Unmarshal([]byte(aiResp.RawResponse), &openaiResp)
	if err != nil {
		return false
	}

	for _, choice := range openaiResp.Choices {
		content := choice.Message.Content
		if strings.Contains(content, "抱歉") ||
			strings.Contains(content, "无法处理") ||
			strings.Contains(content, "无法协助") ||
			strings.Contains(content, "不适合") ||
			strings.Contains(content, "不适当") ||
			strings.Contains(content, "违反") {
			return true
		}
	}

	return false
}

func GetDBFromContext() *gorm.DB {
	return database.GetDB()
}

// processAndSaveTags 处理并保存标签
func processAndSaveTags(tx *gorm.DB, file models.File, tags []string) error {
	// 使用新的全局标签架构处理
	if len(tags) == 0 {
		return nil
	}

	globalTagService := tagService.NewGlobalTagService()
	imageTagService := tagService.NewFileGlobalTagService()

	globalTags, err := globalTagService.CreateTagsFromNames(tags, file.UserID, "ai")
	if err != nil {
		logger.Error("创建AI标签失败: %v", err)
		return err
	}

	var tagIDs []uint
	for _, tag := range globalTags {
		tagIDs = append(tagIDs, tag.ID)
	}

	err = imageTagService.ReplaceFileTagsTx(tx, file.ID, tagIDs, "ai", common.AIDefaultConfidence)
	if err != nil {
		logger.Error("保存文件标签失败: %v", err)
		return err
	}

	return nil
}

// updateFileNSFWStatus 更新文件的 NSFW 状态
func updateFileNSFWStatus(tx *gorm.DB, fileID string, nsfw bool) error {
	return tx.Model(&models.File{}).
		Where("id = ?", fileID).
		Update("nsfw", nsfw).Error
}

// 辅助函数：更新文件状态
func updateFileStatus(tx *gorm.DB, fileID string, status string) error {
	updates := map[string]interface{}{
		"ai_tagging_status": status,
	}
	if status == common.AITaggingStatusFailed {
		updates["ai_tagging_tries"] = gorm.Expr("ai_tagging_tries + 1")
	}

	return tx.Model(&models.File{}).
		Where("id = ?", fileID).
		Updates(updates).Error
}

func TestAIConfigurationWithParams(params map[string]interface{}) (map[string]interface{}, error) {
	return ai.TestAIConfigurationWithParams(params)
}

// TestAIConfiguration 测试AI配置是否正确
func TestAIConfiguration() (map[string]interface{}, error) {
	return ai.TestAIConfiguration()
}

// CompressedFileData 压缩后的文件数据结构（避免循环引用）
type CompressedFileData struct {
	ThumbnailBase64 string // 缩略图的base64数据
}

// AiImageTaggingAndSaveWithCompressedFile 使用压缩后的缩略图base64进行AI分析并保存结果
func AiImageTaggingAndSaveWithCompressedFile(file models.File, imageData *CompressedFileData, imageFormat string) error {
	if imageData == nil {
		logger.Error("fileData为空")
		return fmt.Errorf("imageData为空")
	}

	if imageData.ThumbnailBase64 == "" {
		logger.Error("缺少缩略图base64数据")
		return fmt.Errorf("缺少缩略图base64数据")
	}

	return AiImageTaggingAndSaveWithBase64(file, imageData.ThumbnailBase64, imageFormat)
}
