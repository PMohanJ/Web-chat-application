package messageControllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SendMessageController struct {
	SendMessageUseCase domain.SendMessageUseCase
}

func (sm *SendMessageController) SendMessage(c *gin.Context) {
	var reqData map[string]interface{}

	if err := c.BindJSON(&reqData); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Println(err)
		return
	}

	cId := reqData["chatId"].(string)
	content := reqData["content"].(string)

	chatId, err := primitive.ObjectIDFromHex(cId)
	if err != nil {
		log.Panic(err)
	}

	//senderId refers to the user who's sending the message
	sId, exists := c.Get("_id")
	if !exists {
		log.Panic("User details not available")
	}

	senderId := sId.(primitive.ObjectID)
	newMessage := domain.Message{
		Sender:  senderId,
		Content: content,
		Chat:    chatId,
	}

	insId, err := sm.SendMessageUseCase.Create(c, newMessage)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while inserting data"})
		log.Panic(err)
	}

	// update the latestMessage field
	filter := bson.D{{"_id", chatId}}
	update := bson.D{{"$set", bson.D{{"latestMessage", insId}}}}

	err = sm.SendMessageUseCase.UpdateByFilter(c, filter, update)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while updating data"})
		log.Panic(err)
	}

	// get the inserted message document, and send it to client
	document, err := sm.SendMessageUseCase.FetchById(c, "_id", insId)
	if err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while retrieving data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, document[0])
}
