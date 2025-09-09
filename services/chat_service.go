package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func CreateConversation(db *sql.DB, user1ID, user2ID int) error {
	if err := repositories.CreateConversation(db, user1ID, user2ID); err != nil {
		log.Printf("CreateConversation service error for users %d and %d: %v", user1ID, user2ID, err)
		return err
	}
	return nil
}

func GetAllConversations(db *sql.DB, username string) ([]models.Conversation, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("GetAllConversations user lookup error for %s: %v", username, err)
		return nil, err
	}
	conversations, err := repositories.GetConversationsByUserID(db, userID)
	if err != nil {
		log.Printf("GetAllConversations fetch error for user %d: %v", userID, err)
		return nil, err
	}
	return conversations, nil
}

func GetChatHistory(db *sql.DB, conversationID string) ([]models.ChatMessage, error) {
	messages, err := repositories.GetMessagesByConversationID(db, conversationID)
	if err != nil {
		log.Printf("GetChatHistory service error for conversation %s: %v", conversationID, err)
		return nil, err
	}
	return messages, nil
}
