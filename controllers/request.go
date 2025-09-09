package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
)

func SendFriendRequest(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var request models.FriendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("SendFriendRequest bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.SendFriendRequest(username.(string), request); err != nil {
		if errors.Is(err, services.ErrFriendRequestExists) {
			log.Printf("SendFriendRequest duplicate between %s and %d: %v", username.(string), request.ReceiverID, err)
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		log.Printf("SendFriendRequest service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

func AcceptFriendRequest(ctx *gin.Context) {
	var request models.AcceptRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("AcceptFriendRequest bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.AcceptFriendRequest(request.RequestID); err != nil {
		log.Printf("AcceptFriendRequest service error for request %d: %v", request.RequestID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept request"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request accepted and conversation created"})
}

func RejectFriendRequest(ctx *gin.Context) {
	var request models.RejectRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		log.Printf("RejectFriendRequest bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.RejectFriendRequest(request.RequestID); err != nil {
		log.Printf("RejectFriendRequest service error for request %d: %v", request.RequestID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject request", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

func GetPendingRequests(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	requests, err := frService.GetPendingRequests(username.(string))
	if err != nil {
		log.Printf("GetPendingRequests service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"requests": requests})
}

func CheckReqStatus(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	receiverIDParam := ctx.Param("reciver_id")
	receiverID, err := strconv.Atoi(receiverIDParam)
	if err != nil {
		log.Printf("CheckReqStatus invalid receiver id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid receiver id"})
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	requestSent, err := frService.CheckRequestStatus(username.(string), receiverID)
	if err != nil {
		log.Printf("CheckReqStatus service error for %s and %d: %v", username.(string), receiverID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve request status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"requestStatus": requestSent})
}
