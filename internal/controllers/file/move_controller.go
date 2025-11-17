package file

// 文件移动操作控制器

import (
	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/models"
	"pixelpunk/internal/services/activity"
	fileService "pixelpunk/internal/services/file"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/database"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

// MoveFiles 批量移动文件
func MoveFiles(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req dto.MoveFilesRequest
	if validatedReq, err := common.ValidateRequest[dto.MoveFilesRequest](c); err != nil {
		errors.HandleError(c, err)
		return
	} else {
		req = *validatedReq
	}

	// 获取文件信息以记录日志（批量移动前）
	var files []models.File
	if err := database.DB.Where("id IN ? AND user_id = ?", req.FileIDs, userID).Find(&files).Error; err == nil && len(files) > 0 {
		var targetFolderName = "根目录"
		if req.TargetFolderID != "" {
			var targetFolder models.Folder
			if err := database.DB.Where("id = ?", req.TargetFolderID).First(&targetFolder).Error; err == nil {
				targetFolderName = targetFolder.Name
			}
		}

		if err := fileService.MoveFiles(userID, req.FileIDs, req.TargetFolderID); err != nil {
			errors.HandleError(c, err)
			return
		}

		// 记录每个文件的移动日志
		for _, file := range files {
			oldFolderName := "根目录"
			if file.FolderID != "" {
				var oldFolder models.Folder
				if err := database.DB.Where("id = ?", file.FolderID).First(&oldFolder).Error; err == nil {
					oldFolderName = oldFolder.Name
				}
			}

			fileName := file.DisplayName
			if fileName == "" {
				fileName = file.OriginalName
			}

			activity.LogFileMove(userID, fileName, file.ID, oldFolderName, targetFolderName)
		}
	} else {
		// 如果获取文件信息失败，仍然尝试移动
		if err := fileService.MoveFiles(userID, req.FileIDs, req.TargetFolderID); err != nil {
			errors.HandleError(c, err)
			return
		}
	}

	errors.ResponseSuccess(c, nil, "文件移动成功")
}

// 注意：ReorderFiles 函数已存在于 file_controller.go 中，这里不重复定义以避免函数名冲突

// GetMoveInfo 接口已删除 - 无实际价值，移动时直接报错即可
