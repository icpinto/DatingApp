package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/utils"
)

// OutboxWorker processes conversation outbox events.
type OutboxWorker struct {
	db                  *sql.DB
	outboxRepo          *repositories.OutboxRepository
	profileOutboxRepo   *repositories.ProfileSyncOutboxRepository
	profileRepo         *repositories.ProfileRepository
	frRepo              *repositories.FriendRequestRepository
	client              *http.Client
	baseURL             string
	matchService        *MatchService
	lifecycleOutboxRepo *repositories.UserLifecycleOutboxRepository
	lifecyclePublisher  *RabbitMQPublisher
}

// NewOutboxWorker creates a new OutboxWorker.
func NewOutboxWorker(db *sql.DB, baseURL string, matchService *MatchService, lifecyclePublisher *RabbitMQPublisher) *OutboxWorker {
	return &OutboxWorker{
		db:                  db,
		outboxRepo:          repositories.NewOutboxRepository(db),
		profileOutboxRepo:   repositories.NewProfileSyncOutboxRepository(db),
		profileRepo:         repositories.NewProfileRepository(db),
		frRepo:              repositories.NewFriendRequestRepository(db),
		client:              &http.Client{Timeout: 5 * time.Second},
		baseURL:             baseURL,
		matchService:        matchService,
		lifecycleOutboxRepo: repositories.NewUserLifecycleOutboxRepository(db),
		lifecyclePublisher:  lifecyclePublisher,
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
	if err := w.processConversationEvents(); err != nil {
		return err
	}
	if err := w.processProfileSyncEvents(); err != nil {
		return err
	}
	if err := w.processUserLifecycleEvents(); err != nil {
		return err
	}
	return nil
}

func (w *OutboxWorker) processConversationEvents() error {
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

func (w *OutboxWorker) processProfileSyncEvents() error {
	if w.profileOutboxRepo == nil || w.matchService == nil {
		return nil
	}
	events, err := w.profileOutboxRepo.FetchPending(10)
	if err != nil {
		return err
	}
	for _, event := range events {
		profile, err := w.profileRepo.GetByUserID(event.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("OutboxWorker profile not found for user %d, marking event %s processed", event.UserID, event.EventID)
				if markErr := w.profileOutboxRepo.MarkProcessed(event.EventID); markErr != nil {
					log.Printf("OutboxWorker mark processed error for event %s: %v", event.EventID, markErr)
				}
				continue
			}
			log.Printf("OutboxWorker profile fetch error for user %d: %v", event.UserID, err)
			continue
		}
		if _, err := w.matchService.UpsertProfile(context.Background(), profile.Profile); err != nil {
			log.Printf("OutboxWorker profile sync error for user %d: %v", event.UserID, err)
			continue
		}
		if err := w.profileOutboxRepo.MarkProcessed(event.EventID); err != nil {
			log.Printf("OutboxWorker mark processed error for event %s: %v", event.EventID, err)
			continue
		}
	}
	return nil
}

func (w *OutboxWorker) processUserLifecycleEvents() error {
	if w.lifecycleOutboxRepo == nil || w.lifecyclePublisher == nil {
		return nil
	}
	events, err := w.lifecycleOutboxRepo.FetchPending(25)
	if err != nil {
		return err
	}
	for _, event := range events {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		publishErr := w.lifecyclePublisher.PublishLifecycleEvent(ctx, event)
		cancel()
		if publishErr != nil {
			log.Printf("OutboxWorker lifecycle publish error for event %s: %v", event.EventID, publishErr)
			continue
		}
		if err := w.lifecycleOutboxRepo.MarkProcessed(event.EventID); err != nil {
			log.Printf("OutboxWorker lifecycle mark processed error for event %s: %v", event.EventID, err)
			continue
		}
	}
	return nil
}

func (w *OutboxWorker) handleEvent(eventID string, user1ID, user2ID int) error {
	token, err := utils.GenerateToken(user1ID)
	if err != nil {
		return err
	}
	payload := map[string][]string{
		"participant_ids": []string{strconv.Itoa(user2ID)},
	}
	body, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/conversations", w.baseURL), bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", eventID)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	var res struct {
		ConversationID uuid.UUID `json:"conversation_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return err
	}
	if err := w.frRepo.LinkConversation(user1ID, user2ID, res.ConversationID); err != nil {
		return err
	}
	return w.outboxRepo.MarkProcessed(eventID, res.ConversationID)
}
