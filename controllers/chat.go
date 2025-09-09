package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/services"
)

func CreateConversation(ctx *gin.Context) {
	chatService := ctx.MustGet("chatService").(*services.ChatService)

	var req struct {
		User1ID int `json:"user1_id"`
		User2ID int `json:"user2_id"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		log.Printf("CreateConversation bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := chatService.CreateConversation(req.User1ID, req.User2ID); err != nil {
		log.Printf("CreateConversation service error for users %d and %d: %v", req.User1ID, req.User2ID, err)
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

	chatService := ctx.MustGet("chatService").(*services.ChatService)

	conversations, err := chatService.GetAllConversations(username.(string))
	if err != nil {
		log.Printf("GetAllConversations service error for %s: %v", username.(string), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conversations"})
		return
	}

	ctx.JSON(http.StatusOK, conversations)
}

func GetChatHistory(ctx *gin.Context) {
	conversationID := ctx.Param("id")

	chatService := ctx.MustGet("chatService").(*services.ChatService)

	messages, err := chatService.GetChatHistory(conversationID)
	if err != nil {
		log.Printf("GetChatHistory service error for conversation %s: %v", conversationID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch chat history"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}
