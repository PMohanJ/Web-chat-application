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
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterUser will register the new users to application
func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if user.Pic == "" {
			user.SetDefaultPic()
		}
		log.Printf("User data %+v", user)

		var ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// check if user is already resgistered
		database.DBinstance()
		dbClient := database.Client
		userCollection := dbClient.Database("cluster0").Collection("user")

		var temp models.User
		// If user doesn't exist, the following returns ErrNoDocuments
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&temp)

		// if err is other than ErrNoDocuments, something wrong while querying
		if err != nil && errors.Is(errors.Unwrap(err), mongo.ErrNoDocuments) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err occured while checking for user"})
			log.Panic(err)
		}

		if temp.Email == user.Email {
			c.JSON(http.StatusBadRequest, gin.H{"error": "You've alrady registered with this email"})
			return
		}

		// user doesn't exist in database, so register
		hashedPassowrd := helpers.HashPassowrd(user.Password)
		user.Password = hashedPassowrd

		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "err occured while registering the user"})
			log.Panic(err)
		} else {
			c.JSON(http.StatusOK, user)
		}

	}
}
