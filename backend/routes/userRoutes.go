package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/controllers/userControllers"
	"github.com/pmohanj/web-chat-app/domain"
	"github.com/pmohanj/web-chat-app/middleware"
	"github.com/pmohanj/web-chat-app/mongo"
	"github.com/pmohanj/web-chat-app/repository"
	"github.com/pmohanj/web-chat-app/usecase"
)

func AddUserRoutes(router *gin.RouterGroup, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	userRouter := router.Group("/user")

	signUpRoute(userRouter, "/", env, timeout, db)
	loginRoute(userRouter, "/login", env, timeout, db)
	searchRoute(userRouter, "/search", env, timeout, db)
}

func signUpRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	sc := userControllers.SingUpController{
		SingupUseCase: usecase.NewSignupUseCase(ur, timeout),
		Env:           env,
	}

	r.POST(endPath, sc.SingUp)
}

func loginRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	lc := userControllers.LoginController{
		LoginUseCase: usecase.NewLoginUseCase(ur, timeout),
		Env:          env,
	}

	r.POST(endPath, lc.Login)
}

func searchRoute(r *gin.RouterGroup, endPath string, env *bootstrap.Env, timeout time.Duration, db mongo.Database) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)

	sc := userControllers.SearchController{
		SearchUseCase: usecase.NewSearchUseCase(ur, timeout),
	}

	r.GET(endPath, middleware.Authenticate(env.SecretKey), sc.SearchUsers)
}
