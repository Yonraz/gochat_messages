package consumers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_messages/constants"
	"github.com/yonraz/gochat_messages/initializers"
	"github.com/yonraz/gochat_messages/models"
	"github.com/yonraz/gochat_messages/services"
)

func NewMessageUpdatedConsumer(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel:     channel,
		srv:         services.NewMessagesService(initializers.DB),
		queueName:   string(constants.MessageReadQueue),
		routingKey:  string(constants.MessageReadKey),
		exchange:    string(constants.MessageEventsExchange),
		handlerFunc: MessageUpdatedHandler,
	}
}
func MessageUpdatedHandler(srv *services.MessagesService, msg amqp.Delivery) error {
	var parsed models.WsMessage

	if err := json.Unmarshal(msg.Body, &parsed); err != nil {
		log.Printf("error unmarshalling message: %v\n", err.Error())
		return err
	}

	fmt.Printf("message %v consumed on exchange %v with routing key %v\n", parsed, constants.MessageEventsExchange, constants.MessageReadKey)

	conv, err := srv.GetConversation(parsed.Sender, parsed.Receiver)
	if err != nil || conv == nil {
		log.Printf("error fetching conversation: %v\n", err)
		return err
	}

	existingMessage, err := srv.GetMessageByID(parsed.ID)
	if err != nil {
		log.Printf("error fetching message: %v\n", err)
		return err
	}

	if existingMessage == nil {
		return nil
	}
	isRead := string(parsed.Status) == string(constants.MessageReadKey)

	// Create a new message
	message := &models.Message{
		ID:             parsed.ID,
		Content:        parsed.Content,
		Sender:         parsed.Sender,
		Type:           parsed.Type,
		ConversationID: conv.ID,
		Receiver:       parsed.Receiver,
		Read:           isRead,
		Status:         parsed.Status,
		CreatedAt:      parsed.CreatedAt,
		UpdatedAt:      parsed.UpdatedAt,
		Version:        existingMessage.Version,
	}

	if parsed.Type == constants.MessageUpdate {
		_, err = srv.UpdateMessage(message)
	} else {
		err = fmt.Errorf("error processing: expected message type to be message.update, instead was: %v", parsed.Type)
		log.Printf("%v\n", err)
		return err
	}
	if err != nil {
		log.Printf("error updating message in db: %v\n", err)
		
		return err
	}
	
	


	log.Printf("messages service updated message %v", message)
	return nil
}