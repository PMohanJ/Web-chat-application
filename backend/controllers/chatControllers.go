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
func AddOChatUser() gin.HandlerFunc {
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
			// chat exist, perform aggragrate operations to join document of Chat with respectice chat Users profile
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
			lookupStage := bson.D{
				{
					"$lookup", bson.D{{"from", "user"}, {"localField", "users"}, {"foreignField", "_id"}, {"as", "users"}},
				},
			}

			projectStage := bson.D{
				{
					"$project", bson.D{
						{"users.password", 0}, {"users.isAdmin", 0}, {"created_at", 0}, {"updated_at", 0}, {"users.created_at", 0}, {"users.updated_at", 0},
					},
				},
			}
			var res []bson.M
			cur, err := chatCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
			if err != nil {
				log.Panic(err)
			}

			if err = cur.All(ctx, &res); err != nil {
				log.Panic(err)
			}
			for _, docu := range res {
				log.Printf("docu: %+v", docu)
			}

			c.JSON(http.StatusOK, res)
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

func GetUserChats() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("userId")

		userId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while converting hex to objId"})
		}
		// get chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		matchStage := bson.D{
			{
				"$match", bson.D{
					{
						"users", bson.D{{"$elemMatch", bson.D{{"$eq", userId}}}},
					},
				},
			},
		}

		lookupStage := bson.D{
			{
				"$lookup", bson.D{
					{"from", "user"},
					{"localField", "users"},
					{"foreignField", "_id"},
					{"as", "users"},
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

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		cursor, err := chatCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking documents"})
			log.Panic(err)
		}

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			log.Panic(err)
		}
		for _, docu := range results {
			log.Println(docu)
		}

		c.JSON(http.StatusOK, results)
	}
}

func CreateGroupChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		var groupData map[string]interface{}

		if err := c.BindJSON(&groupData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing data"})
			log.Panic(err)
		}

		groupName := groupData["groupName"].(string)

		// assuming that the admin user id is sent separately as JWT is not implemented yet,
		// we can't exactly distinguish between normal and admin user
		users := groupData["users"].([]interface{})
		aUser := groupData["adminUser"].(string)

		var usersIds []primitive.ObjectID
		adminUser, err := primitive.ObjectIDFromHex(aUser)
		if err != nil {
			log.Panic(err)
		}
		usersIds = append(usersIds, adminUser)

		for _, uId := range users {
			id := uId.(string)

			temp, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				log.Panic(err)
			}
			usersIds = append(usersIds, temp)
		}

		groupChat := models.Chat{
			IsGroupChat: true,
			ChatName:    groupName,
			Users:       usersIds,
			GroupAdmin:  adminUser,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		// get the chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		insId, err := chatCollection.InsertOne(ctx, groupChat)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err while inserting document"})
			log.Panic(insId)
		}

		insertedId := insId.InsertedID.(primitive.ObjectID)
		groupChat.Id = insertedId

		c.JSON(http.StatusOK, groupChat)
	}
}

func RenameGroupChatName() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing data"})
			return
		}

		groupName := reqData["groupName"].(string)
		cId := reqData["chatId"].(string)

		chatId, err := primitive.ObjectIDFromHex(cId)
		if err != nil {
			log.Panic(err)
		}

		// get chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		filter := bson.D{{"_id", chatId}}

		update := bson.D{{"$set", bson.D{{"chatName", groupName}}}}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err = chatCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating document"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, gin.H{"updatedGroupName": groupName})
	}
}

func AddUserToGroupChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing data"})
			return
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

		// get chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		filter := bson.D{{"_id", chatId}}

		update := bson.D{{"$push", bson.D{{"users", userId}}}}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		_, err = chatCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating document"})
			log.Panic(err)
		}

		// User is added to group, now retrieve that document and send into client
		// so that client can update its data, and perfrom necessary rendering
		matchStage := bson.D{
			{
				"$match", bson.D{
					{
						"_id", chatId,
					},
				},
			},
		}

		lookupStage := bson.D{
			{
				"$lookup", bson.D{
					{"from", "user"},
					{"localField", "users"},
					{"foreignField", "_id"},
					{"as", "users"},
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

		cursor, err := chatCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking documents"})
			log.Panic(err)
		}

		// we can only pass array type to cursor even though we know
		// we're just retrieving single document
		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			log.Panic(err)
		}
		for _, docu := range results {
			log.Println(docu)
		}

		// send the document
		c.JSON(http.StatusOK, results[0])
	}
}

func DeleteUserFromGroupChat() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing data"})
			return
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

		// get chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		filter := bson.D{{"_id", chatId}}

		update := bson.D{{"$pull", bson.D{{"users", userId}}}}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		res, err := chatCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating document"})
			log.Panic(err)
		}
		log.Printf("Docu up %v", res.ModifiedCount)
		// User is added to group, now retrieve that document and send into client
		// so that client can update its data, and perfrom necessary rendering
		matchStage := bson.D{
			{
				"$match", bson.D{
					{
						"_id", chatId,
					},
				},
			},
		}

		lookupStage := bson.D{
			{
				"$lookup", bson.D{
					{"from", "user"},
					{"localField", "users"},
					{"foreignField", "_id"},
					{"as", "users"},
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

		cursor, err := chatCollection.Aggregate(ctx, mongo.Pipeline{matchStage, lookupStage, projectStage})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking documents"})
			log.Panic(err)
		}

		var results []bson.M
		if err := cursor.All(ctx, &results); err != nil {
			log.Panic(err)
		}
		for _, docu := range results {
			log.Println(docu)
		}

		// send the document
		c.JSON(http.StatusOK, results[0])
	}

}

func UserExitGroup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqData map[string]interface{}

		if err := c.BindJSON(&reqData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Error while parsing data"})
			return
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

		// get chat collection
		chatCollection := database.OpenCollection(database.Client, "chat")

		filter := bson.D{{"_id", chatId}}

		update := bson.D{{"$pull", bson.D{{"users", userId}}}}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		res, err := chatCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while updating document"})
			log.Panic(err)
		}
		log.Printf("Docu up %v", res.ModifiedCount)

		c.JSON(http.StatusOK, gin.H{"message": "Exited from group"})
	}
}
