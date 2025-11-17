package routes

import (
	folderController "pixelpunk/internal/controllers/folder"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterFolderRoutes(r *gin.RouterGroup) {
	r.Use(middleware.RequireAuth())
	{
		r.POST("/create", folderController.CreateFolder)

		r.GET("/list", folderController.ListFolders)

		r.GET("/tree", folderController.GetFolderTree)

		r.GET("/search", folderController.SearchFolders)

		r.GET("/:folder_id", folderController.GetFolderDetail)

		r.POST("/update", folderController.UpdateFolder)
		r.DELETE("/:folder_id", folderController.DeleteFolder)

		r.GET("/contents", folderController.ListFolderContents)

		r.POST("/:folder_id/toggle-access-level", folderController.ToggleAccessLevel)

		r.POST("/reorder", folderController.ReorderFolders)

		r.POST("/move", folderController.MoveFolders)

		r.GET("/:folder_id/path-chain", folderController.GetFolderPathChain)

		r.POST("/batch-path-chains", folderController.GetBatchFolderPathChains)
	}
}
