package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// ErrFriendRequestExists indicates a duplicate friend request.
var ErrFriendRequestExists = errors.New("friend request already exists")

// FriendRequestService provides operations related to friend requests.
type FriendRequestService struct {
	db         *sql.DB
	repo       *repositories.FriendRequestRepository
	outboxRepo *repositories.OutboxRepository
}

// NewFriendRequestService creates a new FriendRequestService.
func NewFriendRequestService(db *sql.DB) *FriendRequestService {
	return &FriendRequestService{db: db, repo: repositories.NewFriendRequestRepository(db), outboxRepo: repositories.NewOutboxRepository(db)}
}

// SendFriendRequest sends a friend request from a user to another.
func (s *FriendRequestService) SendFriendRequest(username string, request models.FriendRequest) error {
	senderID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("SendFriendRequest user lookup error for %s: %v", username, err)
		return err
	}
	request.SenderID = senderID
	request.SenderUsername = username
	request.Status = "pending"
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	receiverUsername, err := repositories.GetUsernameByID(s.db, request.ReceiverID)
	if err != nil {
		log.Printf("SendFriendRequest receiver lookup error for %d: %v", request.ReceiverID, err)
		return err
	}
	request.ReceiverUsername = receiverUsername

	_, err = s.repo.CheckExisting(request.SenderID, request.ReceiverID)
	if err == nil {
		log.Printf("SendFriendRequest duplicate for sender %d and receiver %d", request.SenderID, request.ReceiverID)
		return ErrFriendRequestExists
	}
	if err != sql.ErrNoRows {
		log.Printf("SendFriendRequest check existing error: %v", err)
		return err
	}
	if err := s.repo.Create(request); err != nil {
		log.Printf("SendFriendRequest insert error for sender %d and receiver %d: %v", request.SenderID, request.ReceiverID, err)
		return err
	}
	return nil
}

// AcceptFriendRequest accepts a pending friend request.
func (s *FriendRequestService) AcceptFriendRequest(requestID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("AcceptFriendRequest begin tx error for request %d: %v", requestID, err)
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`
                UPDATE friend_requests SET status = $1, updated_at = $2 WHERE id = $3`,
		"accepted", time.Now(), requestID); err != nil {
		log.Printf("AcceptFriendRequest update status error for request %d: %v", requestID, err)
		return err
	}

	var user1ID, user2ID int
	if err := tx.QueryRow(`SELECT sender_id, receiver_id FROM friend_requests WHERE id = $1`, requestID).Scan(&user1ID, &user2ID); err != nil {
		log.Printf("AcceptFriendRequest get users error for request %d: %v", requestID, err)
		return err
	}

	event := models.ConversationOutbox{
		EventID: uuid.New().String(),
		User1ID: user1ID,
		User2ID: user2ID,
	}
	if err := s.outboxRepo.CreateTx(tx, event); err != nil {
		log.Printf("AcceptFriendRequest create outbox error for request %d: %v", requestID, err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Printf("AcceptFriendRequest commit error for request %d: %v", requestID, err)
		return err
	}
	return nil
}

// RejectFriendRequest rejects a pending friend request.
func (s *FriendRequestService) RejectFriendRequest(requestID int) error {
	if err := s.repo.UpdateStatus(requestID, "rejected", time.Now()); err != nil {
		log.Printf("RejectFriendRequest update status error for request %d: %v", requestID, err)
		return err
	}
	return nil
}

// GetPendingRequests retrieves all pending friend requests for a user.
func (s *FriendRequestService) GetPendingRequests(username string) ([]models.FriendRequest, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetPendingRequests user lookup error for %s: %v", username, err)
		return nil, err
	}
	requests, err := s.repo.GetPending(userID)
	if err != nil {
		log.Printf("GetPendingRequests repository error for user %d: %v", userID, err)
		return nil, err
	}

	for i := range requests {
		sender, err := repositories.GetUsernameByID(s.db, requests[i].SenderID)
		if err != nil {
			log.Printf("GetPendingRequests sender lookup error for user %d: %v", requests[i].SenderID, err)
			return nil, err
		}
		receiver, err := repositories.GetUsernameByID(s.db, requests[i].ReceiverID)
		if err != nil {
			log.Printf("GetPendingRequests receiver lookup error for user %d: %v", requests[i].ReceiverID, err)
			return nil, err
		}
		requests[i].SenderUsername = sender
		requests[i].ReceiverUsername = receiver
	}

	return requests, nil
}

// GetSentRequests retrieves all friend requests sent by a user.
func (s *FriendRequestService) GetSentRequests(username string) ([]models.FriendRequest, error) {
	userID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("GetSentRequests user lookup error for %s: %v", username, err)
		return nil, err
	}
	requests, err := s.repo.GetSent(userID)
	if err != nil {
		log.Printf("GetSentRequests repository error for user %d: %v", userID, err)
		return nil, err
	}

	for i := range requests {
		sender, err := repositories.GetUsernameByID(s.db, requests[i].SenderID)
		if err != nil {
			log.Printf("GetSentRequests sender lookup error for user %d: %v", requests[i].SenderID, err)
			return nil, err
		}
		receiver, err := repositories.GetUsernameByID(s.db, requests[i].ReceiverID)
		if err != nil {
			log.Printf("GetSentRequests receiver lookup error for user %d: %v", requests[i].ReceiverID, err)
			return nil, err
		}
		requests[i].SenderUsername = sender
		requests[i].ReceiverUsername = receiver
	}

	return requests, nil
}

// CheckRequestStatus checks if a friend request exists between sender and receiver.
func (s *FriendRequestService) CheckRequestStatus(username string, receiverID int) (bool, error) {
	senderID, err := repositories.GetUserIDByUsername(s.db, username)
	if err != nil {
		log.Printf("CheckRequestStatus user lookup error for %s: %v", username, err)
		return false, err
	}
	count, err := s.repo.Count(senderID, receiverID)
	if err != nil {
		log.Printf("CheckRequestStatus count error for sender %d and receiver %d: %v", senderID, receiverID, err)
		return false, err
	}
	return count > 0, nil
}
