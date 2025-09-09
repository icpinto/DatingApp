package services

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

var ErrFriendRequestExists = errors.New("friend request already exists")

func SendFriendRequest(db *sql.DB, username string, request models.FriendRequest) error {
	senderID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("SendFriendRequest user lookup error for %s: %v", username, err)
		return err
	}
	request.SenderID = senderID
	request.Status = "pending"
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err = repositories.CheckExistingRequest(db, request.SenderID, request.ReceiverID)
	if err == nil {
		log.Printf("SendFriendRequest duplicate for sender %d and receiver %d", request.SenderID, request.ReceiverID)
		return ErrFriendRequestExists
	}
	if err != sql.ErrNoRows {
		log.Printf("SendFriendRequest check existing error: %v", err)
		return err
	}
	if err := repositories.InsertFriendRequest(db, request); err != nil {
		log.Printf("SendFriendRequest insert error for sender %d and receiver %d: %v", request.SenderID, request.ReceiverID, err)
		return err
	}
	return nil
}

func AcceptFriendRequest(db *sql.DB, requestID int) error {
	if err := repositories.UpdateFriendRequestStatus(db, requestID, "accepted", time.Now()); err != nil {
		log.Printf("AcceptFriendRequest update status error for request %d: %v", requestID, err)
		return err
	}
	user1ID, user2ID, err := repositories.GetFriendRequestUsers(db, requestID)
	if err != nil {
		log.Printf("AcceptFriendRequest get users error for request %d: %v", requestID, err)
		return err
	}
	if err := repositories.CreateConversation(db, user1ID, user2ID); err != nil {
		log.Printf("AcceptFriendRequest create conversation error for request %d: %v", requestID, err)
		return err
	}
	return nil
}

func RejectFriendRequest(db *sql.DB, requestID int) error {
	if err := repositories.UpdateFriendRequestStatus(db, requestID, "rejected", time.Now()); err != nil {
		log.Printf("RejectFriendRequest update status error for request %d: %v", requestID, err)
		return err
	}
	return nil
}

func GetPendingRequests(db *sql.DB, username string) ([]models.FriendRequest, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("GetPendingRequests user lookup error for %s: %v", username, err)
		return nil, err
	}
	requests, err := repositories.GetPendingRequests(db, userID)
	if err != nil {
		log.Printf("GetPendingRequests repository error for user %d: %v", userID, err)
		return nil, err
	}
	return requests, nil
}

func CheckRequestStatus(db *sql.DB, username string, receiverID int) (bool, error) {
	senderID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		log.Printf("CheckRequestStatus user lookup error for %s: %v", username, err)
		return false, err
	}
	count, err := repositories.CountFriendRequests(db, senderID, receiverID)
	if err != nil {
		log.Printf("CheckRequestStatus count error for sender %d and receiver %d: %v", senderID, receiverID, err)
		return false, err
	}
	return count > 0, nil
}
