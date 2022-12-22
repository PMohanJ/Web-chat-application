package controllers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/pmohanj/web-chat-app/database"
	"github.com/pmohanj/web-chat-app/models"
	"github.com/pmohanj/web-chat-app/routes"
)

var router *gin.Engine
var tempToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJRCI6IjYzOGY1ZWRlMjUyODI2MjBmMDAwYzVmZiIsIk5hbWUiOiIiLCJFbWFpbCI6ImNoZWNraW5nQGdtYWlsLmNvbSIsImV4cCI6MTY3MTcxNjQ1OCwiaWF0IjoxNjcxNjMwMDU4fQ.tY-pU9duArH5Iq_wX6sChfbGhttqcj3U0VVzalZ8Vs0"

func TestMain(m *testing.M) {
	router = gin.Default()

	// setup necessary env variables
	os.Setenv("MONGODB_URL", "mongodb+srv://mohanj:webchatapp01@cluster0.f2pstnw.mongodb.net/?retryWrites=true&w=majority")
	os.Setenv("SECRET_KEY", "itsnotpossibletomanipulate000000")

	// Initiate Databse
	database.DBinstance()

	// setup user routes
	api := router.Group("/api")
	routes.AddUserRoutes(api)
	routes.AddMessageRoutes(api)

	code := m.Run()
	database.CloseDBinstance()
	os.Exit(code)
}

func TestAuthUser(t *testing.T) {

	t.Run("returns user deatils", func(t *testing.T) {
		input := []byte(`{"email":"checking@gmail.com", "password":"haha123"}`)
		req, _ := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(input))

		response := httptest.NewRecorder()
		router.ServeHTTP(response, req)

		// expected response
		exp_data := models.User{
			Name:  "Checking",
			Email: "checking@gmail.com",
			Pic:   "http://res.cloudinary.com/dkqc4za4f/image/upload/v1670340314/clsfmjxnuzsnidzc59np.jpg",
		}

		// decode the response body
		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusOK, response.Code)

		if res["name"] != "Checking" {
			t.Errorf("Unexpectes result: got %v, want %v", res["name"], exp_data.Name)
		}

		if res["email"] != "checking@gmail.com" {
			t.Errorf("Unexpectes result: got %v, want %v", res["email"], exp_data.Email)
		}

		if res["pic"] == "" {
			t.Errorf("Unexpected result: pic field can't be empty")
		}

		if res["token"] == "" {
			t.Errorf("Unexpected result: token field can't be empty")
		}
	})

	t.Run("returns password invalid error ", func(t *testing.T) {
		input := []byte(`{"email":"checking@gmail.com", "password":"haha"}`)
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
		input := []byte(`{"email":"checking@gmail.com "password":"haha123"}`)
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
		input := []byte(`{"name":"Checking", "email":"checking@gmail.com", "password":"haha123"}`)
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

func TestSearchUsers(t *testing.T) {

	t.Run("returns no users found", func(t *testing.T) {
		url := "/api/user/search?search=" + "uoaomaxoasvfa*#20"
		log.Println(url)
		request, _ := http.NewRequest("GET", url, nil)
		request.Header.Set("Authorization", "Bearer "+tempToken)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var res map[string]string
		_ = json.NewDecoder(response.Body).Decode(&res)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}
