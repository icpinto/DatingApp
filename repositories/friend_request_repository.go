package repositories

import (
	"database/sql"
	"time"

	"github.com/icpinto/dating-app/models"
)

func CheckExistingRequest(db *sql.DB, senderID, receiverID int) (string, error) {
	var status string
	err := db.QueryRow(`
        SELECT status FROM friend_requests WHERE sender_id = $1 AND receiver_id = $2`,
		senderID, receiverID).Scan(&status)
	return status, err
}

func InsertFriendRequest(db *sql.DB, request models.FriendRequest) error {
	_, err := db.Exec(`
            INSERT INTO friend_requests (sender_id, receiver_id, status, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5)`,
		request.SenderID, request.ReceiverID, request.Status, request.CreatedAt, request.UpdatedAt)
	return err
}

func UpdateFriendRequestStatus(db *sql.DB, requestID int, status string, updatedAt time.Time) error {
	_, err := db.Exec(`
        UPDATE friend_requests
        SET status = $1, updated_at = $2
        WHERE id = $3`,
		status, updatedAt, requestID)
	return err
}

func GetFriendRequestUsers(db *sql.DB, requestID int) (int, int, error) {
	var user1ID, user2ID int
	err := db.QueryRow(`
        SELECT sender_id, receiver_id
        FROM friend_requests
        WHERE id = $1`, requestID).Scan(&user1ID, &user2ID)
	return user1ID, user2ID, err
}

func GetPendingRequests(db *sql.DB, userID int) ([]models.FriendRequest, error) {
	rows, err := db.Query(`
        SELECT id, sender_id, receiver_id, status, created_at
        FROM friend_requests
        WHERE receiver_id = $1 AND status = 'pending'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.FriendRequest
	for rows.Next() {
		var request models.FriendRequest
		if err := rows.Scan(&request.RequestId, &request.SenderID, &request.ReceiverID, &request.Status, &request.CreatedAt); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, rows.Err()
}

func CountFriendRequests(db *sql.DB, senderID, receiverID int) (int, error) {
	var count int
	err := db.QueryRow(`
                SELECT COUNT(*)
                FROM friend_requests
                WHERE sender_id = $1 AND receiver_id = $2`,
		senderID, receiverID).Scan(&count)
	return count, err
}
