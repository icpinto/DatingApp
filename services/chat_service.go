package services

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func CreateConversation(db *sql.DB, user1ID, user2ID int) error {
	return repositories.CreateConversation(db, user1ID, user2ID)
}

func GetAllConversations(db *sql.DB, username string) ([]models.Conversation, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return nil, err
	}
	return repositories.GetConversationsByUserID(db, userID)
}

func GetChatHistory(db *sql.DB, conversationID string) ([]models.ChatMessage, error) {
	return repositories.GetMessagesByConversationID(db, conversationID)
}
