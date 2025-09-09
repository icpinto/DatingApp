package repositories

import (
	"database/sql"
	"log"

	"github.com/icpinto/dating-app/models"
)

func CreateConversation(db *sql.DB, user1ID, user2ID int) error {
	_, err := db.Exec(`INSERT INTO conversations (user1_id, user2_id, created_at) VALUES ($1, $2, NOW())`, user1ID, user2ID)
	if err != nil {
		log.Printf("CreateConversation exec error for users %d and %d: %v", user1ID, user2ID, err)
	}
	return err
}

func GetConversationsByUserID(db *sql.DB, userID int) ([]models.Conversation, error) {
	rows, err := db.Query(`
                SELECT id, user1_id, user2_id, created_at
                FROM conversations
                WHERE user1_id = $1 OR user2_id = $1;`, userID)
	if err != nil {
		log.Printf("GetConversationsByUserID query error for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(&c.ID, &c.User1ID, &c.User2ID, &c.CreatedAt); err != nil {
			log.Printf("GetConversationsByUserID scan error: %v", err)
			return nil, err
		}
		conversations = append(conversations, c)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetConversationsByUserID rows error: %v", err)
		return nil, err
	}
	return conversations, nil
}

func GetMessagesByConversationID(db *sql.DB, conversationID string) ([]models.ChatMessage, error) {
	rows, err := db.Query(`SELECT conversation_id, sender_id, message, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		log.Printf("GetMessagesByConversationID query error for conversation %s: %v", conversationID, err)
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.ConversationID, &m.SenderID, &m.Message, &m.CreatedAt); err != nil {
			log.Printf("GetMessagesByConversationID scan error: %v", err)
			return nil, err
		}
		messages = append(messages, m)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetMessagesByConversationID rows error: %v", err)
		return nil, err
	}
	return messages, nil
}
