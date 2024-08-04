package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/streadway/amqp"
	"github.com/yonraz/gochat_messages/initializers"
)

func init () {
	fmt.Println("Application starting...")
	time.Sleep(1 * time.Minute)
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
	initializers.ConnectToRabbitmq()
	initializers.ConnectToRedis()
}

func main() {
	router := gin.Default()

	defer func() {
		if err := initializers.RmqChannel.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}()
	defer func() {
		if err := initializers.RmqConn.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ connection: %v", err)
		}
	}()

	// dev - insert 35 lines of mock user data
	// services.CreateMockUsers()
	//

	go startConsumers()


	router.Run()
}

func startConsumers() {
	// targetConsumers := []func(*amqp.Channel) *consumers.Consumer{
	// 	consumers.NewUserRegisteredConsumer,
	// 	consumers.NewUserLoggedinConsumer,
	// 	consumers.NewUserSignedoutConsumer,
	// }

	// for _, consumerFunc := range targetConsumers {
	// 	consumer := consumerFunc(initializers.RmqChannel)
	// 	go func(c *consumers.Consumer) {
	// 		if err := c.Consume(); err != nil {
	// 			log.Fatalf("Error starting consumer: %v", err)
	// 		}
	// 	}(consumer)
	// }
}