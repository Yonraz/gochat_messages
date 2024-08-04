package consumers

import (
	"encoding/json"
	"errors"
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
	var parsed *models.Message
	
	json.Unmarshal(msg.Body, &parsed)
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
	}

	message := &models.Message{
		ID: parsed.ID,
		Content: fmt.Sprintf("%v sent a message", parsed.Sender),
		Sender: parsed.Sender,
		Type: parsed.Type,
		ConversationID: conv.ID,
		Receiver: parsed.Receiver,
		Read: parsed.Read,
		Status: parsed.Status,
	}
	if parsed.Type == constants.MessageCreate {
		err = srv.AddMessage(message)
	} else if parsed.Type == constants.MessageUpdate {
		err = srv.UpdateMessage(message)
	} else {
		return errors.New("no valid message type was published")
	}
	if err != nil {
		log.Printf("error inserting message to db: %v\n", err)
		return err
	}

	log.Printf("messages service added message %v", message)
	return nil
}