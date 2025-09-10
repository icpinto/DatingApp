package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

// ErrFriendRequestExists indicates a duplicate friend request.
var ErrFriendRequestExists = errors.New("friend request already exists")

// FriendRequestService provides operations related to friend requests.
type FriendRequestService struct {
	db        *sql.DB
	repo      *repositories.FriendRequestRepository
	convoRepo *repositories.ConversationRepository
}

// NewFriendRequestService creates a new FriendRequestService.
func NewFriendRequestService(db *sql.DB) *FriendRequestService {
	return &FriendRequestService{db: db, repo: repositories.NewFriendRequestRepository(db), convoRepo: repositories.NewConversationRepository(db)}
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
	if err := s.repo.UpdateStatus(requestID, "accepted", time.Now()); err != nil {
		log.Printf("AcceptFriendRequest update status error for request %d: %v", requestID, err)
		return err
	}
	user1ID, user2ID, err := s.repo.GetUsers(requestID)
	if err != nil {
		log.Printf("AcceptFriendRequest get users error for request %d: %v", requestID, err)
		return err
	}
	user1Username, err := repositories.GetUsernameByID(s.db, user1ID)
	if err != nil {
		log.Printf("AcceptFriendRequest user1 lookup error for %d: %v", user1ID, err)
		return err
	}
	user2Username, err := repositories.GetUsernameByID(s.db, user2ID)
	if err != nil {
		log.Printf("AcceptFriendRequest user2 lookup error for %d: %v", user2ID, err)
		return err
	}
	if err := s.convoRepo.Create(user1ID, user1Username, user2ID, user2Username); err != nil {
		log.Printf("AcceptFriendRequest create conversation error for request %d: %v", requestID, err)
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
