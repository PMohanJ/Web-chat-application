package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers"
)

func AddChatRoutes(r *gin.RouterGroup) {
	chat := r.Group("/chat")
	chat.POST("/", controllers.AddOChatUser())
	chat.GET("/:userId", controllers.GetUserChats())
	chat.POST("/group", controllers.CreateGroupChat())
}
