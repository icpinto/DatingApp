package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

// SaveCorePreferences godoc
// @Summary      Save user core preferences
// @Tags         Core Preferences
// @Accept       json
// @Produce      json
// @Param        preferences  body      models.CorePreferences  true  "Core preferences payload"
// @Success      200          {object}  models.CorePreferences
// @Failure      400          {object}  utils.ErrorResponse
// @Failure      401          {object}  utils.ErrorResponse
// @Failure      500          {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/core-preferences [post]
func SaveCorePreferences(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "SaveCorePreferences unauthorized", "Unauthorized")
		return
	}

	var prefs models.CorePreferences
	if err := ctx.ShouldBindJSON(&prefs); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "SaveCorePreferences bind error", "Invalid request data")
		return
	}

	prefs.UserID = userID.(int)

	matchService := ctx.MustGet("matchService").(*services.MatchService)
	saved, err := matchService.SaveCorePreferences(ctx.Request.Context(), prefs)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "SaveCorePreferences service error", "Failed to save core preferences")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, saved)
}

// UpdateCorePreferences godoc
// @Summary      Update user core preferences
// @Tags         Core Preferences
// @Accept       json
// @Produce      json
// @Param        preferences  body      models.CorePreferences  true  "Core preferences payload"
// @Success      200          {object}  models.CorePreferences
// @Failure      400          {object}  utils.ErrorResponse
// @Failure      401          {object}  utils.ErrorResponse
// @Failure      500          {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/core-preferences [put]
func UpdateCorePreferences(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "UpdateCorePreferences unauthorized", "Unauthorized")
		return
	}

	var prefs models.CorePreferences
	if err := ctx.ShouldBindJSON(&prefs); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "UpdateCorePreferences bind error", "Invalid request data")
		return
	}

	prefs.UserID = userID.(int)

	matchService := ctx.MustGet("matchService").(*services.MatchService)
	updated, err := matchService.UpdateCorePreferences(ctx.Request.Context(), prefs)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "UpdateCorePreferences service error", "Failed to update core preferences")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, updated)
}

// GetCorePreferences godoc
// @Summary      Get user core preferences
// @Tags         Core Preferences
// @Produce      json
// @Success      200  {object}  models.CorePreferences
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      404  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/core-preferences [get]
func GetCorePreferences(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetCorePreferences unauthorized", "Unauthorized")
		return
	}

	matchService := ctx.MustGet("matchService").(*services.MatchService)
	prefs, err := matchService.GetCorePreferences(ctx.Request.Context(), userID.(int))
	if err != nil {
		if err.Error() == "core preferences not found" {
			utils.RespondError(ctx, http.StatusNotFound, err, "GetCorePreferences not found", "Core preferences not found")
			return
		}

		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetCorePreferences service error", "Failed to fetch core preferences")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, prefs)
}
