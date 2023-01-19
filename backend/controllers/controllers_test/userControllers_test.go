package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/joho/godotenv"
	"github.com/pmohanj/web-chat-app/database"
	"github.com/pmohanj/web-chat-app/models"
	"github.com/pmohanj/web-chat-app/routes"
)

var router *gin.Engine
var user1Token string
var chatId string
var user2Id string

func TestMain(m *testing.M) {
	router = gin.Default()

	err := godotenv.Load("./../../.env")
	if err != nil {
		log.Fatal("Error loading env variables ", err)
	}
	MongoDBURL := os.Getenv("MONGODB_URL_TESTING")
	// Initiate Databse
	database.DBinstance(MongoDBURL)

	// setup user routes
	api := router.Group("/api")
	routes.AddUserRoutes(api)
	routes.AddMessageRoutes(api)
	routes.AddChatRoutes(api)

	status := setupPhase()
	if status != 0 {
		os.Exit(1)
	}
	code := m.Run()
	tearDownPhase()
	database.CloseDBinstance()
	os.Exit(code)
}

// setupInitialPhase creates a user to for req resources needed to make
// subsequent operations like creating chat, send message, etc.
func setupPhase() int {
	input1 := []byte(`{"name":"User1", "email":"user1@gmail.com", "password":"haha123"}`)
	req1, _ := http.NewRequest("POST", "/api/user/", bytes.NewBuffer(input1))

	response1 := httptest.NewRecorder()
	router.ServeHTTP(response1, req1)
	if response1.Code != 200 {
		return response1.Code
	}

	var resUser1 map[string]string
	_ = json.NewDecoder(response1.Body).Decode(&resUser1)
	user1Token = resUser1["token"]
	fmt.Println(user1Token)

	input2 := []byte(`{"name":"User2", "email":"user2@gmail.com", "password":"haha123"}`)
	req2, _ := http.NewRequest("POST", "/api/user/", bytes.NewBuffer(input2))

	response2 := httptest.NewRecorder()
	router.ServeHTTP(response2, req2)
	if response2.Code != 200 {
		return response2.Code
	}

	var resUser2 map[string]string
	_ = json.NewDecoder(response2.Body).Decode(&resUser2)
	user2Id = resUser2["_id"]
	fmt.Println("UserID 2 ", user2Id)

	// initiate chat
	data := fmt.Sprintf(`{"userToBeAdded":"%s"}`, user2Id)
	fmt.Println("Str inp initiate chat", data)
	input3 := []byte(data)
	req3, _ := http.NewRequest("POST", "/api/chat/", bytes.NewBuffer(input3))
	req3.Header.Set("Authorization", "Bearer "+user1Token)
	//req3.Header.Set("Content-Type", "application/json")
	response3 := httptest.NewRecorder()
	router.ServeHTTP(response3, req3)
	if response3.Code != 200 {
		return response3.Code
	}

	var resChat map[string]interface{}
	_ = json.NewDecoder(response3.Body).Decode(&resChat)
	chatId, _ = resChat["_id"].(string)
	fmt.Println("Chat id ", chatId)
	return 0
}

func tearDownPhase() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	chatCollection := database.OpenCollection(database.Client, "chat")
	messageCollection := database.OpenCollection(database.Client, "message")
	userCollection := database.OpenCollection(database.Client, "user")

	chatCollection.Drop(ctx)
	messageCollection.Drop(ctx)
	userCollection.Drop(ctx)
}

func TestRegisterUser(t *testing.T) {

	t.Run("returns data decoding error", func(t *testing.T) {
		// improperly structured json input
		input := []byte(`{"name":"Sky, email":"checking@gmail.com "password" "78945626"}`)
		req, _ := http.NewRequest("POST", "/api/user/", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		if res["error"] != "error while decoding user data" {
			t.Errorf("Unexpected results: should return an error, error while decoding user data")
		}
	})

	t.Run("returns user already resgistered", func(t *testing.T) {
		input := []byte(`{"name":"User1", "email":"user1@gmail.com", "password":"haha123"}`)
		request, _ := http.NewRequest("POST", "/api/user/", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		if res["error"] != "You've already registered with this email" {
			t.Error("Unexpected result: should return error, user already registered with this email")
		}
	})
}

func TestAuthUser(t *testing.T) {

	t.Run("returns user deatils", func(t *testing.T) {
		input := []byte(`{"email":"user1@gmail.com", "password":"haha123"}`)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// expected response
		exp_data := models.User{
			Name:  "User1",
			Email: "user1@gmail.com",
			Pic:   "http://res.cloudinary.com/dkqc4za4f/image/upload/v1670340314/clsfmjxnuzsnidzc59np.jpg",
		}

		// decode the response body
		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusOK, response.Code)

		if res["name"] != exp_data.Name {
			t.Errorf("Unexpectes result: got %v, want %v", res["name"], exp_data.Name)
		}

		if res["email"] != exp_data.Email {
			t.Errorf("Unexpectes result: got %v, want %v", res["email"], exp_data.Email)
		}

		if res["pic"] == "" {
			t.Errorf("Unexpected result: pic field can't be empty")
		}

		if res["token"] == "" {
			t.Errorf("Unexpected result: token field can't be empty")
		}
	})

	t.Run("returns password invalid error", func(t *testing.T) {
		input := []byte(`{"email":"user1@gmail.com", "password":"haha"}`)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusUnauthorized, response.Code)

		if res["error"] != "Given password is invalid" {
			t.Errorf("Unexpected results: should return an error, invalid password")
		}
	})

	t.Run("returns user not resgistered error", func(t *testing.T) {
		input := []byte(`{"email":"unknown@gmail.com", "password":"12345678"}`)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusNotFound, response.Code)

		if res["error"] != "User not registered" {
			t.Errorf("Unexpected results: should return an error, user not registered")
		}
	})

	t.Run("returns data decoding error", func(t *testing.T) {
		// improperly structured json input
		input := []byte(`{"email"user1@gmail.com "password":"haha123"}`)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusBadRequest, response.Code)

		if res["error"] != "error while decoding user data" {
			t.Errorf("Unexpected results: should return an error, error while decoding user data")
		}
	})
}

func TestSearchUsers(t *testing.T) {

	t.Run("returns no users found", func(t *testing.T) {
		url := "/api/user/search?search=" + "uoaomaxoasvfa*#20"
		log.Println(url)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}
