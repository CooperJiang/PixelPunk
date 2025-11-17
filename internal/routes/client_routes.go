package routes

import (
	"pixelpunk/internal/middleware"
	"pixelpunk/internal/static"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterClientRoutes(r *gin.Engine) {
	distFS := static.GetDistFS()

	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api/") || strings.HasPrefix(path, "/debug/") {
			c.Next()
			return
		}
		middleware.StaticFileHandler(distFS)(c)
	})
}
