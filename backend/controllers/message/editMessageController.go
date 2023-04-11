package message

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EditMessageController struct {
	EditMessageUseCase domain.EditMessageUseCase
}

func (em *EditMessageController) EditMessage(c *gin.Context) {
	var reqData map[string]interface{}

	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Panic(err)
	}
	content, ok := reqData["content"].(string)
	if !ok {
		log.Panic("Type assertion failed")
	}

	mId, ok := reqData["messageId"].(string)
	if !ok {
		log.Panic("Type assertion failed")
	}

	messageId, err := primitive.ObjectIDFromHex(mId)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{{"_id", messageId}}
	update := bson.D{{"$set", bson.M{"content": content, "isedited": true}}}

	err = em.EditMessageUseCase.UpdateById(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while editing message"})
		log.Panic(err)
	}

	editedDocument, err := em.EditMessageUseCase.FetchById(c, "_id", messageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while retrieving data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, editedDocument[0])
}
