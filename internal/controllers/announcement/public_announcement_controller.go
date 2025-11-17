package announcement

import (
	"pixelpunk/internal/services/announcement"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPublicAnnouncementListHandler(c *gin.Context) {
	// 查询公告列表（配置由后端控制）
	result, err := announcement.GetPublicAnnouncementList()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取公告列表成功")
}

func GetPublicAnnouncementDetailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的公告ID"))
		return
	}

	result, err := announcement.GetPublicAnnouncementDetail(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取公告详情成功")
}
