package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers/user"
	"github.com/pmohanj/web-chat-app/middleware"
)

func AddUserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	userRouter.GET("/search", middleware.Authenticate(), user.SearchUsers())
	userRouter.POST("/", user.RegisterUser())
	userRouter.POST("/login", user.AuthUser())
}
