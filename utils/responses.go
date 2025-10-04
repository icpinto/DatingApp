package utils

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
)

// RespondSuccess sends a JSON response with the provided status and payload.
func RespondSuccess(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}

// RespondError logs the error with the provided log message and sends a JSON response with a client-facing message.
// An optional details parameter can be included for additional client context.
func RespondError(ctx *gin.Context, status int, err error, logMsg, clientMsg string, details ...string) {
	if err != nil {
		log.Printf("%s: %v", logMsg, err)
	} else {
		log.Println(logMsg)
	}

	response := gin.H{"error": clientMsg}
	if len(details) > 0 {
		response["details"] = details[0]
	}

	ctx.JSON(status, response)
}

// MessageResponse represents a generic message payload.
type MessageResponse struct {
	Message string `json:"message"`
}

// TokenResponse represents a JWT token response.
type TokenResponse struct {
	Token  string `json:"token"`
	UserID int    `json:"user_id"`
}

// ErrorResponse represents an error message returned to the client.
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

// FriendRequestsResponse wraps a list of friend requests.
type FriendRequestsResponse struct {
	Requests []models.FriendRequest `json:"requests"`
}

// FriendRequestStatusResponse represents the friend request status between two users.
type FriendRequestStatusResponse struct {
	RequestStatus bool `json:"requestStatus"`
}

// UserStatusResponse represents a user's activation status.
type UserStatusResponse struct {
	Status string `json:"status"`
}
