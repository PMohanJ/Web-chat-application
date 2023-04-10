package chatControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AddUserToGroupChatController struct {
	AddUserToGroupChatUseCase domain.AddUserToGroupChatUseCase
}

func (aug *AddUserToGroupChatController) AddUserToGroupChat(c *gin.Context) {
	var reqData map[string]interface{}

	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Panic(err)
	}

	uId := reqData["userId"].(string)
	cId := reqData["chatId"].(string)

	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	userId, err := primitive.ObjectIDFromHex(uId)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{{"_id", chatId}}

	update := bson.D{{"$push", bson.D{{"users", userId}}}}

	err = aug.AddUserToGroupChatUseCase.UpdateById(c, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while updating data"})
		log.Panic(err)
	}

	updateChat, err := aug.AddUserToGroupChatUseCase.FetchById(c, chatId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error retrieving data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, updateChat[0])
}
