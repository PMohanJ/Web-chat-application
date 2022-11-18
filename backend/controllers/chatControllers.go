package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/database"
	"github.com/pmohanj/web-chat-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AddChatUser lets the user to add a user to chat with
func AddChatUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ids map[string]interface{}

		if err := c.BindJSON(&ids); err != nil {
			log.Panic("error while parsing data")
		}

		log.Printf("ids: %+v", ids)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		chatCollection := database.OpenCollection(database.Client, "chat")

		// get the ids and convert them back to primitive.ObjectID format for querying
		id1 := ids["addingUser"].(string)
		addingUser, err := primitive.ObjectIDFromHex(id1)
		if err != nil {
			log.Panic(err)
		}

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

		// check if the users have chatted before, if so reaturn their chat
		var existedChat models.Chat
		err = chatCollection.FindOne(ctx, filter).Decode(&existedChat)
		if err == nil {
			c.JSON(http.StatusOK, existedChat)
			return
		} else if errors.Is(err, mongo.ErrNoDocuments) {
			log.Println("Panic...., no docs")

		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erroe while checking db"})
			log.Panic(err)
		}

		// No chat existed, so create it
		createChat := models.Chat{
			ChatName:    "sender",
			IsGroupChat: false,
			Users:       []primitive.ObjectID{addingUser, userToBeAdded},
		}

		insId, err := chatCollection.InsertOne(ctx, createChat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err while inserting chat"})
			log.Panic(err)
		}
		insertedId := insId.InsertedID.(primitive.ObjectID)
		log.Println(insertedId)

		var insertedChat models.Chat
		err = chatCollection.FindOne(ctx, bson.D{{"_id", insertedId}}).Decode(&insertedChat)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err while retreving inserted chat"})
		}

		c.JSON(http.StatusOK, insertedChat)
	}
}
