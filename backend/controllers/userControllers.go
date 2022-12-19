package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pmohanj/web-chat-app/database"
	"github.com/pmohanj/web-chat-app/helpers"
	"github.com/pmohanj/web-chat-app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUser will register the new users to application
func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while decoding user data"})
			log.Println(err)
			return
		}

		if user.Pic == "" {
			user.SetDefaultPic()
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// get the collection to perform querying
		userCollection := database.OpenCollection(database.Client, "user")

		// check if user is already resgistered
		var temp models.User

		// If user doesn't exist, the following returns ErrNoDocuments
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&temp)

		// if err is other than ErrNoDocuments, something wrong while querying
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while querying for user"})
			log.Panic(err)
		}

		if temp.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You've already registered with this email"})
			return
		}

		// user doesn't exist in database, so register the user
		hashedPassowrd := helpers.HashPassowrd(user.Password)
		user.Password = hashedPassowrd

		insId, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while registering the user"})
			log.Panic(err)
		}

		user.Id = insId.InsertedID.(primitive.ObjectID)
		id := user.Id.Hex()
		// generate token for the user
		if user.Token, err = helpers.GenerateToken(id, user.Name, user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, user)
	}
}

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "error while decoding user data"})
			log.Println(err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// get the collection to perform querying
		userCollection := database.OpenCollection(database.Client, "user")
		var registeredUser models.UserResponse

		// check if user is a registered user
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&registeredUser)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not registered"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while querying user data"})
			log.Panic(err)
		}

		log.Printf("user data from database %+v", registeredUser)
		// user exist, check for password validation
		errMsg, valid := helpers.VerifyPassword(registeredUser.Password, user.Password)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		// generate token for the user
		if registeredUser.Token, err = helpers.GenerateToken(registeredUser.Id.Hex(), user.Name, user.Email); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to generate token"})
			log.Panic(err)
		}

		c.JSON(http.StatusOK, registeredUser)
	}
}

func SearchUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("search")
		log.Println(query)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// get the user collection
		userCollection := database.OpenCollection(database.Client, "user")
		filter := bson.D{
			{"$or",
				bson.A{
					bson.D{{"name", bson.D{{"$regex", query}}}},
					bson.D{{"email", bson.D{{"$regex", query}}}},
				},
			},
		}

		cursor, err := userCollection.Find(ctx, filter)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			c.Status(http.StatusOK)
			return
		}

		var results []models.UserResponse
		for cursor.Next(ctx) {
			var temp models.UserResponse
			if err := cursor.Decode(&temp); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error in the server"})
				log.Panic(err)
			}
			results = append(results, temp)
		}

		c.JSON(http.StatusOK, results)

	}
}
