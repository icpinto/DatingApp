package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
)

// Create a new conversation between two users
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
	log.Println(req)
	_, err := db.(*sql.DB).Exec(`INSERT INTO conversations (user1_id, user2_id, created_at) VALUES ($1, $2, NOW())`,
		req.User1ID, req.User2ID)

	if err != nil {
		log.Println("id", req.User2ID)
		log.Println(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create conversation"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Conversation created"})
}

// Retrieve all conversations
func GetAllConversations(ctx *gin.Context) {
	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	// Query to get all conversations
	rows, err := db.(*sql.DB).Query(`
		SELECT id, user1_id, user2_id, created_at 
		FROM conversations
	`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve conversations"})
		return
	}
	defer rows.Close()

	// Slice to hold all profiles
	var conversations []models.Conversation

	// Iterate through the rows and scan each row into a conversation struct
	for rows.Next() {
		var conversation models.Conversation
		err := rows.Scan(
			&conversation.ID, &conversation.User1ID, &conversation.User2ID, &conversation.CreatedAt,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan conversation"})
			return
		}
		conversations = append(conversations, conversation)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error during rows iteration"})
		return
	}

	// Return the list of conversations as JSON
	ctx.JSON(http.StatusOK, conversations)
}

// Fetch chat history for a conversation
func GetChatHistory(ctx *gin.Context) {
	conversationID := ctx.Param("id")

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	var messages []models.ChatMessage
	rows, err := db.(*sql.DB).Query(`SELECT conversation_id, sender_id, message, created_at FROM messages WHERE conversation_id = $1 ORDER BY created_at ASC`, conversationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch chat history"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var msg models.ChatMessage
		if err := rows.Scan(&msg.ConversationID, &msg.SenderID, &msg.Message, &msg.CreatedAt); err != nil {
			log.Println(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading message"})
			return
		}
		messages = append(messages, msg)
	}

	ctx.JSON(http.StatusOK, messages)
}
