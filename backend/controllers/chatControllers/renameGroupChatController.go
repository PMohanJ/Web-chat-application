package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RenameGroupChatController struct {
	RenameGroupChatUseCase domain.RenameGroupChatUseCase
}

func (rgc *RenameGroupChatController) RenameGroupChat(c *gin.Context) {
	var reqData map[string]interface{}

	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Panic(err)
	}

	groupName := reqData["groupName"].(string)
	cId := reqData["chatId"].(string)

	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{{"_id", chatId}}

	update := bson.D{{"$set", bson.D{{"chatName", groupName}}}}

	err = rgc.RenameGroupChatUseCase.UpdateById(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "rrror while updating data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"updatedGroupName": groupName})
}
