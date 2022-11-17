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

		// get the collection to perform query
		userCollection := database.OpenCollection(database.Client, "user")

		// check if user is already resgistered
		var temp models.User

		// If user doesn't exist, the following returns ErrNoDocuments
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&temp)

		// if err is other than ErrNoDocuments, something wrong while querying
		if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
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

func AuthUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		// should use other stuct as user login credential only containe email, pass...
		// will modify later
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Err while decoing data"})
			log.Printf("See the data: %+v", err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// get the collection to perform query
		userCollection := database.OpenCollection(database.Client, "user")
		var registeredUser models.UserResponse

		// check if user is a registered user
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&registeredUser)
		if err != nil && errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not registered"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for user data"})
			log.Panic(err)
		}

		log.Printf("UserResp form backed %+v", registeredUser)
		// user exist, check password validation
		errMsg, valid := helpers.VerifyPassword(registeredUser.Password, user.Password)
		if !valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
			return
		}

		// user is authorized
		c.JSON(http.StatusOK, registeredUser)
	}
}
