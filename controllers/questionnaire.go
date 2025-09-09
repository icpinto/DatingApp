package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
)

func GetQuestionnaire(ctx *gin.Context) {
	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	questions, err := services.GetQuestionnaire(db.(*sql.DB))
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

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	if err := services.SubmitAnswer(db.(*sql.DB), username.(string), answer); err != nil {
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

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	answers, err := services.GetUserAnswers(db.(*sql.DB), username.(string))
	if err != nil {
		log.Printf("GetUserAnswers service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve answers"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"answers": answers})
}
