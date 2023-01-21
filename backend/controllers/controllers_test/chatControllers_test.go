package controllers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestAddChatUser(t *testing.T) {
	t.Run("returns users chat", func(t *testing.T) {
		data := fmt.Sprintf(`{"userToBeAdded":"%s"}`, user2Id)
		input := []byte(data)
		request, _ := http.NewRequest("POST", "/api/chat/", bytes.NewBuffer(input))
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		expectedChatId := chatId

		var result map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, response.Code)

		if result["_id"] != expectedChatId {
			t.Errorf("Unexpected result: got %v, want %v", result["_id"], expectedChatId)
		}
	})
}

func TestGetUserChats(t *testing.T) {
	t.Run("returns user chats", func(t *testing.T) {
		request, _ := http.NewRequest("GET", "/api/chat/", nil)
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		var result []map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, response.Code)

		// wait 'ill change the actual handler so that I don't need to send id explicitly
		if len(result) < 1 {
			t.Errorf("Unexpected result: got %v, want %v", len(result), "at least 1 chat document")
		}
	})
}

func TestDeleteUserConversation(t *testing.T) {
	t.Run("returns status ok for delete conversation", func(t *testing.T) {
		request, _ := http.NewRequest("DELETE", "/api/chat/"+chatIdDelete, nil)
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}

func TestCreateGroupChat(t *testing.T) {
	t.Run("returns status ok for create group", func(t *testing.T) {
		data := fmt.Sprintf(`{"groupName":"Temporary testing group", "users":["%s"]}`, user2Id)
		input := []byte(data)
		request, _ := http.NewRequest("POST", "/api/chat/group", bytes.NewBuffer(input))
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		expectedChatName := "Temporary testing group"
		expectedChatLabel := true

		var result map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, response.Code)

		if result["chatName"] != expectedChatName {
			t.Errorf("Unexpected result: got %v, want %v", result["chatName"], expectedChatName)
		}

		if result["isGroupChat"] != expectedChatLabel {
			t.Errorf("Unexpected result: got %v, want %v", result["isGroupChat"], expectedChatLabel)
		}
	})
}

func TestRenameGroupChatName(t *testing.T) {
	t.Run("returns status ok for rename group", func(t *testing.T) {
		data := fmt.Sprintf(`{"groupName":"Group for testing renamed", "chatId":"%s"}`, chatIdGroup)
		input := []byte(data)
		request, _ := http.NewRequest("PUT", "/api/chat/grouprename", bytes.NewBuffer(input))
		request.Header.Set("Authorization", "Bearer "+user1Token)

		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)

		expectedChatName := "Group for testing renamed"

		var result map[string]interface{}
		_ = json.NewDecoder(response.Body).Decode(&result)

		assert.Equal(t, http.StatusOK, response.Code)

		if result["updatedGroupName"] != expectedChatName {
			t.Errorf("Unexpected result: got %v, want %v", result["chatName"], expectedChatName)
		}
	})
}
