package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/models"
)

// OutboxRepository manages conversation outbox events.
type OutboxRepository struct {
	db *sql.DB
}

// NewOutboxRepository creates a new OutboxRepository.
func NewOutboxRepository(db *sql.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

// CreateTx inserts a new outbox event within the given transaction.
func (r *OutboxRepository) CreateTx(tx *sql.Tx, event models.ConversationOutbox) error {
	_, err := tx.Exec(`
        INSERT INTO conversation_outbox (event_id, user1_id, user2_id, processed, created_at)
        VALUES ($1, $2, $3, false, $4)`, event.EventID, event.User1ID, event.User2ID, time.Now())
	if err != nil {
		log.Printf("OutboxRepository.CreateTx exec error for users %d and %d: %v", event.User1ID, event.User2ID, err)
	}
	return err
}

// FetchPending retrieves unprocessed outbox events.
func (r *OutboxRepository) FetchPending(limit int) ([]models.ConversationOutbox, error) {
	rows, err := r.db.Query(`
        SELECT event_id, user1_id, user2_id
        FROM conversation_outbox
        WHERE processed = false
        ORDER BY created_at
        LIMIT $1`, limit)
	if err != nil {
		log.Printf("OutboxRepository.FetchPending query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.ConversationOutbox
	for rows.Next() {
		var e models.ConversationOutbox
		if err := rows.Scan(&e.EventID, &e.User1ID, &e.User2ID); err != nil {
			log.Printf("OutboxRepository.FetchPending scan error: %v", err)
			return nil, err
		}
		events = append(events, e)
	}
	if err := rows.Err(); err != nil {
		log.Printf("OutboxRepository.FetchPending rows error: %v", err)
		return nil, err
	}
	return events, nil
}

// MarkProcessed marks an outbox event as processed and stores the conversation ID.
func (r *OutboxRepository) MarkProcessed(eventID string, conversationID uuid.UUID) error {
	_, err := r.db.Exec(`
        UPDATE conversation_outbox
        SET processed = true, conversation_id = $1
        WHERE event_id = $2`, conversationID, eventID)
	if err != nil {
		log.Printf("OutboxRepository.MarkProcessed exec error for event %s: %v", eventID, err)
	}
	return err
}
