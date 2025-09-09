package repositories

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
)

// ConversationRepository manages conversations and messages.
type ConversationRepository struct {
	db *sql.DB
}

// NewConversationRepository creates a new ConversationRepository.
func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// Create starts a new conversation between two users.
func (r *ConversationRepository) Create(user1ID, user2ID int) error {
	_, err := r.db.Exec(`INSERT INTO conversations (user1_id, user2_id, created_at) VALUES ($1, $2, NOW())`, user1ID, user2ID)
	if err != nil {
		log.Printf("ConversationRepository.Create exec error for users %d and %d: %v", user1ID, user2ID, err)
	}
	return err
}

// GetByUserID retrieves all conversations for a user.
func (r *ConversationRepository) GetByUserID(userID int) ([]models.Conversation, error) {
	rows, err := r.db.Query(`
                SELECT id, user1_id, user2_id, created_at
                FROM conversations
                WHERE user1_id = $1 OR user2_id = $1;`, userID)
	if err != nil {
		log.Printf("ConversationRepository.GetByUserID query error for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(&c.ID, &c.User1ID, &c.User2ID, &c.CreatedAt); err != nil {
			log.Printf("ConversationRepository.GetByUserID scan error: %v", err)
			return nil, err
		}
		conversations = append(conversations, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ConversationRepository.GetByUserID rows error: %v", err)
		return nil, err
	}
	return conversations, nil
}

// GetMessages retrieves messages for a conversation.
func (r *ConversationRepository) GetMessages(conversationID string) ([]models.ChatMessage, error) {
	rows, err := r.db.Query(`SELECT conversation_id, sender_id, message, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		log.Printf("ConversationRepository.GetMessages query error for conversation %s: %v", conversationID, err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.ConversationID, &m.SenderID, &m.Message, &m.CreatedAt); err != nil {
			log.Printf("ConversationRepository.GetMessages scan error: %v", err)
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ConversationRepository.GetMessages rows error: %v", err)
		return nil, err
	}
	return messages, nil
}

// SaveMessage stores a chat message.
func (r *ConversationRepository) SaveMessage(msg models.ChatMessage) error {
	_, err := r.db.Exec(`
        INSERT INTO messages (conversation_id, sender_id, message, created_at)
        VALUES ($1, $2, $3, NOW())`, msg.ConversationID, msg.SenderID, msg.Message)
	if err != nil {
		log.Printf("ConversationRepository.SaveMessage exec error for conversation %d: %v", msg.ConversationID, err)
	}
	return err
}

// Delete removes a conversation and its messages.
func (r *ConversationRepository) Delete(id int) error {
	if _, err := r.db.Exec(`DELETE FROM conversations WHERE id = $1`, id); err != nil {
		log.Printf("ConversationRepository.Delete exec error for conversation %d: %v", id, err)
		return err
	}
	return nil
}
