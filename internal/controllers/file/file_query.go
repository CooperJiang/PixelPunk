package file

import (
	"strings"
	"time"

	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"
	filesvc "pixelpunk/internal/services/file"
	setting "pixelpunk/internal/services/setting"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func GetFileList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.FileListQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	size := req.Size
	if size <= 0 {
		size = 20
	}
	sort := req.Sort
	if sort == "" {
		sort = "newest"
	}

	var tagsArray []string
	if req.Tags != "" {
		tagsArray = strings.Split(req.Tags, ",")
	}

	var dominantColorsArray []string
	if req.DominantColor != "" {
		dominantColorsArray = strings.Split(req.DominantColor, ",")
	}

	var categoryIDsArray []string
	if req.CategoryID != "" {
		categoryIDsArray = strings.Split(req.CategoryID, ",")
	}

	searchParams := filesvc.AdminFileSearchParams{
		Page:          page,
		Size:          size,
		Sort:          sort,
		Keyword:       req.Keyword,
		Tags:          tagsArray,
		CategoryIDs:   categoryIDsArray,
		DominantColor: dominantColorsArray,
		Resolution:    req.Resolution,
		MinWidth:      req.MinWidth,
		MaxWidth:      req.MaxWidth,
		MinHeight:     req.MinHeight,
		MaxHeight:     req.MaxHeight,
		UserID:        userID, // 设置为当前用户ID，限制只查询该用户的文件
	}

	if req.FolderID != "" {
		searchParams.FolderID = req.FolderID
	}

	if req.AccessLevel != "" {
		searchParams.AccessLevel = req.AccessLevel
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
			"size":         size,
			"current_page": page,
			"last_page":    (total + int64(size) - 1) / int64(size),
		},
	}

	errors.ResponseSuccess(c, data, "获取成功")
}

func GetFileDetail(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	imgInfo, err := filesvc.GetFileDetail(userID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, imgInfo, "获取成功")
}
func GetFileStats(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	req, err := common.ValidateRequest[dto.FileStatsQueryDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	stats, err := filesvc.GetFileStats(userID, fileID, req.Period)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, stats, "获取成功")
}
func GetRecommendedFileList(c *gin.Context) {

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

	// 如果没有指定排序方式，游客默认按浏览量排序
	sortOrder := params.Sort
	if sortOrder == "" {
		sortOrder = "views"
	}

	// 强制设置IsRecommended=true，只获取推荐的文件
	isRecommended := true

	searchParams := filesvc.AdminFileSearchParams{
		Page:          params.Page,
		Size:          params.Size,
		Sort:          sortOrder,
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
		UserID:        0, // 禁止访客模式下用户ID过滤
		AccessLevel:   "public",
		IsRecommended: &isRecommended,
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

	errors.ResponseSuccess(c, data, "获取推荐文件列表成功")
}

// UserGetTagList 获取当前用户文件的标签列表
func GetPublicFileCount(c *gin.Context) {
	showImageCount := setting.GetBool("website_info", "show_file_count", false)
	showStorageUsage := setting.GetBool("website_info", "show_storage_usage", false)

	// 如果两个配置都未启用，则不返回任何数据
	if !showImageCount && !showStorageUsage {
		errors.ResponseSuccess(c, gin.H{}, "统计功能已关闭")
		return
	}

	data := gin.H{
		"timestamp": time.Now().Unix(),
	}

	// 如果启用文件数量统计，获取文件总数
	if showImageCount {
		count, err := filesvc.GetTotalFileCount()
		if err != nil {
			errors.HandleError(c, err)
			return
		}
		data["total"] = count
	}

	// 如果启用存储使用统计，获取存储信息
	if showStorageUsage {
		totalStorage, formattedStorage, err := filesvc.GetTotalStorageUsage()
		if err != nil {
			errors.HandleError(c, err)
			return
		}
		data["storage"] = gin.H{
			"total_storage":     totalStorage,
			"formatted_storage": formattedStorage,
		}
	}

	errors.ResponseSuccess(c, data, "获取公开统计数据成功")
}

func GetRandomRecommendedFile(c *gin.Context) {
	imageInfo, err := filesvc.GetRandomRecommendedFile()
	if err != nil {
		// 特殊处理：当没有推荐文件时，返回特殊的成功响应而不是错误
		if err.(*errors.Error).Code == errors.CodeServiceUnavailable {
			errors.ResponseSuccess(c, nil, err.Error())
			return
		}
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, imageInfo, "获取随机推荐文件成功")
}
