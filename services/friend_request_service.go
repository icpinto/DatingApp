package services

import (
	"database/sql"
	"errors"
	"time"

	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
)

func SendFriendRequest(db *sql.DB, username string, request models.FriendRequest) error {
	senderID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return err
	}
	request.SenderID = senderID
	request.Status = "pending"
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	_, err = repositories.CheckExistingRequest(db, request.SenderID, request.ReceiverID)
	if err == nil {
		return errors.New("friend request already exists")
	}
	if err != sql.ErrNoRows {
		return err
	}
	return repositories.InsertFriendRequest(db, request)
}

func AcceptFriendRequest(db *sql.DB, requestID int) error {
	if err := repositories.UpdateFriendRequestStatus(db, requestID, "accepted", time.Now()); err != nil {
		return err
	}
	user1ID, user2ID, err := repositories.GetFriendRequestUsers(db, requestID)
	if err != nil {
		return err
	}
	return repositories.CreateConversation(db, user1ID, user2ID)
}

func RejectFriendRequest(db *sql.DB, requestID int) error {
	return repositories.UpdateFriendRequestStatus(db, requestID, "rejected", time.Now())
}

func GetPendingRequests(db *sql.DB, username string) ([]models.FriendRequest, error) {
	userID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return nil, err
	}
	return repositories.GetPendingRequests(db, userID)
}

func CheckRequestStatus(db *sql.DB, username string, receiverID int) (bool, error) {
	senderID, err := repositories.GetUserIDByUsername(db, username)
	if err != nil {
		return false, err
	}
	count, err := repositories.CountFriendRequests(db, senderID, receiverID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
