package repositories

import (
	"database/sql"

	"github.com/icpinto/dating-app/models"
)

func CreateConversation(db *sql.DB, user1ID, user2ID int) error {
	_, err := db.Exec(`INSERT INTO conversations (user1_id, user2_id, created_at) VALUES ($1, $2, NOW())`, user1ID, user2ID)
	return err
}

func GetConversationsByUserID(db *sql.DB, userID int) ([]models.Conversation, error) {
	rows, err := db.Query(`
                SELECT id, user1_id, user2_id, created_at
                FROM conversations
                WHERE user1_id = $1 OR user2_id = $1;`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []models.Conversation
	for rows.Next() {
		var c models.Conversation
		if err := rows.Scan(&c.ID, &c.User1ID, &c.User2ID, &c.CreatedAt); err != nil {
			return nil, err
		}
		conversations = append(conversations, c)
	}
	return conversations, rows.Err()
}

func GetMessagesByConversationID(db *sql.DB, conversationID string) ([]models.ChatMessage, error) {
	rows, err := db.Query(`SELECT conversation_id, sender_id, message, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage
	for rows.Next() {
		var m models.ChatMessage
		if err := rows.Scan(&m.ConversationID, &m.SenderID, &m.Message, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, rows.Err()
}
