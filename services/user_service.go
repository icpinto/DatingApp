package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/utils"
)

// UserService provides user related operations.
type UserService struct {
	db                  *sql.DB
	lifecycleOutboxRepo *repositories.UserLifecycleOutboxRepository
}

// NewUserService creates a new UserService with the given database handle.
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db, lifecycleOutboxRepo: repositories.NewUserLifecycleOutboxRepository(db)}
}

// GetUsepwd retrieves the hashed password and user ID for a username.
func (s *UserService) GetUsepwd(username string) (string, int, error) {
	pwd, id, err := repositories.GetUserpwdByUsername(s.db, username)
	if err != nil {
		log.Printf("GetUsepwd service error for %s: %v", username, err)
	}
	return pwd, id, err
}

// RegisterUser creates a new user in the database after hashing the password.
func (s *UserService) RegisterUser(user models.User) error {
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Printf("RegisterUser hash password error for %s: %v", user.Username, err)
		return err
	}
	user.Password = hashedPassword
	if err := repositories.CreateUser(s.db, user); err != nil {
		log.Printf("RegisterUser repository error for %s: %v", user.Username, err)
		return err
	}
	return nil
}

// GetUserIDByUsername returns the user ID for a given username.
func (s *UserService) GetUserIDByUsername(username string) (int, error) {
	id, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetUserIDByUsername service error for %s: %v", username, err)
	}
	return id, err
}

// GetUsernameByID returns the username for a given user ID.
func (s *UserService) GetUsernameByID(userID int) (string, error) {
	username, err := repositories.GetUsernameByID(s.db, userID)
	if err != nil {
		log.Printf("GetUsernameByID service error for %d: %v", userID, err)
	}
	return username, err
}

// GetUsernameByIDAllowInactive returns the username for the given user ID regardless of account status.
func (s *UserService) GetUsernameByIDAllowInactive(userID int) (string, error) {
	username, err := repositories.GetUsernameByIDAllowInactive(s.db, userID)
	if err != nil {
		log.Printf("GetUsernameByIDAllowInactive service error for %d: %v", userID, err)
	}
	return username, err
}

// GetUserStatus returns whether the provided user account is active.
func (s *UserService) GetUserStatus(userID int) (bool, error) {
	isActive, err := repositories.GetUserStatusByID(s.db, userID)
	if err != nil {
		log.Printf("GetUserStatus service error for %d: %v", userID, err)
	}
	return isActive, err
}

// DeactivateUser marks the account inactive and enqueues a lifecycle event for downstream cleanup.
func (s *UserService) DeactivateUser(ctx context.Context, userID int, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := repositories.DeactivateUserTx(tx, userID); err != nil {
		return err
	}

	event, err := s.buildLifecycleEvent(userID, models.UserLifecycleEventTypeDeactivated, reason)
	if err != nil {
		return err
	}
	if err := s.lifecycleOutboxRepo.EnqueueTx(tx, event); err != nil {
		return err
	}

	return tx.Commit()
}

// ReactivateUser marks the account active again and enqueues a lifecycle event for downstream restoration.
func (s *UserService) ReactivateUser(ctx context.Context, userID int, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := repositories.ReactivateUserTx(tx, userID); err != nil {
		return err
	}

	event, err := s.buildLifecycleEvent(userID, models.UserLifecycleEventTypeReactivated, reason)
	if err != nil {
		return err
	}
	if err := s.lifecycleOutboxRepo.EnqueueTx(tx, event); err != nil {
		return err
	}

	return tx.Commit()
}

// DeleteUser removes the account and enqueues a deletion lifecycle event.
func (s *UserService) DeleteUser(ctx context.Context, userID int, reason string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	event, err := s.buildLifecycleEvent(userID, models.UserLifecycleEventTypeDeleted, reason)
	if err != nil {
		return err
	}
	if err := s.lifecycleOutboxRepo.EnqueueTx(tx, event); err != nil {
		return err
	}

	if err := repositories.DeleteUserTx(tx, userID); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *UserService) buildLifecycleEvent(userID int, eventType models.UserLifecycleEventType, reason string) (models.UserLifecycleOutbox, error) {
	payload := make(map[string]string)
	trimmed := strings.TrimSpace(reason)
	if trimmed != "" {
		payload["reason"] = trimmed
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return models.UserLifecycleOutbox{}, err
	}
	return models.UserLifecycleOutbox{
		EventID:   uuid.NewString(),
		UserID:    userID,
		EventType: eventType,
		Payload:   body,
		CreatedAt: time.Now().UTC(),
	}, nil
}
