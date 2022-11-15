package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/database"
	"github.com/pmohanj/web-chat-app/routes"
)

func main() {
	r := gin.Default()

	// Initiate Databse
	database.DBinstance()

	// Allows all origins, not suitable for prod environments
	r.Use(cors.Default())
	api := r.Group("/api")
	routes.AddUserRoutes(api)

	r.Run(":8000")
}
