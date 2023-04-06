package userControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
)

type SearchController struct {
	SearchUseCase domain.SearchUseCase
}

func (sc *SearchController) SearchUsers(c *gin.Context) {
	query := c.Query("search")

	users, err := sc.SearchUseCase.SearchUsers(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error occured in the server"})
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, users)
}
