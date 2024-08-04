package initializers

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_messages/events/utils"
	"github.com/yonraz/gochat_messages/constants"
)

var RmqChannel *amqp.Channel
var RmqConn *amqp.Connection

func ConnectToRabbitmq() {
	var err error
	user := os.Getenv("RMQ_USER")
	password := os.Getenv("RMQ_PASSWORD")
	connectionString := fmt.Sprintf("amqp://%v:%v@rabbitmq:5672/", user, password)
	RmqConn, err = amqp.Dial(connectionString)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("Connected to Rabbitmq")

	RmqChannel, err = RmqConn.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// Declaring the topic exchange
	err = RmqChannel.ExchangeDeclare(
		string(constants.UserEventsExchange),
		"topic",             
		true,                
		false,               
		false, 
		false,              
		nil,                 
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	// create registration, login, logout queues
	err = utils.DeclareQueues(RmqChannel)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}