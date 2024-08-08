package services

import (
	"context"
	"errors"
	"log"

	"github.com/lib/pq"
	"github.com/yonraz/gochat_messages/models"
	"gorm.io/gorm"
)
var MESSAGE_PAGINATION_SIZE = 20

type MessagesServiceInterface interface {
    GetConversation(sender, receiver string) (*models.Conversation, error)
	GetConversationWithMessages(sender, receiver string, page int) (*models.Conversation, error)
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

func (srv *MessagesService) GetConversation(sender, receiver string) (*models.Conversation, error) {
var conv models.Conversation
	participants := pq.StringArray{sender, receiver}
	query := srv.DB.WithContext(context.Background()).
		Where("participants @> ? AND participants <@ ?", participants, participants).
		First(&conv)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		conversation := models.Conversation{
			Participants: participants,
			Messages:     []models.Message{},
		}
		
		result := srv.DB.Create(&conversation)
		if result.Error != nil {
			log.Printf("error saving new conversation: %v\n", result.Error)
			return nil, result.Error
		}
		return &conversation, nil 
	} else if query.Error != nil {
		log.Printf("error querying conversation: %v\n", query.Error)
		return nil, query.Error
	}

	return &conv, nil 
}

func (srv *MessagesService) GetConversationWithMessages(sender, receiver string, page int) (*models.Conversation, error) {
	var conv models.Conversation
	participants := pq.StringArray{sender, receiver}
	offset := (page-1) * MESSAGE_PAGINATION_SIZE
	query := srv.DB.WithContext(context.Background()).
		Where("participants @> ? AND participants <@ ?", participants, participants).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at ASC").Offset(offset).Limit(MESSAGE_PAGINATION_SIZE)
		}).
		First(&conv)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		conversation := models.Conversation{
			Participants: participants,
			Messages:     []models.Message{},
		}
		
		result := srv.DB.Create(&conversation)
		if result.Error != nil {
			log.Printf("error saving new conversation: %v\n", result.Error)
			return nil, result.Error
		}
		return &conversation, nil 
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

    // Update only the fields that need to be updated
    updateFields := map[string]interface{}{
        "content":      message.Content,
        "read":         message.Read,
        "status":       message.Status,
        "type":         message.Type,
        "created_at":   message.CreatedAt,
        "version":      existingMessage.Version + 1,
    }
    
    // Save the updated message
    result = s.DB.Model(&existingMessage).Updates(updateFields)
	if result.Error != nil {
		return nil, errors.New("failed to update message")
	}
    
    return &existingMessage, nil
}

func (srv *MessagesService) CreateConversation(sender, receiver string) (*models.Conversation, error) {
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


