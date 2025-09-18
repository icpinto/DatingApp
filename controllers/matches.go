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

// GetUserMatches returns best matches for the specified user by combining microservice scores with profile data.
func GetUserMatches(ctx *gin.Context) {
	matchService := ctx.MustGet("matchService").(*services.MatchService)
	profileService := ctx.MustGet("profileService").(*services.ProfileService)

	userIDParam := ctx.Param("user_id")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "GetUserMatches invalid user id", "Invalid user id")
		return
	}

	matches, err := matchService.GetMatches(ctx.Request.Context(), userID, ctx.Request.URL.RawQuery)
	if err != nil {
		logMsg := fmt.Sprintf("GetUserMatches match service error for user %d", userID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve matches")
		return
	}

	ids := make([]int, 0, len(matches))
	for _, m := range matches {
		ids = append(ids, m.UserID)
	}

	profilesByID, err := profileService.GetProfilesByUserIDs(ids)
	if err != nil {
		logMsg := fmt.Sprintf("GetUserMatches profile lookup error for user %d", userID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve matches")
		return
	}

	matchedProfiles := make([]models.MatchedProfile, 0, len(matches))
	for _, m := range matches {
		profile, ok := profilesByID[m.UserID]
		if !ok {
			continue
		}
		matchedProfiles = append(matchedProfiles, models.MatchedProfile{
			UserProfile: profile,
			Score:       m.Score,
			Reasons:     m.Reasons,
		})
	}

	utils.RespondSuccess(ctx, http.StatusOK, matchedProfiles)
}
