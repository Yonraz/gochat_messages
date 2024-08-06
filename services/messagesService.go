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
	CreateConversation(sender string, receiver string) (*models.Conversation, error)
	UpdateMessage(message *models.Message) (*models.Message, error)
	GetMessageByID(id string) (*models.Message, error)
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
		conversation := models.Conversation{
			Participants: participants,
			Messages: []models.Message{},
		}
		srv.DB.Save(&conversation)
		return &conversation, nil // just return empty, no need for error
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

func (s *MessagesService) UpdateMessage(message *models.Message) (*models.Message, error) {
    var existingMessage models.Message
    
    // Retrieve the existing message by ID
    result := s.DB.First(&existingMessage, "id = ?", message.ID)
    if result.Error != nil {
        if result.Error == gorm.ErrRecordNotFound {
            // Handle case where the record is not found
            return nil, nil
        }
        return nil, result.Error
    }

    // Check version
    if existingMessage.Version != message.Version {
        // Return an error if versions do not match
        return nil, errors.New("version conflict")
    }

    // Update the message
    existingMessage.Content = message.Content
    // Add other fields if needed
    existingMessage.Version++ // Increment the version number
    
    // Save the updated message
    result = s.DB.Save(&existingMessage)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &existingMessage, nil
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

func (srv *MessagesService) GetMessageByID(id string) (*models.Message, error) {
	var msg *models.Message 
	err := srv.DB.Where(models.Message{ID: id}).First(&msg).Error

	if err != nil {
		return nil, err
	}

	return msg, nil
}


