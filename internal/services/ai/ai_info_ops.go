package ai

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/logger"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func saveFileAIInfo(tx *gorm.DB, fileID string, result *AITaggingResult, usage *TokenUsage) (*models.FileAIInfo, error) {
	colorPaletteJSON, err := json.Marshal(result.VisualElements.ColorPalette)
	if err != nil {
		return nil, fmt.Errorf("序列化颜色调色板失败: %v", err)
	}
	tagsJSON, err := json.Marshal(result.Tags)
	if err != nil {
		return nil, fmt.Errorf("序列化标签失败: %v", err)
	}
	semanticKeywordsJSON, err := json.Marshal(result.SemanticKeywords)
	if err != nil {
		return nil, fmt.Errorf("序列化语义关键词失败: %v", err)
	}
	nsfwCategoriesJSON, err := json.Marshal(result.ContentSafety.Categories)
	if err != nil {
		return nil, fmt.Errorf("序列化NSFW分类失败: %v", err)
	}

	aiInfo := &models.FileAIInfo{
		FileID:           fileID,
		Description:      result.Description,
		SearchContent:    result.SearchContent,
		SemanticKeywords: semanticKeywordsJSON,
		Tags:             tagsJSON,
		Width:            result.BasicInfo.Width,
		Height:           result.BasicInfo.Height,
		AspectRatio:      result.BasicInfo.AspectRatio,
		Resolution:       result.BasicInfo.Resolution,
		FileType:         result.BasicInfo.ImageType,
		EstimatedSize:    result.BasicInfo.EstimatedSize,
		DominantColor:    result.VisualElements.DominantColor,
		ColorPalette:     colorPaletteJSON,
		ObjectsCount:     result.VisualElements.ObjectsCount,
		Composition:      result.VisualElements.Composition,
		IsNSFW:           result.ContentSafety.IsNSFW,
		NSFWScore:        result.ContentSafety.NSFWScore,
		NSFWCategories:   nsfwCategoriesJSON,
		NSFWEvaluation:   result.ContentSafety.EvaluationResult,
		NSFWReason:       result.ContentSafety.NSFWReason,
	}

	// 使用 UPSERT 操作，避免并发时的重复插入问题
	// ON CONFLICT (file_id) DO UPDATE 会在冲突时更新所有字段
	if err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "file_id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"description", "search_content", "semantic_keywords", "tags",
			"width", "height", "aspect_ratio", "resolution", "file_type", "estimated_size",
			"dominant_color", "color_palette", "objects_count", "composition",
			"is_nsfw", "nsfw_score", "nsfw_categories", "nsfw_evaluation", "nsfw_reason",
			"updated_at",
		}),
	}).Create(aiInfo).Error; err != nil {
		return nil, fmt.Errorf("保存AI信息失败: %v", err)
	}

	// 同步更新 file.description
	if result.Description != "" {
		if err := tx.Model(&models.File{}).Where("id = ?", fileID).Update("description", result.Description).Error; err != nil {
			logger.Error("更新文件描述失败: %v", err)
		}
	}
	return aiInfo, nil
}
