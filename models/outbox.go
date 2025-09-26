package models

import (
	"time"

	"github.com/google/uuid"
)

// ConversationOutbox represents an event for creating a conversation.
type ConversationOutbox struct {
	EventID        string     `json:"event_id"`
	User1ID        int        `json:"user1_id"`
	User2ID        int        `json:"user2_id"`
	ConversationID *uuid.UUID `json:"conversation_id,omitempty"`
	Processed      bool       `json:"processed"`
	CreatedAt      time.Time  `json:"created_at"`
}

// ProfileSyncOutbox represents a pending profile synchronization event.
type ProfileSyncOutbox struct {
	EventID   string    `json:"event_id"`
	UserID    int       `json:"user_id"`
	Processed bool      `json:"processed"`
	CreatedAt time.Time `json:"created_at"`
}
