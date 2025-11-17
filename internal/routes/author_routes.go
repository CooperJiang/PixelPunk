package routes

import (
	authorController "pixelpunk/internal/controllers/author"

	"github.com/gin-gonic/gin"
)

func RegisterAuthorRoutes(r *gin.RouterGroup) {
	{
		r.GET("/:author_id", authorController.GetAuthorHomepage)

		r.GET("/:author_id/folders/:folder_id", authorController.GetAuthorFolder)
	}
}
