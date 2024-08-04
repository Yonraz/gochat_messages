package initializers

import "github.com/yonraz/gochat_messages/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.Message{}, &models.Conversation{})
}