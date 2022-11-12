package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/data"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.JSON(http.StatusOK, data.ChatsOfUsers[0])
	})

	r.Run(":8000")
}
