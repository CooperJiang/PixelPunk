package routes

import (
	shareController "pixelpunk/internal/controllers/share"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterShareRoutes(r *gin.RouterGroup) {
	r.POST("/download-files", shareController.DownloadFilesBatch)
	userShareGroup := r.Group("")
	userShareGroup.Use(middleware.RequireAuth())

	userShareGroup.POST("", shareController.CreateShare)

	userShareGroup.GET("", shareController.GetShareList)

	userShareGroup.GET("/:id", shareController.GetShareDetail)

	userShareGroup.GET("/:id/visitors", shareController.GetShareVisitors)

	userShareGroup.DELETE("/:id/visitors/:visitor_id", shareController.DeleteShareVisitor)

	userShareGroup.DELETE("/:id", shareController.DeleteShare)

	publicGroup := r.Group("/public")

	publicGroup.GET("/:key", shareController.ViewShare)

	publicGroup.POST("/:key/verify", shareController.VerifySharePassword)

	publicGroup.POST("/:key/visitor", shareController.SubmitVisitorInfo)

	publicGroup.GET("/:key/files/:file_id/download", shareController.DownloadSharedFile)
}
