package services

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// ChatService provides chat-related operations.
type ChatService struct {
	db   *sql.DB
	repo *repositories.ConversationRepository
}

// NewChatService creates a new ChatService.
func NewChatService(db *sql.DB) *ChatService {
	return &ChatService{db: db, repo: repositories.NewConversationRepository(db)}
}

// CreateConversation creates a conversation between two users.
func (s *ChatService) CreateConversation(user1ID, user2ID int) error {
	if err := s.repo.Create(user1ID, user2ID); err != nil {
		log.Printf("CreateConversation service error for users %d and %d: %v", user1ID, user2ID, err)
		return err
	}
	return nil
}

// GetAllConversations retrieves all conversations for a user.
func (s *ChatService) GetAllConversations(username string) ([]models.Conversation, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetAllConversations user lookup error for %s: %v", username, err)
		return nil, err
	}
	conversations, err := s.repo.GetByUserID(userID)
	if err != nil {
		log.Printf("GetAllConversations fetch error for user %d: %v", userID, err)
		return nil, err
	}
	return conversations, nil
}

// GetChatHistory retrieves messages for a conversation.
func (s *ChatService) GetChatHistory(conversationID string) ([]models.ChatMessage, error) {
	messages, err := s.repo.GetMessages(conversationID)
	if err != nil {
		log.Printf("GetChatHistory service error for conversation %s: %v", conversationID, err)
		return nil, err
	}
	return messages, nil
}

// SaveMessage stores a chat message in the database.
func (s *ChatService) SaveMessage(msg models.ChatMessage) error {
	if err := s.repo.SaveMessage(msg); err != nil {
		log.Printf("SaveMessage service error for conversation %d: %v", msg.ConversationID, err)
		return err
	}
	return nil
}
