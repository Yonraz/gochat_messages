package initializers

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
	"github.com/yonraz/gochat_messages/events/utils"
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
	err = utils.DeclareExchanges(RmqChannel)
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