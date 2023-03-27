package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/bootstrap"
	"github.com/pmohanj/web-chat-app/routes"
)

func main() {

	app := bootstrap.App()
	defer app.CloseDBConnection()

	env := app.Env

	r := gin.Default()

	ORIGIN := env.Origin
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ORIGIN},
		AllowMethods: []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	routes.SetupRoutes(r)

	r.Run(env.Port)
}
