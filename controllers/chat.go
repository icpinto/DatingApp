package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
)

func CreateConversation(ctx *gin.Context) {
	chatService := ctx.MustGet("chatService").(*services.ChatService)

	var req struct {
		User1ID int `json:"user1_id"`
		User2ID int `json:"user2_id"`
	}

	if err := ctx.BindJSON(&req); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "CreateConversation bind error", "Invalid input")
		return
	}

	if err := chatService.CreateConversation(req.User1ID, req.User2ID); err != nil {
		logMsg := fmt.Sprintf("CreateConversation service error for users %d and %d", req.User1ID, req.User2ID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to create conversation")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "Conversation created"})
}

func GetAllConversations(ctx *gin.Context) {
	username, exists := ctx.Get("username")
	if !exists {
		utils.RespondError(ctx, http.StatusUnauthorized, nil, "GetAllConversations unauthorized", "Unauthorized")
		return
	}

	chatService := ctx.MustGet("chatService").(*services.ChatService)

	conversations, err := chatService.GetAllConversations(username.(string))
	if err != nil {
		logMsg := fmt.Sprintf("GetAllConversations service error for %s", username.(string))
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Failed to retrieve conversations")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, conversations)
}

func GetChatHistory(ctx *gin.Context) {
	conversationID := ctx.Param("id")

	chatService := ctx.MustGet("chatService").(*services.ChatService)

	messages, err := chatService.GetChatHistory(conversationID)
	if err != nil {
		logMsg := fmt.Sprintf("GetChatHistory service error for conversation %s", conversationID)
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "Could not fetch chat history")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, messages)
}
