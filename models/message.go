package models

import (
	"github.com/yonraz/gochat_messages/constants"
	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	ID             string                `json:"id" gorm:"primaryKey"`
	ConversationID uint                  `json:"conversationId" gorm:"index"`  // Use uint for foreign key
	Content        string                `json:"content"`
	Sender         string                `json:"sender"`
	Receiver       string                `json:"receiver"`
	Status         constants.RoutingKey  `json:"status"`
	Type           constants.MessageType `json:"type"`
	Read           bool                  `json:"read"`
	Sent           bool                  `json:"sent"`
}

type Conversation struct {
	gorm.Model
	Participants   []string       `json:"participants"`  // This should be better defined
	Messages       []Message  `json:"messages" gorm:"foreignKey:ConversationID"`
}
