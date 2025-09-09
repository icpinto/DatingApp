package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func CreateProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "CreateProfile unauthorized", "Unauthorized")
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	var profile models.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "CreateProfile bind error", "Invalid input")
		return
	}

	if err := profileService.CreateOrUpdateProfile(username.(string), profile); err != nil {
		logMsg := fmt.Sprintf("CreateProfile service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to update profile")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetProfile unauthorized", "Unauthorized")
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profile, err := profileService.GetProfile(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetProfile service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve profile")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, profile)
}

func GetProfiles(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profiles, err := profileService.GetProfiles()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetProfiles service error", "Failed to retrieve profiles")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, profiles)
}

func GetUserProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "GetUserProfile invalid user id", "Invalid user id")
		return
	}

	profile, err := profileService.GetProfileByUserID(userID)
	if err != nil {
		logMsg := fmt.Sprintf("GetUserProfile service error for user %d", userID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve profile")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, profile)
}
