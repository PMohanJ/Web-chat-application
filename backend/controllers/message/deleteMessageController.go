package message

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteMessageController struct {
	DeleteMessageUseCase domain.DeleteMessageUseCase
}

func (dm *DeleteMessageController) DeleteMessage(c *gin.Context) {
	mId := c.Param("messageId")

	messageId, err := primitive.ObjectIDFromHex(mId)
	if err != nil {
		log.Panic(err)
	}

	err = dm.DeleteMessageUseCase.DeleteById(c, messageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while deleting data"})
		log.Panic(err)
	}

	c.Status(http.StatusOK)
}
