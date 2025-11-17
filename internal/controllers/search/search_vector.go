package search

import (
	"fmt"
	"pixelpunk/internal/controllers/search/dto"
	"pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/vector"
	"time"

	"pixelpunk/internal/models"

	"github.com/gin-gonic/gin"
)

func getVectorConfig() (threshold float32, maxResults int, err error) {
	thresholdVal := setting.GetFloatDirectFromDB("vector", "vector_similarity_threshold", 0.7)
	threshold = float32(thresholdVal)

	maxResults = setting.GetIntDirectFromDB("vector", "vector_max_results", 50)

	return threshold, maxResults, nil
}

func getSearchThreshold() float32 {
	thresholdVal := setting.GetFloatDirectFromDB("vector", "vector_search_threshold", 0.3)
	return float32(thresholdVal)
}

func VectorSearch(c *gin.Context) {
	startTime := time.Now()

	req, err := common.ValidateRequest[dto.VectorSearchRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.Threshold == 0 {
		req.Threshold = 0.7
	}

	engine := vector.GetGlobalVectorEngine()
	if engine == nil {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索引擎未初始化"))
		return
	}

	if !engine.IsEnabled() {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索功能未启用"))
		return
	}

	searchResults, err := engine.SearchFiles(req.Query, req.Limit, 0, req.Threshold)
	if err != nil {
		logger.Error("向量搜索失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, fmt.Sprintf("搜索失败: %v", err)))
		return
	}

	results := make([]dto.VectorSearchResult, 0, len(searchResults))
	db := database.GetDB()

	for _, result := range searchResults {
		if result.Similarity < req.Threshold {
			continue
		}

		var file models.File
		if err := db.Where("id = ?", result.FileID).
			Where("status <> ?", "pending_deletion").
			First(&file).Error; err != nil {
			results = append(results, dto.VectorSearchResult{
				FileID:      result.FileID,
				Similarity:  result.Similarity,
				Description: result.Description,
				FileInfo:    nil,
			})
			continue
		}

		fileInfo := dto.ConvertFileToInfo(&file)
		results = append(results, dto.VectorSearchResult{
			FileID:      result.FileID,
			Similarity:  result.Similarity,
			Description: result.Description,
			FileInfo:    fileInfo,
		})
	}

	processTime := time.Since(startTime)

	response := dto.VectorSearchResponse{
		Query:       req.Query,
		Total:       len(results),
		Results:     results,
		ProcessTime: processTime.String(),
		UsedCache:   false,
	}

	errors.ResponseSuccess(c, response, "搜索完成")
}

func GetVectorStats(c *gin.Context) {
	userIDUint := uint(0)

	engine := vector.GetGlobalVectorEngine()
	if engine == nil {
		errors.HandleError(c, errors.New(errors.CodeServiceUnavailable, "向量搜索引擎未初始化"))
		return
	}

	db := database.GetDB()

	var stats dto.VectorStatsResponse

	db.Model(&models.FileVector{}).Count(&stats.TotalVectors)

	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusCompleted).Count(&stats.CompletedVectors)
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusPending).Count(&stats.PendingVectors)
	db.Model(&models.FileVector{}).Where("status = ?", common.VectorStatusFailed).Count(&stats.FailedVectors)

	db.Table("file_vector iv").
		Joins("JOIN file i ON iv.file_id = i.id").
		Where("i.user_id = ?", userIDUint).
		Count(&stats.UserVectors)

	var avgDimResult struct {
		AvgDimension float64
	}
	db.Model(&models.FileVector{}).
		Select("AVG(dimension) as avg_dimension").
		Where("status = ?", common.VectorStatusCompleted).
		Scan(&avgDimResult)
	stats.AvgDimension = int(avgDimResult.AvgDimension)

	errors.ResponseSuccess(c, stats, "统计信息获取成功")
}

func GetVectorHealth(c *gin.Context) {

	engine := vector.GetGlobalVectorEngine()

	response := dto.VectorHealthResponse{
		Status:      "error",
		Enabled:     false,
		Provider:    "",
		Model:       "",
		DatabaseOK:  false,
		EmbeddingOK: false,
		Message:     "向量引擎未初始化",
	}

	if engine == nil {
		errors.ResponseSuccess(c, response, "健康检查完成")
		return
	}

	response.Enabled = engine.IsEnabled()

	if !response.Enabled {
		response.Message = "向量搜索功能未启用"
		errors.ResponseSuccess(c, response, "健康检查完成")
		return
	}

	err := engine.HealthCheck()
	if err != nil {
		response.Message = fmt.Sprintf("健康检查失败: %v", err)
		logger.Error("向量引擎健康检查失败: %v", err)
	} else {
		response.Status = "ok"
		response.DatabaseOK = true
		response.EmbeddingOK = true
		response.Provider = "openai"
		response.Model = "text-embedding-3-small"
		response.Message = "向量引擎运行正常"
	}

	errors.ResponseSuccess(c, response, "健康检查完成")
}
