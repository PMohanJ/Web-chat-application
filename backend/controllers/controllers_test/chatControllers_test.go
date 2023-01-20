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
