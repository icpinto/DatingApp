package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/models"
)

// ProfileSyncOutboxRepository manages profile synchronization outbox events.
type ProfileSyncOutboxRepository struct {
	db *sql.DB
}

// NewProfileSyncOutboxRepository creates a new ProfileSyncOutboxRepository.
func NewProfileSyncOutboxRepository(db *sql.DB) *ProfileSyncOutboxRepository {
	return &ProfileSyncOutboxRepository{db: db}
}

// Enqueue stores a new profile synchronization event.
func (r *ProfileSyncOutboxRepository) Enqueue(userID int) error {
	eventID := uuid.New().String()
	_, err := r.db.Exec(`
        INSERT INTO profile_sync_outbox (event_id, user_id, processed, created_at)
        VALUES ($1, $2, false, $3)`, eventID, userID, time.Now())
	if err != nil {
		log.Printf("ProfileSyncOutboxRepository.Enqueue exec error for user %d: %v", userID, err)
	}
	return err
}

// FetchPending retrieves pending profile synchronization events.
func (r *ProfileSyncOutboxRepository) FetchPending(limit int) ([]models.ProfileSyncOutbox, error) {
	rows, err := r.db.Query(`
        SELECT event_id, user_id
        FROM profile_sync_outbox
        WHERE processed = false
        ORDER BY created_at
        LIMIT $1`, limit)
	if err != nil {
		log.Printf("ProfileSyncOutboxRepository.FetchPending query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var events []models.ProfileSyncOutbox
	for rows.Next() {
		var event models.ProfileSyncOutbox
		if err := rows.Scan(&event.EventID, &event.UserID); err != nil {
			log.Printf("ProfileSyncOutboxRepository.FetchPending scan error: %v", err)
			return nil, err
		}
		events = append(events, event)
	}
	if err := rows.Err(); err != nil {
		log.Printf("ProfileSyncOutboxRepository.FetchPending rows error: %v", err)
		return nil, err
	}
	return events, nil
}

// MarkProcessed marks a profile synchronization event as processed.
func (r *ProfileSyncOutboxRepository) MarkProcessed(eventID string) error {
	_, err := r.db.Exec(`
        UPDATE profile_sync_outbox
        SET processed = true
        WHERE event_id = $1`, eventID)
	if err != nil {
		log.Printf("ProfileSyncOutboxRepository.MarkProcessed exec error for event %s: %v", eventID, err)
	}
	return err
}
