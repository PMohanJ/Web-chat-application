package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers/message"
	"github.com/pmohanj/web-chat-app/middleware"
)

func AddMessageRoutes(router *gin.RouterGroup) {
	messageRouter := router.Group("/message")

	messageRouter.POST("/", middleware.Authenticate(), message.SendMessage())
	messageRouter.GET("/:chatId", middleware.Authenticate(), message.GetMessages())
	messageRouter.PUT("/", middleware.Authenticate(), message.EditUserMessage())
	messageRouter.DELETE("/:messageId", middleware.Authenticate(), message.DeleteUserMessage())
}
