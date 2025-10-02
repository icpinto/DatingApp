package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/icpinto/dating-app/models"
)

// UserLifecycleOutboxRepository manages lifecycle events queued for downstream services.
type UserLifecycleOutboxRepository struct {
	db *sql.DB
}

// NewUserLifecycleOutboxRepository creates a repository backed by the provided database handle.
func NewUserLifecycleOutboxRepository(db *sql.DB) *UserLifecycleOutboxRepository {
	return &UserLifecycleOutboxRepository{db: db}
}

// EnqueueTx inserts a lifecycle event inside the supplied transaction to keep it atomic with user mutations.
func (r *UserLifecycleOutboxRepository) EnqueueTx(tx *sql.Tx, event models.UserLifecycleOutbox) error {
	_, err := tx.Exec(`
        INSERT INTO user_lifecycle_outbox (event_id, user_id, event_type, payload, processed, created_at)
        VALUES ($1, $2, $3, $4, false, $5)`,
		event.EventID, event.UserID, event.EventType, event.Payload, event.CreatedAt)
	if err != nil {
		log.Printf("UserLifecycleOutboxRepository.EnqueueTx exec error for user %d: %v", event.UserID, err)
	}
	return err
}

// FetchPending returns at most limit unprocessed lifecycle events ordered by creation time.
func (r *UserLifecycleOutboxRepository) FetchPending(limit int) ([]models.UserLifecycleOutbox, error) {
	rows, err := r.db.Query(`
        SELECT event_id, user_id, event_type, payload, processed, processed_at, created_at
        FROM user_lifecycle_outbox
        WHERE processed = false
        ORDER BY created_at
        LIMIT $1`, limit)
	if err != nil {
		log.Printf("UserLifecycleOutboxRepository.FetchPending query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.UserLifecycleOutbox
	for rows.Next() {
		var event models.UserLifecycleOutbox
		var processedAt sql.NullTime
		if err := rows.Scan(&event.EventID, &event.UserID, &event.EventType, &event.Payload, &event.Processed, &processedAt, &event.CreatedAt); err != nil {
			log.Printf("UserLifecycleOutboxRepository.FetchPending scan error: %v", err)
			return nil, err
		}
		if processedAt.Valid {
			t := processedAt.Time
			event.ProcessedAt = &t
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		log.Printf("UserLifecycleOutboxRepository.FetchPending rows error: %v", err)
		return nil, err
	}
	return events, nil
}

// MarkProcessed flags an event as processed and records when it was delivered.
func (r *UserLifecycleOutboxRepository) MarkProcessed(eventID string) error {
	_, err := r.db.Exec(`
        UPDATE user_lifecycle_outbox
        SET processed = true, processed_at = $1
        WHERE event_id = $2`, time.Now(), eventID)
	if err != nil {
		log.Printf("UserLifecycleOutboxRepository.MarkProcessed exec error for event %s: %v", eventID, err)
	}
	return err
}
