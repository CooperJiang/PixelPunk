package routes

import (
	fileController "pixelpunk/internal/controllers/file"
	randomAPIController "pixelpunk/internal/controllers/random_api"
	"pixelpunk/internal/middleware"
	"pixelpunk/pkg/health"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	r.Use(middleware.IpRefererMiddleware())

	RegisterClientRoutes(r)

	prefix := r.Group("/api")
	version := prefix.Group("/v1")

	RegisterSetupRoutes(version)

	version.Use(middleware.InstallCheckMiddleware())

	version.GET("/health", health.SimpleHealthHandler)
	version.GET("/health/basic", health.BasicHealthHandler)
	version.GET("/health/complete", health.CompleteHealthHandler)

	RegisterMetricsRoutes(version)

	pbRoutes := version.Group("/pb")
	{
		pbRoutes.GET("/stats/files/count", fileController.GetPublicFileCount)
	}

	// 注册公开的认证路由（不需要JWT认证）
	authRoutes := version.Group("/auth")
	RegisterAuthRoutes(authRoutes)

	// 注册公开的用户路由（兼容旧的API路径，不需要JWT认证）
	publicUserRoutes := version.Group("/user")
	RegisterPublicUserRoutes(publicUserRoutes)

	// 注册公开的公告路由（不需要JWT认证）
	RegisterPublicAnnouncementRoutes(version)

	// JWT 中间件必须在所有需要认证的路由之前注册
	version.Use(middleware.JWTAuth())
	version.Use(middleware.TrackUserActivity())

	commonRoutes := version.Group("/common")
	RegisterCommonRoutes(commonRoutes)

	authorRoutes := version.Group("/authors")
	RegisterAuthorRoutes(authorRoutes)

	userRoutes := version.Group("/user")
	RegisterUserRoutes(userRoutes)

	personalRoutes := version.Group("/personal")
	personalRoutes.Use(middleware.RequireAuth())
	RegisterPersonalRoutes(personalRoutes)

	fileRoutes := version.Group("/files")
	RegisterFileRoutes(fileRoutes)

	RegisterChunkedUploadRoutes(fileRoutes)

	RegisterConfigRoutes(version)

	folderRoutes := version.Group("/folders")
	RegisterFolderRoutes(folderRoutes)

	tagRoutes := version.Group("/tags")
	RegisterTagRoutes(tagRoutes)

	RegisterCategoryTemplateRoutes(version)

	RegisterUserCategoryRoutes(version)

	RegisterUserTagRoutes(version)

	RegisterAutomationRoutes(version)

	apiKeyRoutes := version.Group("/apikey")
	RegisterAPIKeyRoutes(apiKeyRoutes)

	randomAPIRoutes := version.Group("/random-api")
	RegisterRandomAPIRoutes(randomAPIRoutes)

	adminRoutes := version.Group("/admin")
	RegisterAdminRoutes(adminRoutes)

	RegisterWebSocketRoutes(adminRoutes)

	adminShareRoutes := version.Group("/admin/shares")
	RegisterAdminShareRoutes(adminShareRoutes)

	adminContentReviewRoutes := version.Group("/admin")
	RegisterAdminContentReviewRoutes(adminContentReviewRoutes)

	aiRoutes := version.Group("/admin/ai")
	RegisterAIRoutes(aiRoutes)

	storageRoutes := version.Group("/storage")
	RegisterStorageRoutes(storageRoutes)

	settingRoutes := version.Group("/settings")
	RegisterSettingRoutes(settingRoutes)

	shareRoutes := version.Group("/shares")
	RegisterShareRoutes(shareRoutes)

	RegisterSearchRoutes(version)

	vectorRoutes := version.Group("/admin")
	RegisterVectorRoutes(vectorRoutes)

	RegisterMessageRoutes(version)

	// 注册公告管理端路由（需要管理员权限）
	RegisterAdminAnnouncementRoutes(version)

	// 注意：静态文件已通过 embed.FS 嵌入，由 RegisterClientRoutes 的 NoRoute 处理
	// 不需要额外的 r.Static 路由

	{
		fileIDGroup := r.Group("/f")
		fileIDGroup.Use(middleware.FileInfoExtractorMiddleware())
		fileIDGroup.Use(middleware.FileAccessControlMiddleware())
		fileIDGroup.Use(middleware.BandwidthLimitMiddleware())
		fileIDGroup.Use(middleware.BandwidthTrackingMiddleware())
		fileIDGroup.GET("/:fileID", fileController.ServeFileByID)
		fileIDGroup.GET("/:fileID/*displayName", fileController.ServeFileByID)

		thumbGroup := r.Group("/t")
		thumbGroup.Use(middleware.FileInfoExtractorMiddleware())
		thumbGroup.Use(middleware.FileAccessControlMiddleware())
		thumbGroup.Use(middleware.BandwidthLimitMiddleware())
		thumbGroup.Use(middleware.BandwidthTrackingMiddleware())
		thumbGroup.GET("/:fileID", fileController.ServeThumbByID)
		thumbGroup.GET("/:fileID/*displayName", fileController.ServeThumbByID)

		shortLinkGroup := r.Group("/s")
		shortLinkGroup.Use(middleware.FileInfoExtractorMiddleware())
		shortLinkGroup.Use(middleware.FileAccessControlMiddleware())
		shortLinkGroup.Use(middleware.BandwidthLimitMiddleware())
		shortLinkGroup.Use(middleware.BandwidthTrackingMiddleware())
		shortLinkGroup.GET("/:shortURL", fileController.ServeFileByShortURL)
	}

	r.POST("/api/v1/avatar/upload", middleware.RequireAuth(), fileController.UploadAvatar)
	r.GET("/file/avatar/:fileName", fileController.ServeAvatarFile)

	r.GET("/file/admin/:fileName", fileController.ServeAdminFile)

	apiUploadRoutes := r.Group("/api/v1/external")
	apiUploadRoutes.Use(middleware.APIKeyAuthMiddleware())
	apiUploadRoutes.POST("/upload", fileController.UploadForApiKey)

	// 随机图片API公开接口（不需要认证）
	randomImageRoutes := r.Group("/api/v1/r")
	randomImageRoutes.GET("/:api_key", randomAPIController.GetRandomImage)
}
