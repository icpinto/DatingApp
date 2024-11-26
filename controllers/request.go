package controllers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/models"
)

// SendFriendRequest sends a friend request
func SendFriendRequest(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request models.FriendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	var senderID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}
	request.SenderID = senderID
	request.Status = "pending"
	request.CreatedAt = time.Now()
	request.UpdatedAt = time.Now()

	// Check if the request already exists
	var existingStatus string
	err = db.DB.QueryRow(`
        SELECT status FROM friend_requests WHERE sender_id = $1 AND receiver_id = $2`,
		request.SenderID, request.ReceiverID).Scan(&existingStatus)

	if err == sql.ErrNoRows {
		// No existing request, proceed to insert
		_, err = db.DB.Exec(`
            INSERT INTO friend_requests (sender_id, receiver_id, status, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5)`,
			request.SenderID, request.ReceiverID, request.Status, request.CreatedAt, request.UpdatedAt)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
	} else if err == nil {
		// If there’s already a request between the users
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Friend request already exists"})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing request"})
	}
}

// AcceptFriendRequest accepts a friend request
func AcceptFriendRequest(ctx *gin.Context) {

	var request models.AcceptRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		//ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	// Update the friend request status to "accepted"
	_, err := db.DB.Exec(`
        UPDATE friend_requests
        SET status = $1, updated_at = $2
        WHERE id = $3 `,
		"accepted", time.Now(), request.RequestID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
}

// RejectFriendRequest rejects a friend request
func RejectFriendRequest(ctx *gin.Context) {
	var request models.RejectRequest

	// Parse and validate incoming JSON
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	// Update the friend request status to "rejected"
	_, err := db.DB.Exec(`
        UPDATE friend_requests
        SET status = $1, updated_at = $2
        WHERE id = $3`,
		"rejected", time.Now(), request.RequestID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject request", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

// GetPendingRequests retrieves all pending friend requests for a user
func GetPendingRequests(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var userID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	rows, err := db.DB.Query(`
        SELECT id,sender_id, receiver_id, status, created_at
        FROM friend_requests
        WHERE receiver_id = $1 AND status = 'pending'`,
		userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		//ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve requests"})
		return
	}
	defer rows.Close()

	var requests []models.FriendRequest
	for rows.Next() {
		var request models.FriendRequest
		if err := rows.Scan(&request.RequestId, &request.SenderID, &request.ReceiverID, &request.Status, &request.CreatedAt); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan request data"})
			return
		}
		requests = append(requests, request)
	}

	ctx.JSON(http.StatusOK, gin.H{"requests": requests})
}

// check request status between two users
func CheckReqStatus(ctx *gin.Context) {
	// Extract username from context (authenticated user)
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get the receiver ID from the route parameter
	receiverID := ctx.Param("reciver_id")

	// Get the sender's user ID based on the authenticated username
	var senderID int
	err := db.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&senderID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve sender ID"})
		return
	}

	// Query to check if a friend request exists between the sender and receiver
	var requestCount int
	err = db.DB.QueryRow(`
		SELECT COUNT(*) 
		FROM friend_requests 
		WHERE sender_id = $1 AND receiver_id = $2`,
		senderID, receiverID).Scan(&requestCount)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve request status"})
		return
	}

	// If requestCount > 0, a request has already been sent
	requestSent := requestCount > 0

	// Respond with the request status
	ctx.JSON(http.StatusOK, gin.H{"requestStatus": requestSent})
}
