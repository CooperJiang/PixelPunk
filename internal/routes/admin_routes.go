package routes

import (
	adminController "pixelpunk/internal/controllers/admin"
	aiController "pixelpunk/internal/controllers/ai"
	fileController "pixelpunk/internal/controllers/file"
	statsController "pixelpunk/internal/controllers/stats"
	userController "pixelpunk/internal/controllers/user"

	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAdminRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAuth())
	{
		r.GET("/user", statsController.UserStats)
	}

	statsAdmin := r.Group("/stats")
	statsAdmin.Use(middleware.RequireAdmin())
	{
		statsAdmin.GET("/latest-files", statsController.LatestFiles)

		statsAdmin.GET("/users", statsController.DashboardUserStats)
		statsAdmin.GET("/files", statsController.DashboardFileStats)
		statsAdmin.GET("/storage", statsController.DashboardStorageStats)
		statsAdmin.GET("/system-resources", statsController.DashboardSystemResources)

		statsAdmin.GET("/upload-trends", statsController.DashboardUploadTrends)
		statsAdmin.GET("/ai-services", statsController.DashboardAIServices)

		statsAdmin.GET("/shares", statsController.DashboardShareStats)
		statsAdmin.GET("/tags", statsController.DashboardTagStats)
		statsAdmin.GET("/system-info", statsController.DashboardSystemInfo)
	}

	userRoutes := r.Group("/user")
	userRoutes.Use(middleware.RequireAdmin())
	{
		userRoutes.GET("/list", userController.AdminGetUserList)
		userRoutes.POST("/create", middleware.RequireSuperAdmin(), userController.AdminCreateUser)
		userRoutes.GET("/detail/:id", userController.AdminGetUserDetail)
		userRoutes.POST("/update", middleware.RequireSuperAdmin(), userController.AdminUpdateUser)
		userRoutes.POST("/storage", middleware.RequireSuperAdmin(), userController.AdminUpdateUserStorage)
		userRoutes.POST("/reset-password/:id", middleware.RequireSuperAdmin(), userController.AdminResetUserPassword)
		userRoutes.POST("/send-email", middleware.RequireSuperAdmin(), userController.AdminSendUserEmail)
		userRoutes.POST("/toggle-status", middleware.RequireSuperAdmin(), userController.AdminToggleUserStatus)
		userRoutes.POST("/delete/:id", middleware.RequireSuperAdmin(), userController.AdminDeleteUser)
		userRoutes.POST("/batch", middleware.RequireSuperAdmin(), userController.AdminBatchOperateUsers)
	}

	imageRoutes := r.Group("/files")
	imageRoutes.Use(middleware.RequireAdmin())
	{
		imageRoutes.GET("/list", fileController.AdminGetFileList)
		imageRoutes.GET("/tags", fileController.AdminGetTagList)
		imageRoutes.GET("/colors", fileController.AdminGetColorList)
		imageRoutes.POST("/recommend", fileController.AdminRecommendFile)
		imageRoutes.POST("/batch-recommend", fileController.AdminBatchRecommendFiles)
		imageRoutes.POST("/delete", fileController.AdminDeleteFile)
		imageRoutes.POST("/batch-delete", fileController.AdminBatchDeleteFiles)
	}

	aiRoutes := r.Group("/ai")
	aiRoutes.Use(middleware.RequireAdmin())
	{
		aiRoutes.POST("/trigger-tagging", aiController.TriggerFileTagging)

		aiRoutes.POST("/reset-stuck", aiController.ResetStuckFiles)

		aiRoutes.POST("/test-config", aiController.TestAIConfig)
	}

	vectorVerificationRoutes := r.Group("/vector-verification")
	vectorVerificationRoutes.Use(middleware.RequireAdmin())
	{
		controller := adminController.NewVectorVerificationController()

		vectorVerificationRoutes.POST("/start", controller.StartVerificationTask)

		vectorVerificationRoutes.GET("/task/:taskId", controller.GetTaskStatus)

		vectorVerificationRoutes.GET("/tasks", controller.GetTaskList)

		vectorVerificationRoutes.GET("/statistics", controller.GetVerificationStatistics)

		vectorVerificationRoutes.GET("/progress", controller.GetRunningTaskProgress)

		vectorVerificationRoutes.POST("/repair", controller.RepairMissingVector)

		vectorVerificationRoutes.POST("/batch-repair", controller.BatchRepairMissingVectors)

		vectorVerificationRoutes.POST("/verify/:fileIds", controller.VerifySpecificVectors)

		vectorVerificationRoutes.DELETE("/task/:taskId", controller.CancelTask)

		vectorVerificationRoutes.POST("/smart-fix", controller.SmartFix)
	}

	fileRoutes := r.Group("/file")
	fileRoutes.Use(middleware.RequireAdmin())
	{
		fileRoutes.POST("/upload", fileController.UploadAdminFile)
	}

}
