package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yonraz/gochat_messages/controllers"
	"github.com/yonraz/gochat_messages/models"
	"gorm.io/gorm"
)

type MockService struct {
	DB *gorm.DB
}

func newMockMessagesService() *MockService {
	return &MockService{
		DB: &gorm.DB{},
	}
}
func (s *MockService) GetConversation(sender, receiver string) (*models.Conversation, error) {
	// Return a mock conversation for testing purposes
	return nil, nil
}

func (s *MockService) AddMessage(msg *models.Message) error {
	return nil
}
func (s *MockService) UpdateMessage(msg *models.Message) error {
	return nil
}
func (s *MockService) CreateConversation(sender string, receiver string) (*models.Conversation, error) {
	return nil, nil
}

func TestGetMessages(t *testing.T) {
    
	// Initialize MessagesService with the test database
	mockService := newMockMessagesService()
	controller := controllers.NewMessagesController(mockService)

	// Set up the router
	r := gin.Default()
	r.GET("/messages", controller.GetMessages)

	// Define test cases
	testCases := []struct {
		name               string
		queryParams        string
		expectedStatusCode int
		expectedBody       string
		expectConvNil bool
	}{
		{
			name:               "Missing query params",
			queryParams:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error":"missing sender and receiver query params"}`,
			expectConvNil: false,
		},
		{
			name:               "Valid query params",
			queryParams:        "?sender=foo&receiver=bar",
			expectedStatusCode: http.StatusOK,
			expectedBody:       `{"conv":{}}`,
			expectConvNil: 		true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("GET", "/messages"+tc.queryParams, nil)

			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatusCode, w.Code)

			var responseBody map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &responseBody)
            require.NoError(t, err, "Failed to unmarshal response body")

            if tc.expectConvNil {
                conv, exists := responseBody["conv"]
                assert.True(t, exists, "The 'conv' key should be present")
                assert.Nil(t, conv, "The 'conv' value should be nil")
            } else {
                _, exists := responseBody["conv"]
                assert.False(t, exists, "The 'conv' key should not be present")
            }
		})
	}
}
