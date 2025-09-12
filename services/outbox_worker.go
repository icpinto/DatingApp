package services

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/icpinto/dating-app/repositories"
)

// OutboxWorker processes conversation outbox events.
type OutboxWorker struct {
	db         *sql.DB
	outboxRepo *repositories.OutboxRepository
	frRepo     *repositories.FriendRequestRepository
	client     *http.Client
	baseURL    string
}

// NewOutboxWorker creates a new OutboxWorker.
func NewOutboxWorker(db *sql.DB, baseURL string) *OutboxWorker {
	return &OutboxWorker{
		db:         db,
		outboxRepo: repositories.NewOutboxRepository(db),
		frRepo:     repositories.NewFriendRequestRepository(db),
		client:     &http.Client{Timeout: 5 * time.Second},
		baseURL:    baseURL,
	}
}

// Start begins processing outbox events periodically.
func (w *OutboxWorker) Start() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		if err := w.process(); err != nil {
			log.Printf("OutboxWorker process error: %v", err)
		}
	}
}

func (w *OutboxWorker) process() error {
	events, err := w.outboxRepo.FetchPending(10)
	if err != nil {
		return err
	}
	for _, e := range events {
		if err := w.handleEvent(e.EventID, e.User1ID, e.User2ID); err != nil {
			log.Printf("OutboxWorker handle event %s error: %v", e.EventID, err)
			continue
		}
	}
	return nil
}

func (w *OutboxWorker) handleEvent(eventID string, user1ID, user2ID int) error {
	payload := map[string]int{"user1_id": user1ID, "user2_id": user2ID}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/internal/create_conversation", w.baseURL), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", eventID)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	var res struct {
		ConversationID int `json:"conversation_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if err := w.frRepo.LinkConversation(user1ID, user2ID, res.ConversationID); err != nil {
		return err
	}
	return w.outboxRepo.MarkProcessed(eventID, res.ConversationID)
}
