package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

// GetQuestionnaire godoc
// @Summary      Retrieve matchmaking questionnaire
// @Tags         Questionnaire
// @Produce      json
// @Success      200  {object}  utils.QuestionnaireResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/questionnaire [get]
func GetQuestionnaire(ctx *gin.Context) {
	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	questions, err := questionnaireService.GetQuestionnaire()
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "GetQuestionnaire service error", "Failed to fetch questions")
		return
	}
	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"questions": questions})
}

// SubmitQuestionnaire godoc
// @Summary      Submit questionnaire answers
// @Tags         Questionnaire
// @Accept       json
// @Produce      json
// @Param        answer  body      models.Answer  true  "Questionnaire answer"
// @Success      200     {object}  utils.MessageResponse
// @Failure      400     {object}  utils.ErrorResponse
// @Failure      401     {object}  utils.ErrorResponse
// @Failure      500     {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/submitQuestionnaire [post]
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

// GetUserAnswers godoc
// @Summary      Retrieve questionnaire answers for the authenticated user
// @Tags         Questionnaire
// @Produce      json
// @Success      200  {object}  utils.AnswersResponse
// @Failure      401  {object}  utils.ErrorResponse
// @Failure      500  {object}  utils.ErrorResponse
// @Security     BearerAuth
// @Router       /user/questionnaireAnswers [get]
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
