package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
)

func GetQuestionnaire(ctx *gin.Context) {
	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	questions, err := questionnaireService.GetQuestionnaire()
	if err != nil {
		log.Printf("GetQuestionnaire service error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"questions": questions})
}

func SubmitQuestionnaire(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var answer models.Answer
	if err := ctx.ShouldBindJSON(&answer); err != nil {
		log.Printf("SubmitQuestionnaire bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	if err := questionnaireService.SubmitAnswer(username.(string), answer); err != nil {
		log.Printf("SubmitQuestionnaire service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Submit Answer"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Answers submitted successfully"})
}

func GetUserAnswers(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	questionnaireService := ctx.MustGet("questionnaireService").(*services.QuestionnaireService)

	answers, err := questionnaireService.GetUserAnswers(username.(string))
	if err != nil {
		log.Printf("GetUserAnswers service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answers"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"answers": answers})
}
