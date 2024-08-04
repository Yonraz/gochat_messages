package consumers

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_messages/constants"
	"github.com/yonraz/gochat_messages/initializers"
	"github.com/yonraz/gochat_messages/services"
)

type Consumer struct {
	channel     *amqp.Channel
	srv 		*services.MessagesService
	queueName   string
	routingKey  string
	exchange    string
	handlerFunc func(*services.MessagesService, amqp.Delivery) error
}

func NewConsumer(channel *amqp.Channel, queueName constants.Queues, routingKey constants.RoutingKey, exchange constants.Exchange, handlerFunc func(*services.MessagesService, amqp.Delivery) error) *Consumer {
	return &Consumer {
		channel: channel,
		srv: services.NewMessagesService(initializers.DB),
		queueName: string(queueName),
		routingKey: string(routingKey),
		exchange: string(exchange),
		handlerFunc: handlerFunc,
	}
}

func (c *Consumer) Consume() error {
	msgs, err := c.channel.Consume(
		c.queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming %w", err)
	}

	go func () {
		for msg := range msgs {
			if err := c.handlerFunc(c.srv, msg); err != nil {
				fmt.Printf("error consuming message %v: %v\n", msg, err)
				msg.Nack(false, true)		
			} else {

			msg.Ack(false)
			}
		} 
	}()

	fmt.Printf("Started consuming on queue: %s\n", c.queueName)
	return nil
}