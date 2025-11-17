package file

/* Admin queries for files (no behavior change). */

import (
	"encoding/json"
	"fmt"
	"pixelpunk/internal/models"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"strings"
)

/* AdminGetFileList 管理员获取文件列表（语义化命名） */
func AdminGetFileList(params AdminFileSearchParams) ([]AdminFileDetailResponse, int64, error) {
	var total int64
	var images []models.File
	var responses []AdminFileDetailResponse

	query := database.DB.Model(&models.File{}).Where("status <> ?", StatusPendingDeletion)

	if len(params.Tags) > 0 {
		var imageIDs []string
		if err := database.DB.Model(&models.FileGlobalTagRelation{}).Where("tag_id IN ?", params.Tags).Distinct("file_id").Pluck("file_id", &imageIDs).Error; err != nil {
			return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询标签关系失败")
		}
		if len(imageIDs) > 0 {
			query = query.Where("id IN ?", imageIDs)
		} else {
			return []AdminFileDetailResponse{}, 0, nil
		}
	}

	if len(params.CategoryIDs) > 0 {
		// 将字符串数组转换为 uint 数组
		categoryIDs := make([]uint, 0, len(params.CategoryIDs))
		for _, idStr := range params.CategoryIDs {
			var id uint64
			if _, err := fmt.Sscanf(idStr, "%d", &id); err == nil {
				categoryIDs = append(categoryIDs, uint(id))
			}
		}
		if len(categoryIDs) > 0 {
			query = query.Where("category_id IN ?", categoryIDs)
		} else {
			return []AdminFileDetailResponse{}, 0, nil
		}
	}

	if params.Keyword != "" {
		nameQuery := database.DB.Where("original_name LIKE ? OR display_name LIKE ?", "%"+params.Keyword+"%", "%"+params.Keyword+"%")
		var aiMatchingIDs []string
		database.DB.Model(&models.FileAIInfo{}).Where("description LIKE ?", "%"+params.Keyword+"%").Pluck("file_id", &aiMatchingIDs)
		var tagIDs []uint
		database.DB.Model(&models.GlobalTag{}).Where("name LIKE ?", "%"+params.Keyword+"%").Pluck("id", &tagIDs)
		var tagMatchingIDs []string
		if len(tagIDs) > 0 {
			database.DB.Model(&models.FileGlobalTagRelation{}).Where("tag_id IN ?", tagIDs).Pluck("file_id", &tagMatchingIDs)
		}
		var allMatchingIDs []string
		allMatchingIDs = append(allMatchingIDs, aiMatchingIDs...)
		allMatchingIDs = append(allMatchingIDs, tagMatchingIDs...)
		matchingIDMap := make(map[string]bool)
		for _, id := range allMatchingIDs {
			matchingIDMap[id] = true
		}
		uniqueMatchingIDs := make([]string, 0, len(matchingIDMap))
		for id := range matchingIDMap {
			uniqueMatchingIDs = append(uniqueMatchingIDs, id)
		}
		if len(uniqueMatchingIDs) > 0 {
			query = query.Where(database.DB.Where(nameQuery).Or("id IN ?", uniqueMatchingIDs))
		} else {
			query = query.Where(nameQuery)
		}
	}

	if params.UserID > 0 {
		query = query.Where("user_id = ?", params.UserID)
	}
	if params.FolderID != "" {
		query = query.Where("folder_id = ?", params.FolderID)
	}
	if params.AccessLevel != "" {
		query = query.Where("access_level = ?", params.AccessLevel)
	}
	if params.IsRecommended != nil {
		query = query.Where("is_recommended = ?", *params.IsRecommended)
	}
	if params.StorageType != "" {
		query = query.Where("storage_type = ?", params.StorageType)
	}
	if params.MinWidth > 0 {
		query = query.Where("width >= ?", params.MinWidth)
	}
	if params.MaxWidth > 0 {
		query = query.Where("width <= ?", params.MaxWidth)
	}
	if params.MinHeight > 0 {
		query = query.Where("height >= ?", params.MinHeight)
	}
	if params.MaxHeight > 0 {
		query = query.Where("height <= ?", params.MaxHeight)
	}

	if len(params.DominantColor) > 0 {
		var allColorFormats []string
		for _, color := range params.DominantColor {
			if strings.HasPrefix(color, "#") {
				allColorFormats = append(allColorFormats, color, strings.TrimPrefix(color, "#"))
			} else {
				allColorFormats = append(allColorFormats, color, "#"+color)
			}
		}
		var colorMatchFileIDs []string
		database.DB.Model(&models.FileAIInfo{}).Where("dominant_color IN ?", allColorFormats).Pluck("file_id", &colorMatchFileIDs)
		if len(colorMatchFileIDs) > 0 {
			query = query.Where("id IN ?", colorMatchFileIDs)
		} else {
			return []AdminFileDetailResponse{}, 0, nil
		}
	}

	var aiFiltered bool
	var aiFilterFileIDs []string
	aiQuery := database.DB.Model(&models.FileAIInfo{})
	if params.Resolution != "" {
		aiQuery = aiQuery.Where("resolution = ?", params.Resolution)
	}
	if params.NSFWMinScore > 0 {
		aiQuery = aiQuery.Where("nsfw_score >= ?", params.NSFWMinScore)
		aiFiltered = true
	}
	if params.NSFWMaxScore > 0 && params.NSFWMaxScore <= 1 {
		aiQuery = aiQuery.Where("nsfw_score <= ?", params.NSFWMaxScore)
		aiFiltered = true
	}
	if params.IsNSFW != nil {
		aiQuery = aiQuery.Where("is_nsfw = ?", *params.IsNSFW)
		aiFiltered = true
	}
	if aiFiltered {
		aiQuery.Pluck("file_id", &aiFilterFileIDs)
		if len(aiFilterFileIDs) > 0 {
			query = query.Where("id IN ?", aiFilterFileIDs)
		} else {
			return []AdminFileDetailResponse{}, 0, nil
		}
	}

	var countQuery = query
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "获取文件总数失败")
	}

	switch strings.ToLower(params.Sort) {
	case "newest":
		query = query.Order("created_at DESC")
	case "oldest":
		query = query.Order("created_at ASC")
	case "name":
		query = query.Order("display_name ASC")
	case "size":
		query = query.Order("size DESC")
	case "width":
		query = query.Order("width DESC")
	case "height":
		query = query.Order("height DESC")
	case "views":
		query = query.Joins("LEFT JOIN file_stats ON file_stats.file_id = file.id").Order("COALESCE(file_stats.views, 0) DESC")
	default:
		query = query.Order("created_at DESC")
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 20
	}
	offset := (params.Page - 1) * params.Size

	selectFields := []string{"id", "user_id", "folder_id", "original_name", "display_name",
		"url", "thumb_url", "size", "width", "height", "format", "access_level",
		"is_recommended", "storage_provider_id", "is_duplicate", "md5_hash",
		"created_at", "updated_at", "remote_url", "remote_thumb_url",
		"storage_duration", "expires_at"}
	if err := query.Select(selectFields).Offset(offset).Limit(params.Size).Find(&images).Error; err != nil {
		return nil, 0, errors.Wrap(err, errors.CodeDBQueryFailed, "查询文件列表失败")
	}

	imageIDs := make([]string, 0, len(images))
	for _, file := range images {
		imageIDs = append(imageIDs, file.ID)
	}
	var aiInfoList []models.FileAIInfo
	if len(imageIDs) > 0 {
		database.DB.Where("file_id IN ?", imageIDs).Find(&aiInfoList)
	}
	aiInfoMap := make(map[string]models.FileAIInfo)
	for _, ai := range aiInfoList {
		aiInfoMap[ai.FileID] = ai
	}
	var statsList []models.FileStats
	if len(imageIDs) > 0 {
		database.DB.Where("file_id IN ?", imageIDs).Find(&statsList)
	}
	statsMap := make(map[string]int64)
	for _, s := range statsList {
		statsMap[s.FileID] = s.Views
	}

	for _, file := range images {
		var userName string
		if file.UserID > 0 {
			var user models.User
			if err := database.DB.Select("username").Where("id = ?", file.UserID).First(&user).Error; err == nil {
				userName = user.Username
			}
		}
		var aiInfo *AIInfoResponse
		if ai, ok := aiInfoMap[file.ID]; ok {
			aiInfo = convertToAIResponse(ai)
		}
		views := statsMap[file.ID]
		resp := BuildAdminFileDetailResponse(file, views, userName, aiInfo)
		responses = append(responses, resp)
	}
	return responses, total, nil
}

/* AdminGetImageList -> AdminGetFileList */

func convertToAIResponse(ai models.FileAIInfo) *AIInfoResponse {
	var tagNames []string
	if len(ai.Tags) > 0 {
		if err := json.Unmarshal(ai.Tags, &tagNames); err != nil {
			logger.Error("解析文件标签JSON失败: %v", err)
		}
	}
	if len(tagNames) == 0 {
		var tags []models.GlobalTag
		database.DB.Model(&models.GlobalTag{}).Joins("JOIN file_global_tag_relation ON file_global_tag_relation.tag_id = global_tag.id").Where("file_global_tag_relation.file_id = ?", ai.FileID).Find(&tags)
		for _, tag := range tags {
			tagNames = append(tagNames, tag.Name)
		}
	}
	var colorPalette []string
	if len(ai.ColorPalette) > 0 {
		if err := json.Unmarshal(ai.ColorPalette, &colorPalette); err != nil {
			logger.Error("解析颜色调色板JSON失败: %v", err)
		}
	}
	var nsfwCategories map[string]float64
	if len(ai.NSFWCategories) > 0 {
		if err := json.Unmarshal(ai.NSFWCategories, &nsfwCategories); err != nil {
			logger.Error("解析NSFW分类JSON失败: %v", err)
		}
	}
	return &AIInfoResponse{
		Description: ai.Description, Tags: tagNames, DominantColor: ai.DominantColor, Resolution: ai.Resolution,
		IsNSFW: ai.IsNSFW, NSFWScore: ai.NSFWScore, NSFWEvaluation: ai.NSFWEvaluation,
		ColorPalette: colorPalette, AspectRatio: ai.AspectRatio, Composition: ai.Composition, ObjectsCount: ai.ObjectsCount,
		NSFWCategories: nsfwCategories,
	}
}
