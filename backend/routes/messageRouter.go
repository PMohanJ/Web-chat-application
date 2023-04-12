package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/controllers/messageControllers"
	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/middleware"
	"github.com/pmohanj/web-chat-app/mongo"
	"github.com/pmohanj/web-chat-app/repository"
	"github.com/pmohanj/web-chat-app/usecase"
)

func AddMessageRoutes(r *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	messageRouter := r.Group("/message")
	messageRouter.Use(middleware.Authenticate(env.SecretKey))

	sendMessageRoute(messageRouter, "/", env, timeout, db)
	getMessagesRoute(messageRouter, "/:chatId", env, timeout, db)
	editMessageRoute(messageRouter, "/", env, timeout, db)
	deleteMessageRoute(messageRouter, "/:messageId", env, timeout, db)
}

func sendMessageRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	sm := &messageControllers.SendMessageController{
		SendMessageUseCase: usecase.NewSendMessageUseCase(cr, mr, timeout),
	}

	r.POST(endPath, sm.SendMessage)
}

func getMessagesRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	gm := &messageControllers.GetMessagesController{
		GetMessagesUseCase: usecase.NewGetMessagesUseCase(mr, timeout),
	}

	r.GET(endPath, gm.GetMessages)
}

func editMessageRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	em := &messageControllers.EditMessageController{
		EditMessageUseCase: usecase.NewEditMessageUseCase(mr, timeout),
	}

	r.PUT(endPath, em.EditMessage)
}

func deleteMessageRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	mr := repository.NewMessageRepository(db, domain.ColelctionMessage)

	dm := &messageControllers.DeleteMessageController{
		DeleteMessageUseCase: usecase.NewDeleteMessageUseCase(mr, timeout),
	}

	r.DELETE(endPath, dm.DeleteMessage)
}
