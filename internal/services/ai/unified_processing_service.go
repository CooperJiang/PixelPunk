package ai

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"pixelpunk/internal/models"
	tagService "pixelpunk/internal/services/tag"
	ai "pixelpunk/pkg/ai"
	"pixelpunk/pkg/ai/prompts"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"

	"gorm.io/gorm"
)

// buildTaggingContext 构建标签识别的上下文信息
func buildTaggingContext(categoryResult *ai.FileCategorizationResponse) (string, string, *uint) {
	var categoryName, categoryDescription string
	var categoryID *uint

	if categoryResult != nil && categoryResult.Success {
		if categoryResult.CategoryID > 0 {
			categoryID = &categoryResult.CategoryID
			categoryName = categoryResult.CategoryName
			categoryDescription = categoryResult.CategoryDescription
		} else {
			// 建议创建新分类，使用通用标签识别
			categoryID = nil
			categoryName = ""
			categoryDescription = ""
		}
	} else {
		categoryName = ""
		categoryDescription = ""
	}

	return categoryName, categoryDescription, categoryID
}

// performAITagging 执行AI标签识别
func performAITagging(file models.File, base64Data, imageFormat, categoryName, categoryDescription string, categoryID *uint) (*AIFileResponse, error) {
	imageTagService := tagService.NewFileGlobalTagService()
	availableTags, err := imageTagService.BuildTagsForAI(file.UserID, categoryID)
	if err != nil {
		availableTags = []ai.TagInfo{} // 使用空列表继续处理
	}

	// 将可用标签限制为前N个（按使用次数降序），降低提示词长度与成本
	const maxPromptTags = 80
	if len(availableTags) > maxPromptTags {
		sort.Slice(availableTags, func(i, j int) bool { return availableTags[i].UsageCount > availableTags[j].UsageCount })
		availableTags = availableTags[:maxPromptTags]
	}

	// 转换标签类型以避免循环导入
	promptTags := make([]prompts.TagInfo, len(availableTags))
	for i, tag := range availableTags {
		promptTags[i] = prompts.TagInfo{
			ID:          tag.ID,
			Name:        tag.Name,
			Description: tag.Description,
			Source:      tag.Source,
			UsageCount:  tag.UsageCount,
		}
	}

	// 使用统一的提示词管理获取增强提示词
	enhancedPrompt := prompts.GetEnhancedImageAnalysisPrompt(categoryName, categoryDescription, promptTags)

	aiResp, err := ai.AnalyzeImageByBase64(base64Data, imageFormat, enhancedPrompt)
	if err != nil {
		return nil, err
	}

	return convertAIResponse(aiResp), nil
}

// convertAIResponse 转换AI响应格式
func convertAIResponse(aiResp *ai.AIResponse) *AIFileResponse {
	convertedResp := &AIFileResponse{
		Success:      aiResp.Success,
		Data:         aiResp.Data,
		ErrMsg:       aiResp.ErrMsg,
		RawResponse:  aiResp.Data,
		HttpDuration: aiResp.HttpDuration,
	}

	// 转换Usage类型（如果存在）
	if aiResp.Usage != nil {
		convertedResp.Usage = &TokenUsage{
			PromptTokens:     aiResp.Usage.PromptTokens,
			CompletionTokens: aiResp.Usage.CompletionTokens,
			TotalTokens:      aiResp.Usage.TotalTokens,
		}
	}

	return convertedResp
}

// saveCategoryResultIfNeeded 如果需要则保存分类结果
func saveCategoryResultIfNeeded(tx *gorm.DB, fileID string, categoryResult *ai.FileCategorizationResponse) error {
	if categoryResult != nil && categoryResult.Success {
		var file models.File
		err := tx.Where("id = ?", fileID).First(&file).Error
		if err != nil {
			return fmt.Errorf("获取图片信息失败: %v", err)
		}

		var finalCategoryID uint

		if categoryResult.CategoryID > 0 {
			finalCategoryID = categoryResult.CategoryID
		} else {
			// AI建议创建新分类，先检查是否已存在同名分类

			// 首先检查是否已存在同名的用户分类
			var existingCategory models.FileCategory
			err = tx.Where("user_id = ? AND name = ? AND status = ?", file.UserID, categoryResult.CategoryName, "active").
				First(&existingCategory).Error

			if err == nil {
				// 找到现有分类，直接使用
				finalCategoryID = existingCategory.ID

				// 更新现有分类的使用次数（忽略错误）
				_ = tx.Model(&models.FileCategory{}).
					Where("id = ?", existingCategory.ID).
					UpdateColumn("file_count", gorm.Expr("file_count + 1")).Error
			} else if err == gorm.ErrRecordNotFound {
				// 不存在同名分类，创建新的用户分类
				newCategory := &models.FileCategory{
					Name:        categoryResult.CategoryName,
					Description: categoryResult.CategoryDescription,
					UserID:      file.UserID,
					Source:      "ai_suggestion",
					Status:      "active",
				}

				err := tx.Create(newCategory).Error
				if err != nil {
					return fmt.Errorf("创建AI建议分类失败: %v", err)
				}

				finalCategoryID = newCategory.ID
			} else {
				// 查询出错（非记录不存在错误）
				return fmt.Errorf("查询现有分类失败: %v", err)
			}

			// 更新AI返回结果中的CategoryID，方便后续使用
			categoryResult.CategoryID = finalCategoryID
		}

		// 直接更新图片的分类ID，不再创建关联表记录（避免数据库锁冲突）
		err = tx.Model(&models.File{}).
			Where("id = ?", fileID).
			Updates(map[string]interface{}{
				"category_id":     finalCategoryID,
				"category_source": "ai",
			}).Error
		if err != nil {
			return fmt.Errorf("保存分类结果到图片失败: %v", err)
		}

		// 异步更新分类使用次数（避免在主事务中造成锁冲突）
		if finalCategoryID > 0 {
			go func(categoryID uint, userID uint) {
				// 使用新的数据库连接，避免事务冲突
				_ = updateCategoryUsageCountAsync(categoryID, userID)
			}(finalCategoryID, file.UserID)
		}
	}
	return nil
}

// performAIImageCategorizationOutsideTx 在事务外执行AI分类识别（用于优化SQLite并发）
func performAIImageCategorizationOutsideTx(file models.File, base64Data string, imageFormat string) (*ai.FileCategorizationResponse, error) {
	maxRetries := common.AICategorizationMaxRetries

	for attempt := 1; attempt <= maxRetries; attempt++ {
		result, err := performAIImageCategorizationSync(file, base64Data, imageFormat)
		if err == nil && result != nil && result.Success {
			return result, nil
		}

		if attempt < maxRetries && err != nil {
			// 指数退避+抖动
			backoff := time.Duration(200*attempt) * time.Millisecond
			jitter := time.Duration((time.Now().UnixNano()%100)+50) * time.Millisecond
			sleep := backoff + jitter
			time.Sleep(sleep)
		}
	}

	return nil, fmt.Errorf("AI分类识别重试 %d 次均失败", maxRetries)
}

// performAIImageCategorizationSync 同步执行AI图片分类识别
func performAIImageCategorizationSync(file models.File, base64Data string, imageFormat string) (*ai.FileCategorizationResponse, error) {

	db := GetDBFromContext()
	if db == nil {
		return nil, fmt.Errorf("无法获取数据库连接")
	}

	var categories []models.CategoryTemplate
	var userCategories []models.FileCategory

	// 获取所有系统分类模板（不限制数量，按热门和排序优先）
	err := db.Where("1 = 1"). // 获取全部系统分类
					Order("is_popular DESC, sort_order ASC, usage_count DESC"). // 热门优先，再按排序和使用次数
					Find(&categories).Error
	if err != nil {
		return nil, fmt.Errorf("获取系统分类模板失败: %v", err)
	}

	// 获取用户自定义分类（最多150个，忽略错误）
	_ = db.Where("user_id = ? AND status = ?", file.UserID, "active").
		Order("sort_order ASC").
		Limit(150).
		Find(&userCategories).Error

	var aiCategories []ai.CategoryInfo

	for _, cat := range categories {
		aiCategories = append(aiCategories, ai.CategoryInfo{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Source:      "system_template",
		})
	}

	for _, cat := range userCategories {
		aiCategories = append(aiCategories, ai.CategoryInfo{
			ID:          cat.ID,
			Name:        cat.Name,
			Description: cat.Description,
			Source:      "user",
		})
	}

	if len(aiCategories) == 0 {
		return nil, nil
	}

	resp, err := ai.CategorizeImageByBase64(base64Data, imageFormat, aiCategories)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// updateCategoryUsageCountAsync 异步更新分类使用次数（使用独立的数据库连接）
func updateCategoryUsageCountAsync(categoryID uint, userID uint) error {
	db := database.GetDB()

	// 使用重试机制处理锁等待超时
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := updateCategoryUsageCountInternal(db, categoryID, userID)
		if err == nil {
			return nil
		}

		if strings.Contains(err.Error(), "Lock wait timeout") {
			if attempt < maxRetries-1 {
				time.Sleep(time.Duration(100*(attempt+1)) * time.Millisecond) // 递增延迟
				continue
			}
		}

		return err
	}

	return nil
}

// updateCategoryUsageCountInternal 内部更新逻辑
func updateCategoryUsageCountInternal(db *gorm.DB, categoryID uint, userID uint) error {
	// 首先检查是否是用户自定义分类
	var userCategory models.FileCategory
	err := db.Where("id = ? AND user_id = ?", categoryID, userID).First(&userCategory).Error

	if err == nil {
		// 是用户分类，更新FileCategory的file_count
		err = db.Model(&models.FileCategory{}).
			Where("id = ?", categoryID).
			UpdateColumn("file_count", gorm.Expr("file_count + 1")).Error

		if err != nil {
			return fmt.Errorf("更新用户分类使用次数失败: %v", err)
		}

		return nil
	} else if err != gorm.ErrRecordNotFound {
		// 查询出错（非记录不存在）
		return fmt.Errorf("查询用户分类失败: %v", err)
	}

	// 不是用户分类，检查是否是系统分类模板
	var template models.CategoryTemplate
	err = db.Where("id = ?", categoryID).First(&template).Error

	if err == nil {
		// 是系统分类模板，更新CategoryTemplate的usage_count
		err = db.Model(&models.CategoryTemplate{}).
			Where("id = ?", categoryID).
			UpdateColumn("usage_count", gorm.Expr("usage_count + 1")).Error

		if err != nil {
			return fmt.Errorf("更新系统分类模板使用次数失败: %v", err)
		}

		return nil
	} else if err == gorm.ErrRecordNotFound {
		// 既不是用户分类也不是系统模板
		return nil
	}

	return fmt.Errorf("查询系统分类模板失败: %v", err)
}

// AiImageTaggingWithBase64AndPrompt 使用自定义提示词进行AI图像标记
func AiImageTaggingWithBase64AndPrompt(base64Data, imageFormat, customPrompt string) (*AIFileResponse, error) {
	// 使用新的统一AI客户端，传入自定义提示词
	aiResp, err := ai.AnalyzeImageByBase64(base64Data, imageFormat, customPrompt)
	if err != nil {
		return nil, fmt.Errorf("调用AI客户端失败: %v", err)
	}

	parsedResp := convertToOldResponse(aiResp)

	return parsedResp, nil
}
