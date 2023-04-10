package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DeleteChatController struct {
	DeleteChatUseCase domain.DeleteChatUseCase
}

func (dc *DeleteChatController) DeleteUserChat(c *gin.Context) {
	cId := c.Param("chatId")
	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	// also required to delete messages that refer to this chat document
	err = dc.DeleteChatUseCase.DeleteById(c, chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while deleting chat document"})
		log.Panic(err)
	}

	c.Status(http.StatusOK)
}
