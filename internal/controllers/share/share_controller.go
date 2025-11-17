package share

import (
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/internal/services/share"
	"pixelpunk/pkg/cache"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"
	"pixelpunk/pkg/storage"
	"pixelpunk/pkg/utils"

	"fmt"
	"io"
	"net/http"
	"time"

	"crypto/md5"

	"github.com/gin-gonic/gin"
)

func CreateShare(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.CreateShareDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	result, err := share.CreateShare(userID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	shareType := "mixed"
	if len(req.Items) == 1 {
		shareType = req.Items[0].ItemType
	}
	activity.LogShareCreate(userID, result.ID, shareType)

	data := gin.H{
		"id":        result.ID,
		"share_key": result.ShareKey,
		"share_url": getShareURL(c, result.ShareKey),
	}

	errors.ResponseSuccess(c, data, "创建分享成功")
}

func GetShareList(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var query dto.ShareListQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误: "+err.Error()))
		return
	}

	shareList, total, err := share.GetUserShares(userID, &query)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"list":  shareList,
		"total": total,
	}

	errors.ResponseSuccess(c, data, "获取分享列表成功")
}

func GetShareDetail(c *gin.Context) {
	shareID := c.Param("id")

	userID := middleware.GetCurrentUserID(c)

	var shareObj models.Share
	if err := database.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&shareObj).Error; err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "分享不存在或您无权访问"))
		return
	}

	items, err := share.GetShareItems(shareID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"share": shareObj,
		"items": items,
		"url":   getShareURL(c, shareObj.ShareKey),
	}

	errors.ResponseSuccess(c, data, "获取分享详情成功")
}

func DeleteShare(c *gin.Context) {
	shareID := c.Param("id")

	userID := middleware.GetCurrentUserID(c)

	forceStr := c.Query("force")
	force := forceStr == "true" || forceStr == "1"

	// 先获取分享信息以记录日志
	var shareInfo models.Share
	if err := database.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&shareInfo).Error; err == nil {
		activity.LogShareDelete(userID, shareID, "分享")
	}

	if err := share.DeleteShare(shareID, userID, force); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除分享成功")
}

func ViewShare(c *gin.Context) {
	shareKey := c.Param("key")

	shareInfo, err := share.GetShareByKey(shareKey)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, err.Error()))
		return
	}

	if shareInfo.Password != "" {
		accessToken := c.Query("access_token")
		if accessToken != "" {
			valid, err := share.ValidateAccessToken(shareKey, accessToken)
			if err != nil {
				errors.HandleError(c, err)
				return
			}

			if !valid {
				errors.ResponseSuccess(c, gin.H{
					"require_password": true,
					"share_id":         shareInfo.ID,
					"name":             shareInfo.Name,
				}, "需要密码验证")
				return
			}

		} else {
			errors.ResponseSuccess(c, gin.H{
				"require_password": true,
				"share_id":         shareInfo.ID,
				"name":             shareInfo.Name,
			}, "需要密码验证")
			return
		}
	}

	folderID := c.Query("folder_id")

	data, err := share.GetShareForView(shareKey, folderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	referer := c.Request.Referer()

	viewedItems := []map[string]string{}

	userAgentHash := fmt.Sprintf("%x", md5.Sum([]byte(userAgent)))
	cacheKey := fmt.Sprintf("share_view_%s_%s_%s", shareInfo.ID, clientIP, userAgentHash[:8])

	isFirstView := !cache.Exists(cacheKey)

	if isFirstView {
		share.IncreaseShareViews(shareInfo.ID)

		err := cache.Set(cacheKey, "1", 24*time.Hour)
		if err != nil {
			logger.Error("设置访问缓存失败: %v", err)
		}
	} else {
	}

	share.LogShareAccess(shareInfo.ID, viewedItems, nil, clientIP, userAgent, referer)

	errors.ResponseSuccess(c, data, "获取分享内容成功")
}

func VerifySharePassword(c *gin.Context) {
	shareKey := c.Param("key")

	req, err := common.ValidateRequest[dto.VerifySharePasswordDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()

	accessToken, err := share.GenerateAccessToken(shareKey, req.Password, clientIP, userAgent)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"access_token": accessToken,
	}

	errors.ResponseSuccess(c, data, "密码验证成功")
}

func DownloadFilesBatch(c *gin.Context) {
	var req struct {
		FileIDs     []string `json:"file_ids" binding:"required,min=1"`
		AccessToken string   `json:"access_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误"))
		return
	}

	data := gin.H{
		"accepted":     true,
		"file_ids":     req.FileIDs,
		"download_url": "", // 未来扩展
		"message":      "批量下载任务已受理",
	}
	errors.ResponseSuccess(c, data, "已受理批量下载请求")
}

func SubmitVisitorInfo(c *gin.Context) {
	shareKey := c.Param("key")

	req, err := common.ValidateRequest[dto.VisitorInfoDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	shareInfo, err := share.GetShareByKey(shareKey)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, err.Error()))
		return
	}

	if !shareInfo.CollectVisitorInfo {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "该分享不需要收集访客信息"))
		return
	}

	clientIP := c.ClientIP()
	userAgent := c.Request.UserAgent()
	referer := c.Request.Referer()

	if err := share.SaveVisitorInfo(shareKey, req, clientIP, userAgent, referer); err != nil {
		errors.HandleError(c, err)
		return
	}

	visitorKey := "share_visitor_" + shareKey
	c.Set(visitorKey, req)

	errors.ResponseSuccess(c, nil, "提交访客信息成功")
}

func getShareURL(c *gin.Context, shareKey string) string {
	baseUrl := utils.GetBaseUrl()
	return baseUrl + "/share/" + shareKey
}

func GetShareVisitors(c *gin.Context) {
	shareID := c.Param("id")

	userID := middleware.GetCurrentUserID(c)

	var shareObj models.Share
	if err := database.DB.Where("id = ? AND user_id = ?", shareID, userID).First(&shareObj).Error; err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "分享不存在或您无权访问"))
		return
	}

	var query dto.VisitorQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误: "+err.Error()))
		return
	}

	visitors, total, err := share.GetShareVisitors(shareID, &query)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	data := gin.H{
		"list":  visitors,
		"total": total,
	}

	errors.ResponseSuccess(c, data, "获取访客信息列表成功")
}

func DeleteShareVisitor(c *gin.Context) {
	shareID := c.Param("id")
	visitorID := c.Param("visitor_id")

	userID := middleware.GetCurrentUserID(c)

	if err := share.DeleteVisitorInfo(shareID, visitorID, userID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除访客信息成功")
}

func DownloadSharedFile(c *gin.Context) {
	shareKey := c.Param("key")
	fileID := c.Param("file_id")

	if shareKey == "" || fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "分享密钥和文件ID不能为空"))
		return
	}

	accessToken := c.Query("access_token")

	shareInfo, err := share.GetShareByKey(shareKey)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "分享不存在或已失效"))
		return
	}

	if shareInfo.Password != "" {
		if accessToken == "" {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, "需要提供访问令牌"))
			return
		}

		valid, err := share.ValidateAccessToken(shareKey, accessToken)
		if err != nil || !valid {
			errors.HandleError(c, errors.New(errors.CodeUnauthorized, "访问令牌无效或已过期"))
			return
		}
	}

	hasAccess, err := share.ValidateSharedFileAccess(shareInfo.ID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if !hasAccess {
		errors.HandleError(c, errors.New(errors.CodeFileAccessDenied, "该文件不在分享内容中"))
		return
	}

	var file models.File
	if err := database.DB.Where("id = ?", fileID).
		Where("status <> ?", "pending_deletion").
		First(&file).Error; err != nil {
		errors.HandleError(c, errors.New(errors.CodeFileNotFound, "文件不存在"))
		return
	}

	result, isLocal, isProxy, err := filesvc.ServeFile(file, false)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	go func() {
		downloadLog := &models.FileDownloadLog{
			UserID:    0, // 分享下载设置为0，表示游客下载
			FileID:    fileID,
			FileSize:  file.Size,
			IPAddress: c.ClientIP(),
			UserAgent: c.GetHeader("User-Agent"),
			ShareKey:  shareKey, // 记录分享密钥
		}
		if err := database.DB.Create(downloadLog).Error; err != nil {
			logger.Error("记录分享下载日志失败: %v", err)
		}
	}()

	fileName := file.DisplayName
	if fileName == "" {
		fileName = file.OriginalName
	}

	fileName = utils.GetSafeFilename(fileName)

	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Disposition", utils.SetContentDispositionFilename(fileName))
	c.Header("Content-Length", fmt.Sprintf("%d", file.Size))
	c.Header("Accept-Ranges", "bytes")

	switch {
	case isLocal:
		c.File(result.(string))
	case isProxy:
		proxyResp := result.(*filesvc.ProxyResponse)
		defer proxyResp.Content.Close()
		c.Status(http.StatusOK)
		io.Copy(c.Writer, proxyResp.Content)
	default:
		storageService := storage.NewGlobalStorage()
		fileReader, err := storageService.ReadFile(c.Request.Context(), file.StorageProviderID, file.URL)
		if err != nil {
			errors.HandleError(c, errors.New(errors.CodeFileNotFound, "无法读取文件文件: "+err.Error()))
			return
		}
		defer fileReader.Close()
		c.Status(http.StatusOK)
		io.Copy(c.Writer, fileReader)
	}
}
