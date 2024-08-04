package cache

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/yonraz/gochat_messages/models"
)

type Redis struct {
	client *redis.Client
	isConnected bool
}

func NewClient() (*Redis, error) {
	host := os.Getenv("REDIS_HOST")
	password := os.Getenv("REDIS_PASSWORD")


	client := redis.NewClient(&redis.Options{
	Addr: host,
	Password: password,
	DB: 0,
	})
	var isConnected bool = false
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Printf("could not connect to redis: %v\n", err)

	} else {
		isConnected = true
		log.Printf("Connected to redis!\n")
	}
		
	return &Redis{
		client,
		isConnected,
	}, err
}

func (r *Redis) GetQuery(query string) (*models.Conversation, error) {
	if !r.isConnected {
		var err error = errors.New("redis is not connected") 
		return nil, err
	}
	data, err := r.client.Get(context.Background(), "queries:users:"+query).Result()
	if err != nil {
		log.Printf("error retreiving query %v from redis: %v\n",query, err)
		return nil, err
	}
	var conversation *models.Conversation
	err = json.Unmarshal([]byte(data), &conversation)
	if err != nil {
		log.Printf("error unmarshalling data: %v\n", err)
		return nil, err
	}

	return conversation, err
}

func (r *Redis) SetQuery(query string, conversation *models.Conversation) error {
	if !r.isConnected {
		return errors.New("redis is not connected")
	}
	ctx := context.Background()
	data, err := json.Marshal(conversation)
	if err != nil {
		log.Printf("error marshaling data: %v\n", err)
		return err
	}

	err = r.client.Set(ctx, "queries:users:"+query, data, 0).Err()
	if err != nil {
		log.Printf("error in SetQuery function: %v\n", err)
		return err
	}

	return nil
}