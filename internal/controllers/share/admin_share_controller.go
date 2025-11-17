package share

import (
	"pixelpunk/internal/controllers/share/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/share"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/utils"

	"github.com/gin-gonic/gin"
)

func AdminGetShareList(c *gin.Context) {
	var query dto.AdminShareListQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误: "+err.Error()))
		return
	}

	shareList, total, err := share.GetAllShares(&query)
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

func AdminGetShareDetail(c *gin.Context) {
	shareID := c.Param("id")

	shareInfo, items, err := share.AdminGetShareDetail(shareID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	baseURL := "http://" + c.Request.Host
	shareURL := baseURL + "/share/" + shareInfo.ShareKey

	data := gin.H{
		"share": shareInfo,
		"items": items,
		"url":   shareURL,
	}

	if shareInfo.Password != "" {
		data["has_password"] = true
		data["access_token_hint"] = "该分享需要密码访问，可以通过 GET /api/v1/admin/shares/" + shareID + "/access-token 获取访问令牌"
	}

	errors.ResponseSuccess(c, data, "获取分享详情成功")
}

func AdminUpdateShareStatus(c *gin.Context) {
	shareID := c.Param("id")

	adminID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.AdminShareStatusUpdateDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := share.AdminUpdateShareStatus(shareID, req, adminID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "更新分享状态成功")
}

func AdminDeleteShare(c *gin.Context) {
	shareID := c.Param("id")

	adminID := middleware.GetCurrentUserID(c)

	forceStr := c.Query("force")
	force := forceStr == "true" || forceStr == "1"

	if err := share.AdminDeleteShare(shareID, adminID, force); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除分享成功")
}

func AdminGetShareStats(c *gin.Context) {
	stats, err := share.GetShareSystemStats()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, stats, "获取分享统计数据成功")
}

func AdminGetAllVisitors(c *gin.Context) {
	var query dto.VisitorQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误: "+err.Error()))
		return
	}

	visitors, total, err := share.AdminGetAllVisitors(&query)
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

func AdminDeleteVisitor(c *gin.Context) {
	visitorID := c.Param("id")

	if err := share.AdminDeleteVisitorInfo(visitorID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除访客信息成功")
}

func AdminGenerateAccessToken(c *gin.Context) {
	shareID := c.Param("id")

	adminID := middleware.GetCurrentUserID(c)

	accessToken, err := share.GenerateAdminAccessToken(shareID, adminID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	var shareInfo models.Share
	if err := database.DB.Select("share_key").Where("id = ?", shareID).First(&shareInfo).Error; err != nil {
		errors.HandleError(c, err)
		return
	}

	baseUrl := utils.GetBaseUrl()
	shareURL := baseUrl + "/share/" + shareInfo.ShareKey
	fullURL := shareURL + "?access_token=" + accessToken

	data := dto.AdminGenerateAccessTokenResponse{
		AccessToken: accessToken,
		ShareKey:    shareInfo.ShareKey,
		ShareURL:    shareURL,
		FullURL:     fullURL,
	}

	errors.ResponseSuccess(c, data, "生成访问令牌成功")
}
