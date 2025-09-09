package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/services"
)

func CreateConversation(ctx *gin.Context) {
	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	var req struct {
		User1ID int `json:"user1_id"`
		User2ID int `json:"user2_id"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := services.CreateConversation(db.(*sql.DB), req.User1ID, req.User2ID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Conversation created"})
}

func GetAllConversations(ctx *gin.Context) {
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

	conversations, err := services.GetAllConversations(db.(*sql.DB), username.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conversations"})
		return
	}

	ctx.JSON(http.StatusOK, conversations)
}

func GetChatHistory(ctx *gin.Context) {
	conversationID := ctx.Param("id")

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	messages, err := services.GetChatHistory(db.(*sql.DB), conversationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch chat history"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
