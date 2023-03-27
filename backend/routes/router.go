package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/websocket"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")

	AddUserRoutes(api)
	AddChatRoutes(api)
	AddMessageRoutes(api)

	// create websocketserver
	websocket := websocket.CreateWebSocketsServer()

	go websocket.SendMessage()
	AddWebScoketRouter(api, websocket)
}
