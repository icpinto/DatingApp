package repositories

import (
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/icpinto/dating-app/models"
)

// FriendRequestRepository manages CRUD operations for friend requests.
type FriendRequestRepository struct {
	db *sql.DB
}

// NewFriendRequestRepository creates a new FriendRequestRepository.
func NewFriendRequestRepository(db *sql.DB) *FriendRequestRepository {
	return &FriendRequestRepository{db: db}
}

// CheckExisting returns the status of an existing friend request.
func (r *FriendRequestRepository) CheckExisting(senderID, receiverID int) (string, error) {
	var status string
	err := r.db.QueryRow(`
        SELECT status FROM friend_requests WHERE sender_id = $1 AND receiver_id = $2`,
		senderID, receiverID).Scan(&status)
	if err != nil {
		log.Printf("FriendRequestRepository.CheckExisting query error for sender %d and receiver %d: %v", senderID, receiverID, err)
	}
	return status, err
}

// Create inserts a new friend request.
func (r *FriendRequestRepository) Create(request models.FriendRequest) error {
	_, err := r.db.Exec(`
            INSERT INTO friend_requests (sender_id, sender_username, receiver_id, receiver_username, status, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		request.SenderID, request.SenderUsername, request.ReceiverID, request.ReceiverUsername, request.Status, request.CreatedAt, request.UpdatedAt)
	if err != nil {
		log.Printf("FriendRequestRepository.Create exec error for sender %d and receiver %d: %v", request.SenderID, request.ReceiverID, err)
	}
	return err
}

// UpdateStatus updates the status of a friend request.
func (r *FriendRequestRepository) UpdateStatus(requestID int, status string, updatedAt time.Time) error {
	_, err := r.db.Exec(`
        UPDATE friend_requests
        SET status = $1, updated_at = $2
        WHERE id = $3`,
		status, updatedAt, requestID)
	if err != nil {
		log.Printf("FriendRequestRepository.UpdateStatus exec error for request %d: %v", requestID, err)
	}
	return err
}

// GetUsers returns the sender and receiver IDs for a request.
func (r *FriendRequestRepository) GetUsers(requestID int) (int, int, error) {
	var user1ID, user2ID int
	err := r.db.QueryRow(`
        SELECT sender_id, receiver_id
        FROM friend_requests
        WHERE id = $1`, requestID).Scan(&user1ID, &user2ID)
	if err != nil {
		log.Printf("FriendRequestRepository.GetUsers query error for request %d: %v", requestID, err)
	}
	return user1ID, user2ID, err
}

// GetPending retrieves all pending friend requests for a user.
func (r *FriendRequestRepository) GetPending(userID int) ([]models.FriendRequest, error) {
	rows, err := r.db.Query(`
       SELECT id, sender_id, sender_username, receiver_id, receiver_username, status, created_at
       FROM friend_requests
       WHERE receiver_id = $1 AND status = 'pending'`, userID)
	if err != nil {
		log.Printf("FriendRequestRepository.GetPending query error for user %d: %v", userID, err)
		return nil, err
	}
	defer rows.Close()

	var requests []models.FriendRequest
	for rows.Next() {
		var request models.FriendRequest
		if err := rows.Scan(&request.RequestId, &request.SenderID, &request.SenderUsername, &request.ReceiverID, &request.ReceiverUsername, &request.Status, &request.CreatedAt); err != nil {
			log.Printf("FriendRequestRepository.GetPending scan error: %v", err)
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		log.Printf("FriendRequestRepository.GetPending rows error: %v", err)
		return nil, err
	}
	return requests, nil
}

// Count returns the number of requests from sender to receiver.
func (r *FriendRequestRepository) Count(senderID, receiverID int) (int, error) {
	var count int
	err := r.db.QueryRow(`
                SELECT COUNT(*)
                FROM friend_requests
                WHERE sender_id = $1 AND receiver_id = $2`,
		senderID, receiverID).Scan(&count)
	if err != nil {
		log.Printf("FriendRequestRepository.Count query error for sender %d and receiver %d: %v", senderID, receiverID, err)
	}
	return count, err
}

// Delete removes a friend request.
func (r *FriendRequestRepository) Delete(requestID int) error {
	if _, err := r.db.Exec(`DELETE FROM friend_requests WHERE id = $1`, requestID); err != nil {
		log.Printf("FriendRequestRepository.Delete exec error for request %d: %v", requestID, err)
		return err
	}
	return nil
}

// LinkConversation stores the conversation ID for an accepted match.
func (r *FriendRequestRepository) LinkConversation(senderID, receiverID int, conversationID uuid.UUID) error {
	_, err := r.db.Exec(`
        UPDATE friend_requests
        SET conversation_id = $1
        WHERE sender_id = $2 AND receiver_id = $3`, conversationID, senderID, receiverID)
	if err != nil {
		log.Printf("FriendRequestRepository.LinkConversation exec error for users %d and %d: %v", senderID, receiverID, err)
	}
	return err
}
