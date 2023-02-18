package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers/chat"
	"github.com/pmohanj/web-chat-app/middleware"
)

func AddChatRoutes(r *gin.RouterGroup) {
	chatRouter := r.Group("/chat")
	chatRouter.POST("/", middleware.Authenticate(), chat.AddChatUser())
	chatRouter.GET("/", middleware.Authenticate(), chat.GetUserChats())
	chatRouter.DELETE("/:chatId", middleware.Authenticate(), chat.DeleteUserConversation())
	chatRouter.POST("/group", middleware.Authenticate(), chat.CreateGroupChat())
	chatRouter.PUT("/grouprename", middleware.Authenticate(), chat.RenameGroupChatName())
	chatRouter.PUT("/groupadd", middleware.Authenticate(), chat.AddUserToGroupChat())
	chatRouter.PUT("/groupremove", middleware.Authenticate(), chat.DeleteUserFromGroupChat())
	chatRouter.PUT("/groupexit", middleware.Authenticate(), chat.UserExitGroup())
}
