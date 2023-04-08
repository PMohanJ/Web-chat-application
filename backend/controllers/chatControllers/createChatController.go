package chatControllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateChatController struct {
	CreateChatUsecase domain.ChatRepository
}

func (cc *CreateChatController) CreateChat(c *gin.Context) {
	var ids map[string]interface{}

	if err := c.BindJSON(&ids); err != nil {
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{Message: "error while parsing data"})
		log.Println(err)
		return
	}

	// get the ids and convert them back to primitive.ObjectID format for querying
	id1, exists := c.Get("_id")
	if !exists {
		log.Panic("User details not available")
	}
	addingUser := id1.(primitive.ObjectID)

	id2 := ids["userToBeAdded"].(string)
	userToBeAdded, err := primitive.ObjectIDFromHex(id2)
	if err != nil {
		log.Panic(err)
	}

	filter := bson.D{
		{"isGroupChat", false},
		{"$and",
			bson.A{
				bson.D{{"users", bson.D{{"$elemMatch", bson.D{{"$eq", addingUser}}}}}},
				bson.D{{"users", bson.D{{"$elemMatch", bson.D{{"$eq", userToBeAdded}}}}}},
			},
		}}

	err = cc.CreateChatUsecase.FindByFilter(c, filter)
	if err == nil {
		// chat exist, return the chat
		matchStage := bson.D{
			{
				"$match", bson.D{{"isGroupChat", false},
					{"$and",
						bson.A{
							bson.D{{"users", bson.D{{"$elemMatch", bson.D{{"$eq", addingUser}}}}}},
							bson.D{{"users", bson.D{{"$elemMatch", bson.D{{"$eq", userToBeAdded}}}}}},
						}}},
			},
		}

		chat, err := cc.CreateChatUsecase.FetchByFilter(c, matchStage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while retrieving data"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, chat)
		return
	} else if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error in the server"})
		log.Panic(err)
	}

	// Chat doesn't exist, so create it
	createChat := domain.Chat{
		ChatName:    "sender",
		IsGroupChat: false,
		Users:       []primitive.ObjectID{addingUser, userToBeAdded},
	}

	insId, err := cc.CreateChatUsecase.Create(c, createChat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while inserting chat"})
		log.Panic(err)
	}

	chat, err := cc.CreateChatUsecase.FetchById(c, insId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, domain.ErrorResponse{Message: "error while fetching data"})
		log.Panic(err)
	}

	c.JSON(http.StatusOK, chat[0])
}
