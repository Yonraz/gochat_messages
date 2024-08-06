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

func NewMessageSentConsumer(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel: channel,
		srv: services.NewMessagesService(initializers.DB),
		queueName: string(constants.MessageSentQueue),
		routingKey: string(constants.MessageSentKey),
		exchange: string(constants.MessageEventsExchange),
		handlerFunc: MessageSentHanlder,
	}
}

func MessageSentHanlder(srv *services.MessagesService, msg amqp.Delivery) error {
	var parsed models.WsMessage

	if err := json.Unmarshal(msg.Body, &parsed); err != nil {
		log.Printf("error unmarshalling message: %v\n", err.Error())
		return err
	}

	fmt.Printf("message %v consumed on exchange %v with routing key %v\n", parsed, constants.MessageEventsExchange, constants.MessageSentKey)

	conv, err := srv.GetConversation(parsed.Sender, parsed.Receiver)
	if err != nil {
		log.Printf("error fetching conversation: %v\n", err)
		return err
	}
	if conv == nil {
		conv, err = srv.CreateConversation(parsed.Sender, parsed.Receiver)
		if err != nil {
			log.Printf("error creating new conversation %v\n", err)
			return err
		}
		log.Printf("Created conversation: %v", conv)
	}

	// Create a new message
	message := &models.Message{
		ID:             parsed.ID,
		Content:        parsed.Content,
		Sender:         parsed.Sender,
		Type:           parsed.Type,
		ConversationID: conv.ID,
		Receiver:       parsed.Receiver,
		Read:           parsed.Read,
		Status:         parsed.Status,
		CreatedAt: parsed.CreatedAt,
		UpdatedAt: parsed.UpdatedAt,
	}

	// Add or update the message
	if parsed.Type == constants.MessageCreate {
		err = srv.AddMessage(message)
	} else {
		err = fmt.Errorf("error processing: expected message type to be message.create, instead was: %v", parsed.Type)
		log.Printf("%v\n", err)
		return err
	}
	if err != nil {
		log.Printf("error inserting message to db: %v\n", err)
		return err
	}

	log.Printf("messages service added message %v", message)
	return nil
}
