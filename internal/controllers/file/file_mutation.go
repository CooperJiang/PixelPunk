package file

import (
	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func UpdateFile(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	req, err := common.ValidateRequest[dto.UpdateFileDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 先获取原文件信息以记录日志
	oldFileInfo, err := filesvc.GetFileDetail(currentUser.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	imgInfo, err := filesvc.UpdateFile(currentUser.UserID, fileID, req.Name, req.FolderID, req.AccessLevel)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 记录文件重命名活动日志
	if req.Name != "" && req.Name != oldFileInfo.DisplayName {
		oldName := oldFileInfo.DisplayName
		if oldName == "" {
			oldName = oldFileInfo.OriginalName
		}
		activity.LogFileRename(currentUser.UserID, oldName, req.Name, fileID)
	}

	if req.FolderID != "" && req.FolderID != "null" && req.FolderID != oldFileInfo.FolderID {
		fileName := imgInfo.DisplayName
		if fileName == "" {
			fileName = imgInfo.OriginalName
		}

		oldFolderName := "根目录"
		if oldFileInfo.FolderID != "" {
			var oldFolder models.Folder
			if err := database.DB.Where("id = ?", oldFileInfo.FolderID).First(&oldFolder).Error; err == nil {
				oldFolderName = oldFolder.Name
			}
		}

		newFolderName := "根目录"
		if req.FolderID != "" {
			var newFolder models.Folder
			if err := database.DB.Where("id = ?", req.FolderID).First(&newFolder).Error; err == nil {
				newFolderName = newFolder.Name
			}
		}

		activity.LogFileMove(currentUser.UserID, fileName, fileID, oldFolderName, newFolderName)
	}

	errors.ResponseSuccess(c, imgInfo, "更新成功")
}

func DeleteFile(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	// 先获取文件信息以记录日志
	fileInfo, err := filesvc.GetFileDetail(currentUser.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = filesvc.DeleteFile(currentUser.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	fileName := fileInfo.DisplayName
	if fileName == "" {
		fileName = fileInfo.OriginalName
	}
	activity.LogFileDelete(currentUser.UserID, fileName, fileID)

	errors.ResponseSuccess(c, gin.H{"id": fileID}, "删除成功")
}
func BatchDeleteFiles(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	var req struct {
		FileIDs []string `json:"file_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "请求参数错误: "+err.Error()))
		return
	}

	if len(req.FileIDs) == 0 {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID列表不能为空"))
		return
	}

	successIds, failIds, err := filesvc.BatchDeleteUserFiles(currentUser.UserID, req.FileIDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if len(successIds) > 0 {
		activity.LogBatchDelete(currentUser.UserID, len(successIds), "")
	}

	response := gin.H{
		"success_count": len(successIds),
		"fail_count":    len(failIds),
		"success_ids":   successIds,
		"fail_ids":      failIds,
	}

	errors.ResponseSuccess(c, response, "批量删除完成")
}
func ToggleAccessLevel(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	fileID := c.Param("file_id")
	if fileID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "文件ID不能为空"))
		return
	}

	// 先获取原文件信息以记录权限变化
	oldFileInfo, err := filesvc.GetFileDetail(currentUser.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	imgInfo, err := filesvc.ToggleFileAccessLevel(currentUser.UserID, fileID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 记录文件权限切换活动日志
	fileName := oldFileInfo.DisplayName
	if fileName == "" {
		fileName = oldFileInfo.OriginalName
	}
	activity.LogFileAccessLevelChange(currentUser.UserID, fileName, oldFileInfo.AccessLevel, imgInfo.AccessLevel)

	errors.ResponseSuccess(c, imgInfo, "访问级别切换成功")
}

// UploadForApiKey 通过API密钥上传文件（支持单张和多张，无需JWT认证）

func ReorderFiles(c *gin.Context) {
	currentUser := middleware.GetCurrentUser(c)

	req, err := common.ValidateRequest[dto.ReorderFilesDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	err = filesvc.ReorderFiles(currentUser.UserID, req.FolderID, req.FileIDs)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "文件排序成功")
}

// CheckDuplicate MD5预检查重复文件

// InstantUpload 秒传上传
