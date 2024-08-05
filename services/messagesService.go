package services

import (
	"context"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/yonraz/gochat_messages/models"
	"gorm.io/gorm"
)

type MessagesServiceInterface interface {
    GetConversation(sender, receiver string) (*models.Conversation, error)
    AddMessage(msg *models.Message) error
	UpdateMessage(msg *models.Message) error
	CreateConversation(sender string, receiver string) (*models.Conversation, error)
}

type MessagesService struct {
	DB *gorm.DB
}

func NewMessagesService(db *gorm.DB) *MessagesService {
	return &MessagesService{
		DB: db,
	}
}

func (srv *MessagesService) GetConversation(sender string, receiver string) (*models.Conversation, error) {
	var conv models.Conversation
	participants := pq.StringArray{sender, receiver}
	query := srv.DB.WithContext(context.Background()).
		Where("participants @> ? AND participants <@ ?", participants, participants).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC")
		}).
		First(&conv)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return nil, nil // just return empty, no need for error
	} else if query.Error != nil {
		log.Printf("error querying conversation: %v\n", query.Error)
		return nil, query.Error
	}

	return &conv, nil 
}

func (srv *MessagesService) AddMessage(msg *models.Message) error {
	err := srv.DB.Create(&msg).Error
	if err != nil {
		return err
	}

	return nil
}

func (srv *MessagesService) UpdateMessage(msg *models.Message) error {
	err := srv.DB.Model(&models.Message{}).Where("id = ?", msg.ID).Updates(msg).Error
	if err != nil {
		return err
	}

	return nil
}

func (srv *MessagesService) CreateConversation(sender string, receiver string) (*models.Conversation, error) {
	conv := &models.Conversation{
		Participants: pq.StringArray{sender, receiver},
		Messages: []models.Message{},
	}

	result := srv.DB.Model(&models.Conversation{}).Create(&conv)
	if result.Error != nil {
		return nil, result.Error
	}

	return conv, nil
}


