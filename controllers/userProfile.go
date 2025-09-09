package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
)

func CreateProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	var profile models.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		log.Printf("CreateProfile bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := profileService.CreateOrUpdateProfile(username.(string), profile); err != nil {
		log.Printf("CreateProfile service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func GetProfile(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profile, err := profileService.GetProfile(username.(string))
	if err != nil {
		log.Printf("GetProfile service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func GetProfiles(ctx *gin.Context) {
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	profiles, err := profileService.GetProfiles()
	if err != nil {
		log.Printf("GetProfiles service error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
		return
	}
	ctx.JSON(http.StatusOK, profiles)
}

func GetUserProfile(ctx *gin.Context) {
	userIDParam := ctx.Param("user_id")

	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		log.Printf("GetUserProfile invalid user id: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user id"})
		return
	}

	profile, err := profileService.GetProfileByUserID(userID)
	if err != nil {
		log.Printf("GetUserProfile service error for user %d: %v", userID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
