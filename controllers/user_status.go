package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

// GetUserStatus godoc
// @Summary      Retrieve the authenticated user's activation status
// @Description  Returns whether the caller's account is currently activated or deactivated.
// @Tags         User
// @Produce      json
// @Success      200  {object}  utils.UserStatusResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Router       /user/status [get]
// @Security     BearerAuth
func GetUserStatus(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)

	userIDVal, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetUserStatus missing user id", "Unauthorized")
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetUserStatus invalid user id", "Unauthorized")
		return
	}

	isActive, err := userService.GetUserStatus(userID)
	if err != nil {
		if err == repositories.ErrUserNotFound {
			utils.RespondError(ctx, http.StatusNotFound, err, "GetUserStatus user not found", "user not found")
			return
		}
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetUserStatus service error", "Could not retrieve user status")
		return
	}

	status := "deactivated"
	if isActive {
		status = "activated"
	}

	utils.RespondSuccess(ctx, http.StatusOK, utils.UserStatusResponse{Status: status})
}
