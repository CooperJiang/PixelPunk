package folder

import (
	"pixelpunk/internal/controllers/folder/dto"
	"pixelpunk/internal/middleware"
	folderService "pixelpunk/internal/services/folder"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

func MoveFolders(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	var req dto.MoveFoldersRequest
	if validatedReq, err := common.ValidateRequest[dto.MoveFoldersRequest](c); err != nil {
		errors.HandleError(c, err)
		return
	} else {
		req = *validatedReq
	}

	if err := folderService.MoveFolders(userID, req.FolderIDs, req.NewParentID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, nil, "文件夹移动成功")
}

/* ValidateMove 接口已删除 - 移动时直接验证即可，无需单独接口 */
