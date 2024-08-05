package initializers

import "github.com/yonraz/gochat_messages/models"

func SyncDatabase() {
	
	DB.AutoMigrate(&models.Conversation{})
	
	DB.AutoMigrate(&models.Message{})
}
