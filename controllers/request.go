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

// SendFriendRequest godoc
// @Summary      Send a friend request
// @Description  Sends a friend request from the authenticated user to another user.
// @Tags         Friend Requests
// @Accept       json
// @Produce      json
// @Param        request  body      models.FriendRequest  true  "Friend request payload"
// @Success      200      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      401      {object}  utils.ErrorResponse
// @Failure      409      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/sendRequest [post]
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

// AcceptFriendRequest godoc
// @Summary      Accept a friend request
// @Tags         Friend Requests
// @Accept       json
// @Produce      json
// @Param        request  body      models.AcceptRequest  true  "Friend request to accept"
// @Success      200      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/acceptRequest [post]
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

// RejectFriendRequest godoc
// @Summary      Reject a friend request
// @Tags         Friend Requests
// @Accept       json
// @Produce      json
// @Param        request  body      models.RejectRequest  true  "Friend request to reject"
// @Success      200      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/rejectRequest [post]
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

// GetPendingRequests godoc
// @Summary      List pending friend requests
// @Tags         Friend Requests
// @Produce      json
// @Success      200  {object}  utils.FriendRequestsResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/requests [get]
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

// GetSentRequests godoc
// @Summary      List sent friend requests
// @Tags         Friend Requests
// @Produce      json
// @Success      200  {object}  utils.FriendRequestsResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/sentRequests [get]
func GetSentRequests(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetSentRequests unauthorized", "Unauthorized")
		return
	}

	frService := ctx.MustGet("friendRequestService").(*services.FriendRequestService)

	requests, err := frService.GetSentRequests(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetSentRequests service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, err.Error())
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"requests": requests})
}

// CheckReqStatus godoc
// @Summary      Check the friend request status with another user
// @Tags         Friend Requests
// @Produce      json
// @Param        reciver_id  path      int  true  "Receiver user ID"
// @Success      200         {object}  utils.FriendRequestStatusResponse
// @Failure      400         {object}  utils.ErrorResponse
// @Failure      401         {object}  utils.ErrorResponse
// @Failure      500         {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/checkReqStatus/{reciver_id} [get]
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
