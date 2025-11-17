package routes

import (
	shareController "pixelpunk/internal/controllers/share"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminShareRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAdmin())

	r.GET("/list", shareController.AdminGetShareList)

	r.GET("/:id", shareController.AdminGetShareDetail)

	r.GET("/:id/access-token", shareController.AdminGenerateAccessToken)

	r.PUT("/:id/status", shareController.AdminUpdateShareStatus)

	r.DELETE("/:id", shareController.AdminDeleteShare)

	r.GET("/stats", shareController.AdminGetShareStats)

	visitorGroup := r.Group("/visitors")
	{
		visitorGroup.GET("", shareController.AdminGetAllVisitors)

		visitorGroup.DELETE("/:id", shareController.AdminDeleteVisitor)
	}
}
