package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func GetQuestionnaire(ctx *gin.Context) {
	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	questions, err := questionnaireService.GetQuestionnaire()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetQuestionnaire service error", "Failed to fetch questions")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"questions": questions})
}

func SubmitQuestionnaire(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "SubmitQuestionnaire unauthorized", "Unauthorized")
		return
	}

	var answer models.Answer
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "SubmitQuestionnaire bind error", "Invalid request data")
		return
	}

	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	if err := questionnaireService.SubmitAnswer(username.(string), answer); err != nil {
		logMsg := fmt.Sprintf("SubmitQuestionnaire service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to Submit Answer")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Answers submitted successfully"})
}

func GetUserAnswers(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetUserAnswers unauthorized", "Unauthorized")
		return
	}

	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	answers, err := questionnaireService.GetUserAnswers(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetUserAnswers service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve answers")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"answers": answers})
}
