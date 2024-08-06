package models

import (
	"time"

	"github.com/lib/pq"
	"github.com/yonraz/gochat_messages/constants"
	"gorm.io/gorm"
)
type Message struct {
    ID             string                `json:"id" gorm:"type:uuid;primary_key"`
    ConversationID uint                  `json:"conversationId" gorm:"index"`
    Content        string                `json:"content"`
    Sender         string                `json:"sender"`
    Receiver       string                `json:"receiver"`
    Status         constants.RoutingKey  `json:"status"`
    Type           constants.MessageType `json:"type"`
    Read           bool                  `json:"read"`
    Sent           bool                  `json:"sent"`
    CreatedAt      time.Time             `json:"createdAt"`
    UpdatedAt      time.Time             `json:"updatedAt"`
    Version        uint                  `json:"version" gorm:"version"`
}
type WsMessage struct {
	ID      	string 					`json:"id" gorm:"primary key"`
	Content 	string					`json:"content"`	
	Sender 		string 					`json:"sender"`	
	Receiver 	string					`json:"receiver"`
	Status  	constants.RoutingKey	`json:"status"`	
	Type 		constants.MessageType	`json:"type"`
	Read 		bool					`json:"read"`
	Sent 		bool					`json:"sent"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Conversation struct {
	gorm.Model
	Participants   pq.StringArray `json:"participants" gorm:"type:text[]"`
	Messages       []Message    `json:"messages" gorm:"foreignKey:ConversationID"`
}
