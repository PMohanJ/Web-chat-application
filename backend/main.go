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

	env := app.Env

	db := app.Mongo.Database(env.DatabaseName)
	defer app.CloseDBConnection()

	r := gin.Default()

	ORIGIN := env.Origin
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{ORIGIN},
		AllowMethods: []string{"PUT", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Content-Type", "Authorization"},
		MaxAge:       12 * time.Hour,
	}))

	timeout := time.Duration(env.ContextTimeout) * time.Second
	routes.SetupRoutes(r, env, timeout, db)

	r.Run(env.Port)
}
