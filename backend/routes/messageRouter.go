package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/controllers/message"
	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/mongo"
	"github.com/pmohanj/web-chat-app/repository"
	"github.com/pmohanj/web-chat-app/usecase"
)

func AddMessageRoutes(r *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	messageRouter := r.Group("/message")
	sendMessageRoute(messageRouter, "/", env, timeout, db)
	getMessagesRoute(messageRouter, "/:chatId", env, timeout, db)
	/*
			messageRouter.POST("/", middleware.Authenticate(), message.SendMessage())
			messageRouter.GET("/:chatId", middleware.Authenticate(), message.GetMessages())
			messageRouter.PUT("/", middleware.Authenticate(), message.EditUserMessage())
			messageRouter.DELETE("/:messageId", middleware.Authenticate(), message.DeleteUserMessage())
		}
	*/
}

func sendMessageRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	sm := &message.SendMessageController{
		SendMessageUseCase: usecase.NewSendMessageUseCase(cr, mr, timeout),
	}

	r.POST(endPath, sm.SendMessage)
}

func getMessagesRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	gm := &message.GetMessagesController{
		GetMessagesUseCase: usecase.NewGetMessagesUseCase(mr, timeout),
	}

	r.GET(endPath, gm.GetMessages)
}
