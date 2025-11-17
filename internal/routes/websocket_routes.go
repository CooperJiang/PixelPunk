package routes

import (
	"pixelpunk/internal/controllers/websocket"
	"pixelpunk/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWebSocketRoutes(r *gin.RouterGroup) {
	wsGroup := r.Group("/ws")
	wsGroup.Use(middleware.JWTAuth())
	wsGroup.Use(middleware.RequireAuth())
	{
		wsGroup.GET("/admin", websocket.HandleWebSocket)

		wsGroup.GET("/stats", websocket.GetStats)
	}
}
