package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

// GetQuestionnaire returns the list of questions for the questionnaire
func GetQuestionnaire(ctx *gin.Context) {
	rows, err := db.DB.Query(`SELECT id, question_text, question_type, options FROM questions`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch questions"})
		return
	}
	defer rows.Close()

	var questions []models.Question
	for rows.Next() {
		var question models.Question
		if err := rows.Scan(&question.ID, &question.QuestionText, &question.QuestionType, pq.Array(&question.Options)); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan question data"})
			return
		}
		questions = append(questions, question)
	}

	ctx.JSON(http.StatusOK, gin.H{"questions": questions})
}

// SubmitQuestionnaire stores user answers in the database
func SubmitQuestionnaire(ctx *gin.Context) {
	userID, exists := ctx.Get("userID") // Extract user ID from JWT
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req models.SubmitAnswersRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
		return
	}

	tx, err := db.DB.Begin() // Start a transaction
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	for _, answer := range req.Answers {
		_, err := tx.Exec(`
            INSERT INTO user_answers (user_id, question_id, answer_text, answer_value, created_at)
            VALUES ($1, $2, $3, $4, $5)`,
			userID, answer.QuestionID, answer.AnswerText, answer.AnswerValue, time.Now())

		if err != nil {
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store answers"})
			return
		}
	}

	if err := tx.Commit(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Answers submitted successfully"})
}

// GetUserAnswers retrieves the answers for a given user
func GetUserAnswers(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	rows, err := db.DB.Query(`
        SELECT question_id, answer_text, answer_value
        FROM user_answers WHERE user_id = $1`, userID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answers"})
		return
	}
	defer rows.Close()

	var answers []models.Answer
	for rows.Next() {
		var answer models.Answer
		if err := rows.Scan(&answer.QuestionID, &answer.AnswerText, &answer.AnswerValue); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan answer data"})
			return
		}
		answers = append(answers, answer)
	}

	ctx.JSON(http.StatusOK, gin.H{"answers": answers})
}
