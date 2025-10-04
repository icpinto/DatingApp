package models

import (
	"encoding/json"
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

// UserLifecycleEventType enumerates supported lifecycle transitions broadcast to downstream services.
type UserLifecycleEventType string

const (
	// UserLifecycleEventTypeDeactivated indicates that the account has been deactivated but not removed.
	UserLifecycleEventTypeDeactivated UserLifecycleEventType = "deactivated"
	// UserLifecycleEventTypeReactivated indicates that the account has been reactivated and should be restored downstream.
	UserLifecycleEventTypeReactivated UserLifecycleEventType = "reactivated"
	// UserLifecycleEventTypeDeleted indicates that the account and related data have been removed from the core service.
	UserLifecycleEventTypeDeleted UserLifecycleEventType = "deleted"
)

// UserLifecycleOutbox represents lifecycle events (deactivation/deletion) queued for downstream processing.
type UserLifecycleOutbox struct {
	EventID     string                 `json:"event_id"`
	UserID      int                    `json:"user_id"`
	EventType   UserLifecycleEventType `json:"event_type"`
	Payload     json.RawMessage        `json:"payload"`
	Processed   bool                   `json:"processed"`
	ProcessedAt *time.Time             `json:"processed_at,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}
