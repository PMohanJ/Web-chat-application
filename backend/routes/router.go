package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/mongo"
	"github.com/pmohanj/web-chat-app/websocket"
)

func SetupRoutes(gin *gin.Engine, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	api := gin.Group("/api")

	AddUserRoutes(api, env, timeout, db)
	AddChatRoutes(api, env, timeout, db)
	AddMessageRoutes(api, env, timeout, db)

	// create websocketserver
	websocket := websocket.CreateWebSocketsServer()

	go websocket.SendMessage()
	AddWebScoketRouter(api, websocket)
}
