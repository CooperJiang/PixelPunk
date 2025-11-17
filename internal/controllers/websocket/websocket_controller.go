package websocket

import (
	"net/http"
	"pixelpunk/internal/services/auth"
	ws "pixelpunk/internal/websocket"
	"pixelpunk/pkg/errors"
	"pixelpunk/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	globalManager *ws.Manager
)

func InitWebSocketManager() {
	config := ws.DefaultConfig()
	globalManager = ws.NewManager(config)
	globalManager.Start()
}

func GetWebSocketManager() *ws.Manager {
	return globalManager
}

func HandleWebSocket(c *gin.Context) {

	claims, exists := c.Get("payload")
	if !exists {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未找到用户信息"))
		return
	}

	jwtClaims, ok := claims.(*auth.JWTClaims)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeInvalidRequest, "用户信息格式错误"))
		return
	}

	if jwtClaims.Role != 1 { // 假设1为管理员角色
		errors.HandleError(c, errors.New(errors.CodeForbidden, "需要管理员权限"))
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WebSocket升级失败: %v", err)
		return
	}

	client := ws.NewClient(conn, jwtClaims.UserID, jwtClaims.Role == 1)

	globalManager.RegisterClient(client)

	go client.WritePump()
	go client.ReadPump(globalManager)
}

func BroadcastMessage(msgType ws.MessageType, data interface{}) {
	if globalManager == nil {
		return
	}

	msg := ws.NewMessage(msgType, data)
	globalManager.BroadcastMessage(msg)
}

func BroadcastToAdmins(msgType ws.MessageType, data interface{}) {
	if globalManager == nil {
		return
	}

	msg := ws.NewMessage(msgType, data)
	globalManager.SendToAdmins(msg)

}

func SendToClient(clientID string, msgType ws.MessageType, data interface{}) error {
	if globalManager == nil {
		return errors.New(errors.CodeInternal, "WebSocket管理器未初始化")
	}

	msg := ws.NewMessage(msgType, data)
	return globalManager.SendToClient(clientID, msg)
}

func GetStats(c *gin.Context) {
	if globalManager == nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "WebSocket管理器未初始化"))
		return
	}

	stats := globalManager.GetStats()
	errors.ResponseSuccess(c, stats, "获取统计信息成功")
}
