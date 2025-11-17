package file

import (
	"net/http"

	"pixelpunk/internal/controllers/file/dto"
	"pixelpunk/internal/middleware"
	filesvc "pixelpunk/internal/services/file"
	"pixelpunk/pkg/common"
	"pixelpunk/pkg/errors"

	"github.com/gin-gonic/gin"
)

// InitChunkedUpload 初始化分片上传
func InitChunkedUpload(c *gin.Context) {
	userID := middleware.GetCurrentUserID(c)

	req, err := common.ValidateRequest[dto.InitChunkedUploadDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	session, err := filesvc.InitChunkedUpload(userID, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	response := &dto.ChunkedUploadResponse{
		SessionID:      session.SessionID,
		Status:         session.Status,
		Progress:       session.Progress,
		TotalChunks:    session.TotalChunks,
		UploadedChunks: 0,
		Message:        "分片上传会话创建成功",
	}

	errors.ResponseSuccess(c, response, "初始化分片上传成功")
}

// UploadChunk 上传分片
func UploadChunk(c *gin.Context) {
	req, err := common.ValidateRequest[dto.UploadChunkDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "获取上传文件失败: "+err.Error()))
		return
	}

	if err := filesvc.UploadChunk(req, file); err != nil {
		errors.HandleError(c, err)
		return
	}

	status, err := filesvc.GetChunkedUploadStatus(req.SessionID)
	if err != nil {
		// 分片上传成功，但获取状态失败，仍返回成功
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "分片上传成功",
			"data": gin.H{
				"session_id":    req.SessionID,
				"chunk_number":  req.ChunkNumber,
				"upload_status": "success",
			},
		})
		return
	}

	errors.ResponseSuccess(c, status, "分片上传成功")
}

// CompleteChunkedUpload 完成分片上传
func CompleteChunkedUpload(c *gin.Context) {
	req, err := common.ValidateRequest[dto.CompleteChunkedUploadDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	fileResponse, err := filesvc.CompleteChunkedUpload(req.SessionID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 直接返回文件信息，与正常上传保持一致的格式
	errors.ResponseSuccess(c, fileResponse, "分片上传完成")
}

func GetChunkedUploadStatus(c *gin.Context) {
	req, err := common.ValidateRequest[dto.ChunkedUploadStatusDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	status, err := filesvc.GetChunkedUploadStatus(req.SessionID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, status, "获取上传状态成功")
}

// CancelChunkedUpload 取消分片上传
func CancelChunkedUpload(c *gin.Context) {
	req, err := common.ValidateRequest[dto.CancelChunkedUploadDTO](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	if err := filesvc.CancelChunkedUpload(req.SessionID); err != nil {
		errors.HandleError(c, err)
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"session_id": req.SessionID,
		"status":     "cancelled",
	}, "取消分片上传成功")
}
