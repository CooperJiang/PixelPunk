package routes

import (
	announcementController "pixelpunk/internal/controllers/announcement"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

/* RegisterPublicAnnouncementRoutes 注册公告公开路由（无需认证） */
func RegisterPublicAnnouncementRoutes(r *gin.RouterGroup) {
	// 公开路由 - 获取公告列表和详情
	public := r.Group("/announcements")
	{
		public.GET("", announcementController.GetPublicAnnouncementListHandler)
		public.GET("/:id", announcementController.GetPublicAnnouncementDetailHandler)
	}
}

/* RegisterAdminAnnouncementRoutes 注册公告管理端路由（需要管理员权限） */
func RegisterAdminAnnouncementRoutes(r *gin.RouterGroup) {
	// 管理端路由 - 需要管理员权限
	admin := r.Group("/admin/announcements")
	admin.Use(middleware.RequireAdmin())
	{
		admin.POST("", announcementController.CreateAnnouncementHandler)
		admin.PUT("/:id", announcementController.UpdateAnnouncementHandler)
		admin.DELETE("/:id", announcementController.DeleteAnnouncementHandler)
		admin.GET("/:id", announcementController.GetAnnouncementHandler)
		admin.GET("", announcementController.GetAnnouncementListHandler)

		admin.PUT("/:id/toggle-pin", announcementController.TogglePinAnnouncementHandler)

		admin.GET("/settings", announcementController.GetAnnouncementSettingsHandler)
		admin.PUT("/settings", announcementController.UpdateAnnouncementSettingsHandler)
	}
}
