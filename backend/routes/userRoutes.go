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
)

func AddUserRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	userRouter := router.Group("/user")

	signUpRoute(userRouter, "/", env, timeout, db)
	loginRoute(userRouter, "/login", env, timeout, db)
	/* userRouter.GET("/search", middleware.Authenticate(), user.SearchUsers())
	userRouter.POST("/", user.RegisterUser())
	userRouter.POST("/login", user.AuthUser()) */

}

func signUpRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	sc := chatControllers.SingUpController{
		SingupUseCase: usecase.NewSignupUseCase(ur, timeout),
		Env:           env,
	}

	r.POST(endPath, sc.SingUp)
}

func loginRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	lc := chatControllers.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}

	r.POST(endPath, lc.Login)
}
