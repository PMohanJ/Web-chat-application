package controllers

import (
	"context"
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

func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while parsing data "})
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
		newMessage := models.Message{
			Sender:  senderId,
			Content: content,
			Chat:    chatId,
		}

		// get the message collection
		messageCollection := database.OpenCollection(database.Client, "message")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		insId, err := messageCollection.InsertOne(ctx, newMessage)
		insertedId := insId.InsertedID.(primitive.ObjectID)

		// get chat collection to update the latestMessage field
		chatCollection := database.OpenCollection(database.Client, "chat")

		filter := bson.D{{"_id", chatId}}
		update := bson.D{{"$set", bson.D{{"latestMessage", insertedId}}}}
		_, err = chatCollection.UpdateOne(ctx, filter, update)

		// get the inserted message document, ans send it to client
		matchStage := bson.D{
			{
				"$match", bson.D{
					{
						"_id", insertedId,
					},
				},
			},
		}

		lookupStage := bson.D{
			{
				"$lookup", bson.D{
					{"from", "user"},
					{"localField", "sender"},
					{"foreignField", "_id"},
					{"as", "sender"},
				},
			},
		}

		projectStage := bson.D{
			{
				"$project", bson.D{
					{"sender.password", 0},
					{"created_at", 0},
					{"updated_at", 0},
					{"sender.created_at", 0},
					{"sender.updated_at", 0},
				},
			},
		}

		cursor, err := messageCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, results[0])
	}
}

func GetMessages() gin.HandlerFunc {
	return func(c *gin.Context) {
		cId := c.Param("chatId")

		chatId, err := primitive.ObjectIDFromHex(cId)
		if err != nil {
			log.Panic(err)
		}

		// get messages collection
		messageCollection := database.OpenCollection(database.Client, "message")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		matchStage := bson.D{
			{
				"$match", bson.D{
					{
						"chat", chatId,
					},
				},
			},
		}

		lookupStage := bson.D{
			{
				"$lookup", bson.D{
					{"from", "user"},
					{"localField", "sender"},
					{"foreignField", "_id"},
					{"as", "sender"},
				},
			},
		}

		projectStage := bson.D{
			{
				"$project", bson.D{
					{"users.password", 0},
					{"created_at", 0},
					{"updated_at", 0},
					{"users.created_at", 0},
					{"users.updated_at", 0},
				},
			},
		}

		cursor, err := messageCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, results)
	}
}

func DeleteUserConversation() gin.HandlerFunc {
	return func(c *gin.Context) {
		cId := c.Param("chatId")
		chatId, err := primitive.ObjectIDFromHex(cId)
		if err != nil {
			log.Panic(err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// delete all the message that refer the chatId
		messageCollection := database.OpenCollection(database.Client, "message")

		filter := bson.D{
			{"chat", chatId},
		}
		_, err = messageCollection.DeleteMany(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting chat messages"})
			log.Panic(err)
		}

		// delete the chat document too
		chatCollection := database.OpenCollection(database.Client, "chat")
		_, err = chatCollection.DeleteOne(ctx, bson.D{{"_id", chatId}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting chat document"})
			log.Panic(err)
		}

		c.Status(http.StatusOK)
	}
}
