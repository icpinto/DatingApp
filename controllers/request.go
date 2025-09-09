package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func SendFriendRequest(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "SendFriendRequest unauthorized", "Unauthorized")
		return
	}

	var request models.FriendRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "SendFriendRequest bind error", "Invalid request data")
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.SendFriendRequest(username.(string), request); err != nil {
		if errors.Is(err, services.ErrFriendRequestExists) {
			logMsg := fmt.Sprintf("SendFriendRequest duplicate between %s and %d", username.(string), request.ReceiverID)
			utils.RespondError(ctx, http.StatusConflict, err, logMsg, err.Error())
			return
		}
		logMsg := fmt.Sprintf("SendFriendRequest service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusBadRequest, err, logMsg, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

func AcceptFriendRequest(ctx *gin.Context) {
	var request models.AcceptRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "AcceptFriendRequest bind error", err.Error())
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.AcceptFriendRequest(request.RequestID); err != nil {
		logMsg := fmt.Sprintf("AcceptFriendRequest service error for request %d", request.RequestID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to accept request")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Friend request accepted and conversation created"})
}

func RejectFriendRequest(ctx *gin.Context) {
	var request models.RejectRequest

	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "RejectFriendRequest bind error", "Invalid request data", err.Error())
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	if err := frService.RejectFriendRequest(request.RequestID); err != nil {
		logMsg := fmt.Sprintf("RejectFriendRequest service error for request %d", request.RequestID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to reject request", err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Friend request rejected successfully"})
}

func GetPendingRequests(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetPendingRequests unauthorized", "Unauthorized")
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	requests, err := frService.GetPendingRequests(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetPendingRequests service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"requests": requests})
}

func CheckReqStatus(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "CheckReqStatus unauthorized", "Unauthorized")
		return
	}

	receiverIDParam := ctx.Param("reciver_id")
	receiverID, err := strconv.Atoi(receiverIDParam)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "CheckReqStatus invalid receiver id", "Invalid receiver id")
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	requestSent, err := frService.CheckRequestStatus(username.(string), receiverID)
	if err != nil {
		logMsg := fmt.Sprintf("CheckReqStatus service error for %s and %d", username.(string), receiverID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve request status")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"requestStatus": requestSent})
}
