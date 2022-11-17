package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/controllers"
)

func AddUserRoutes(router *gin.RouterGroup) {
	userRouter := router.Group("/user")

	userRouter.GET("/search", controllers.SearchUsers())
	userRouter.POST("/", controllers.RegisterUser())
	userRouter.POST("/login", controllers.AuthUser())
}
