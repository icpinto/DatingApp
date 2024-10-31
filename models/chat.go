package models

import (
	"time"

	"github.com/gorilla/websocket"
)

type ChatMessage struct {
	SenderID       int    `json:"sender_id"`
	ReceiverID     int    `json:"receiver_id"`
	ConversationID int    `json:"conversation_id"`
	Message        string `json:"message"`
	CreatedAt      time.Time
}

type Client struct {
	Conn   *websocket.Conn
	UserID int
}

type Conversation struct {
	ID        int       `json:"id"`
	User1ID   int       `json:"user1_id"`
	User2ID   int       `json:"user2_id"`
	CreatedAt time.Time `json:"created_at"`
}
