package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yonraz/gochat_messages/constants"
	"github.com/yonraz/gochat_messages/controllers"
	"github.com/yonraz/gochat_messages/models"
	"gorm.io/gorm"
)

var sender = "foo"
var receiver = "bar"
var result = &models.Conversation{
    Participants: []string{sender, receiver},
    Messages: generateMockMessages(sender, receiver, 50),
}
var expected, err = json.Marshal(result)
var expectedStr = string(expected)

type MockService struct {
    DB *gorm.DB
}

func newMockMessagesService() *MockService {
    return &MockService{
        DB: &gorm.DB{},
    }
}

func (s *MockService) GetConversation(sender, receiver string) (*models.Conversation, error) {
    return nil, nil
}

func (s *MockService) GetConversationWithMessages(sender, receiver string, page int) (*models.Conversation, error) {
    pageSize := 25
    start := (page - 1) * pageSize
    end := start + pageSize
    
    if start >= len(result.Messages) {
        return &models.Conversation{
            Participants: result.Participants,
            Messages:     []models.Message{},
        }, nil
    }
    
    if end > len(result.Messages) {
        end = len(result.Messages)
    }

    msgs := result.Messages[start:end]
    ret := &models.Conversation{
        Participants: result.Participants,
        Messages:     msgs,
    }
    return ret, nil
}


func (s *MockService) AddMessage(msg *models.Message) error {
    return nil
}

func (s *MockService) UpdateMessage(message *models.Message) (*models.Message, error) {
    return nil, nil
}

func (s *MockService) CreateConversation(sender string, receiver string) (*models.Conversation, error) {
    return nil, nil
}

func (s *MockService) GetMessageByID(id string) (*models.Message, error) {
    return nil, nil
}

func TestGetMessages(t *testing.T) {
    mockService := newMockMessagesService()
    controller := controllers.NewMessagesController(mockService)

    r := gin.Default()
    r.GET("/api/messages", controller.GetMessages)

    testCases := []struct {
        name               string
        queryParams        string
        expectedStatusCode int
        expectedBody       string
        expectConvResponse      bool
    }{
        {
            name:               "Missing query params",
            queryParams:        "",
            expectedStatusCode: http.StatusBadRequest,
            expectedBody:       `{"error":"missing sender and receiver query params"}`,
            expectConvResponse:      false,
        },
        {
            name:               "Valid query params without pagination",
            queryParams:        "?sender=foo&receiver=bar",
            expectedStatusCode: http.StatusOK,
            expectedBody:       expectedStr,
            expectConvResponse:      true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            req, _ := http.NewRequest("GET", "/api/messages"+tc.queryParams, nil)
            w := httptest.NewRecorder()

            r.ServeHTTP(w, req)

            assert.Equal(t, tc.expectedStatusCode, w.Code)

            var responseBody map[string]interface{}
            err := json.Unmarshal(w.Body.Bytes(), &responseBody)
            require.NoError(t, err, "Failed to unmarshal response body")

            _, exists := responseBody["conv"]
            if tc.expectConvResponse {
				var res map[string]models.Conversation
        		err := json.Unmarshal(w.Body.Bytes(), &res)
				require.NoError(t, err)
                assert.True(t, exists, "The 'conv' key should be present")
                assert.Equal(t, len(res["conv"].Messages), 25)

				assert.Equal(t, "msg-1", res["conv"].Messages[0].ID)
				assert.Equal(t, "msg-25", res["conv"].Messages[24].ID)
            } else {
                assert.False(t, exists, "The 'conv' key should not be present")
            }
        })
    }
	t.Run("Return 25 messages with pagination 2 value", func(t *testing.T) {
        req, _ := http.NewRequest("GET", "/api/messages?sender=foo&receiver=bar&page=2", nil)
        w := httptest.NewRecorder()

        r.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)

       var responseBody map[string]models.Conversation
        err := json.Unmarshal(w.Body.Bytes(), &responseBody)
        require.NoError(t, err, "Failed to unmarshal response body")

        assert.Equal(t, len(responseBody["conv"].Messages), 25)

        // Additional checks for pagination correctness
        assert.Equal(t, "msg-26", responseBody["conv"].Messages[0].ID)
        assert.Equal(t, "msg-50", responseBody["conv"].Messages[24].ID)
    })
}

func generateMockMessages(sender, receiver string, amount int) []models.Message {
	var messages []models.Message

	for i := 1; i <= amount; i++ {
		msg := models.Message{
			ID:            fmt.Sprintf("msg-%d", i),
			ConversationID: 1, // Assuming all messages belong to the same conversation
			Content:       fmt.Sprintf("Message content %d", i),
			Sender:        sender,
			Receiver:      receiver,
			Status:        constants.MessageSentKey,
			Type:          constants.MessageCreate,
			Read:          i%2 == 0,  // Every second message is read
			Sent:          true,
			CreatedAt:     time.Now().Add(-time.Duration(50-i) * time.Minute), // Spread messages over time
			UpdatedAt:     time.Now(),
			Version:       1,
		}
		messages = append(messages, msg)
	}

	return messages
}

