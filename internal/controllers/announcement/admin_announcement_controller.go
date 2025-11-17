package announcement

import (
	"pixelpunk/internal/controllers/announcement/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/services/announcement"
	"pixelpunk/pkg/errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateAnnouncementHandler(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)
	if userID == 0 {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	var createDTO dto.AnnouncementCreateDTO
	if err := c.ShouldBindJSON(&createDTO); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := announcement.CreateAnnouncement(userID, &createDTO)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "创建公告成功")
}

func UpdateAnnouncementHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的公告ID"))
		return
	}

	var updateDTO dto.AnnouncementUpdateDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := announcement.UpdateAnnouncement(uint(id), &updateDTO)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "更新公告成功")
}

func DeleteAnnouncementHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的公告ID"))
		return
	}

	if err := announcement.DeleteAnnouncement(uint(id)); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "删除公告成功")
}

func GetAnnouncementHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的公告ID"))
		return
	}

	result, err := announcement.GetAnnouncementByID(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取公告成功")
}

func GetAnnouncementListHandler(c *gin.Context) {
	var query dto.AnnouncementQueryDTO
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := announcement.GetAnnouncementList(&query)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "获取公告列表成功")
}

func GetAnnouncementSettingsHandler(c *gin.Context) {
	settings, err := announcement.GetAnnouncementSettings()
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, settings, "获取配置成功")
}

func UpdateAnnouncementSettingsHandler(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	if err := announcement.UpdateAnnouncementSettings(settings); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "配置更新成功")
}

func TogglePinAnnouncementHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "无效的公告ID"))
		return
	}

	result, err := announcement.TogglePinAnnouncement(uint(id))
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, result, "操作成功")
}
