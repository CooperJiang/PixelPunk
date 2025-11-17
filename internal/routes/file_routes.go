package routes

import (
	fileController "pixelpunk/internal/controllers/file"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFileRoutes(r *gin.RouterGroup) {
	guestGroup := r.Group("/guest")
	guestGroup.GET("/list", fileController.GetRecommendedFileList)
	guestGroup.GET("/random", fileController.GetRandomRecommendedFile)

	guestGroup.POST("/upload", middleware.UploadConcurrencyLimit(), fileController.GuestUpload)

	guestGroup.POST("/check-duplicate", fileController.CheckDuplicate)
	guestGroup.POST("/instant-upload", fileController.InstantUpload)

	r.GET("/:file_id/download",
		middleware.JWTAuth(),
		middleware.OptionalAuthForFileDownload(),
		fileController.DownloadFile)

	authGroup := r.Group("")
	authGroup.Use(middleware.RequireAuth())

	authGroup.POST("/upload", middleware.UploadConcurrencyLimit(), fileController.Upload)
	authGroup.POST("/batch-upload", middleware.UploadConcurrencyLimit(), fileController.BatchUpload)

	authGroup.POST("/check-duplicate", fileController.CheckDuplicate)
	authGroup.POST("/instant-upload", fileController.InstantUpload)

	authGroup.GET("/list", fileController.GetFileList)

	authGroup.POST("/batch-delete", fileController.BatchDeleteFiles)

	authGroup.POST("/reorder", fileController.ReorderFiles)

	authGroup.POST("/move", fileController.MoveFiles)

	authGroup.GET("/:file_id/link", fileController.GenerateFileLink)
	authGroup.POST("/:file_id/toggle-access-level", fileController.ToggleAccessLevel)

	authGroup.GET("/:file_id", fileController.GetFileDetail)

	authGroup.PUT("/:file_id", fileController.UpdateFile)

	authGroup.DELETE("/:file_id", fileController.DeleteFile)
}
