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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while parsing data"})
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

func EditUserMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while parsing data"})
			log.Panic(err)
		}
		content, ok := reqData["content"].(string)
		if !ok {
			log.Panic("Error type assertion")
		}

		mId, ok := reqData["messageId"].(string)
		if !ok {
			log.Panic("Error type assertion")
		}

		messageId, err := primitive.ObjectIDFromHex(mId)
		if err != nil {
			log.Panic(err)
		}

		messageCollection := database.OpenCollection(database.Client, "message")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		filter := bson.D{{"_id", messageId}}
		update := bson.D{{"$set", bson.M{"content": content, "isedited": true}}}

		// return the document after it's modified
		options := options.FindOneAndUpdate().SetReturnDocument(options.After)

		var updatedDoc bson.M
		err = messageCollection.FindOneAndUpdate(ctx, filter, update, options).Decode(&updatedDoc)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting message"})
			log.Panic(err)
		}

		pipeline := mongo.Pipeline{
			bson.D{
				{
					"$match", updatedDoc,
				},
			},

			bson.D{
				{
					"$lookup", bson.D{
						{"from", "user"},
						{"localField", "sender"},
						{"foreignField", "_id"},
						{"as", "sender"},
					},
				},
			},

			bson.D{
				{
					"$project", bson.D{
						{"sender.password", 0},
						{"created_at", 0},
						{"updated_at", 0},
						{"sender.created_at", 0},
						{"sender.updated_at", 0},
					},
				},
			},
		}

		cursor, err := messageCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		var res []bson.M
		if err := cursor.All(ctx, &res); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while retrieving data"})
			log.Panic(err)
		}

		log.Printf("res[0] edited msg %+v", res[0])
		c.JSON(http.StatusOK, res[0])
	}
}

func DeleteUserMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		mId := c.Param("messageId")

		messageId, err := primitive.ObjectIDFromHex(mId)
		if err != nil {
			log.Panic(err)
		}

		messageCollection := database.OpenCollection(database.Client, "message")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		deleteRes, err := messageCollection.DeleteOne(ctx, bson.D{{"_id", messageId}})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while deleting message"})
			log.Panic(err)
		}

		log.Println("Delete the msg: ", deleteRes.DeletedCount)

		c.Status(http.StatusOK)
	}
}
