package messageControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GetMessagesController struct {
	GetMessagesUseCase domain.GetMessagesUseCase
}

func (gm *GetMessagesController) GetMessages(c *gin.Context) {
	cId := c.Param("chatId")

	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	messages, err := gm.GetMessagesUseCase.FetchById(c, "chat", chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while retrieving data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, messages)
}
