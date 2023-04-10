package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/controllers/chatControllers"
	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/mongo"
	"github.com/pmohanj/web-chat-app/repository"
	"github.com/pmohanj/web-chat-app/usecase"
	/* "github.com/pmohanj/web-chat-app/controllers/chat" */ /* "github.com/pmohanj/web-chat-app/middleware" */)

func AddChatRoutes(r *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	chatRouter := r.Group("/chat")

	createChatRoute(chatRouter, "/", env, timeout, db)
	getUserChatsRoute(chatRouter, "/", env, timeout, db)
	deleteUserChatRoute(chatRouter, "/:chatId", env, timeout, db)
	createGroupChatRoute(chatRouter, "/group", env, timeout, db)
	/*
		chatRouter.PUT("/grouprename", middleware.Authenticate(), chat.RenameGroupChatName())
		chatRouter.PUT("/groupadd", middleware.Authenticate(), chat.AddUserToGroupChat())
		chatRouter.PUT("/groupremove", middleware.Authenticate(), chat.DeleteUserFromGroupChat())
		chatRouter.PUT("/groupexit", middleware.Authenticate(), chat.UserExitGroup()) */
}

func createChatRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)

	cc := &chatControllers.CreateChatController{
		CreateChatUsecase: usecase.NewCreateChatUseCase(cr, timeout),
	}

	r.POST(endPath, cc.CreateChat)
}

func getUserChatsRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)

	uc := &chatControllers.UserChatsController{
		UserChatsUseCase: usecase.NewUserChatsUseCase(cr, timeout),
	}

	r.GET(endPath, uc.GetUserChats)
}

func deleteUserChatRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)

	dc := &chatControllers.DeleteChatController{
		DeleteChatUseCase: usecase.NewDeleteChatUseCase(cr, timeout),
	}

	r.DELETE(endPath, dc.DeleteUserChat)
}

func createGroupChatRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	cr := repository.NewChatRepository(db, domain.CollectionChat)

	dc := &chatControllers.GroupChatController{
		GroupChatUseCase: usecase.NewGroupChatUseCase(cr, timeout),
	}

	r.POST(endPath, dc.CreateGroupChat)
}
