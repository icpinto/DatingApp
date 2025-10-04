package controllers

import (
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

type lifecycleRequest struct {
	Reason string `json:"reason"`
}

// DeactivateCurrentUser godoc
// @Summary      Deactivate the authenticated user
// @Description  Marks the caller's account as inactive and schedules downstream cleanup.
// @Tags         User
// @Produce      json
// @Param        payload  body      lifecycleRequest  false  "Optional deactivation reason"
// @Success      202      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      401      {object}  utils.ErrorResponse
// @Failure      404      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Router       /user/deactivate [post]
// @Security     BearerAuth
func DeactivateCurrentUser(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Deactivate missing user id", "Unauthorized")
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Deactivate invalid user id", "Unauthorized")
		return
	}

	var req lifecycleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			req = lifecycleRequest{}
		} else {
			utils.RespondError(ctx, http.StatusBadRequest, err, "Deactivate bind error", "Invalid input")
			return
		}
	}

	if err := userService.DeactivateUser(ctx.Request.Context(), userID, req.Reason); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.RespondError(ctx, http.StatusNotFound, err, "Deactivate user not found", "user not found")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err, "Deactivate service error", "Could not deactivate user")
		return
	}

	utils.RespondSuccess(ctx, http.StatusAccepted, utils.MessageResponse{Message: "Account deactivation scheduled"})
}

// ReactivateCurrentUser godoc
// @Summary      Reactivate the authenticated user
// @Description  Marks the caller's account as active and schedules downstream restoration.
// @Tags         User
// @Produce      json
// @Param        payload  body      lifecycleRequest  false  "Optional reactivation reason"
// @Success      202      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      401      {object}  utils.ErrorResponse
// @Failure      404      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Router       /user/reactivate [post]
// @Security     BearerAuth
func ReactivateCurrentUser(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Reactivate missing user id", "Unauthorized")
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Reactivate invalid user id", "Unauthorized")
		return
	}

	var req lifecycleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			req = lifecycleRequest{}
		} else {
			utils.RespondError(ctx, http.StatusBadRequest, err, "Reactivate bind error", "Invalid input")
			return
		}
	}

	if err := userService.ReactivateUser(ctx.Request.Context(), userID, req.Reason); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.RespondError(ctx, http.StatusNotFound, err, "Reactivate user not found", "user not found")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err, "Reactivate service error", "Could not reactivate user")
		return
	}

	utils.RespondSuccess(ctx, http.StatusAccepted, utils.MessageResponse{Message: "Account reactivation scheduled"})
}

// DeleteCurrentUser godoc
// @Summary      Permanently delete the authenticated user
// @Description  Removes the caller's account and enqueues delete events for downstream services.
// @Tags         User
// @Produce      json
// @Param        payload  body      lifecycleRequest  false  "Optional deletion reason"
// @Success      202      {object}  utils.MessageResponse
// @Failure      400      {object}  utils.ErrorResponse
// @Failure      401      {object}  utils.ErrorResponse
// @Failure      404      {object}  utils.ErrorResponse
// @Failure      500      {object}  utils.ErrorResponse
// @Router       /user [delete]
// @Security     BearerAuth
func DeleteCurrentUser(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Delete missing user id", "Unauthorized")
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "Delete invalid user id", "Unauthorized")
		return
	}

	var req lifecycleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		if errors.Is(err, io.EOF) {
			req = lifecycleRequest{}
		} else {
			utils.RespondError(ctx, http.StatusBadRequest, err, "Delete bind error", "Invalid input")
			return
		}
	}

	if err := userService.DeleteUser(ctx.Request.Context(), userID, req.Reason); err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.RespondError(ctx, http.StatusNotFound, err, "Delete user not found", "user not found")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err, "Delete service error", "Could not delete user")
		return
	}

	utils.RespondSuccess(ctx, http.StatusAccepted, utils.MessageResponse{Message: "Account deletion scheduled"})
}
