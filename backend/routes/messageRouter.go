package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers"
)

func AddMessageRoutes(router *gin.RouterGroup) {
	messageRouter := router.Group("/message")

	messageRouter.POST("/", controllers.SendMessage())
	messageRouter.GET("/:chatId", controllers.GetMessages())
}
