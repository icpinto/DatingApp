package websocket

import (
	"database/sql"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin, in production you'd check domain here
	},
}

var clients = make(map[int]*models.Client) // UserID -> WebSocket connection
var broadcast = make(chan models.ChatMessage)
var mutex = &sync.Mutex{}

// Handle WebSocket connections and register user
func HandleConnections(ctx *gin.Context) {
	// Get the authenticated user ID from JWT context
	/*userID, exists := ctx.Get("userID")
	if !exists {
		log.Println("Missing userID in WebSocket connection")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}*/

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

	// Retrieve the user's ID from the users table
	var userID int
	err := db.(*sql.DB).QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
	}
	defer ws.Close()

	client := &models.Client{Conn: ws, UserID: userID}

	// Register client
	mutex.Lock()
	clients[client.UserID] = client
	mutex.Unlock()

	for {
		var msg models.ChatMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("Error reading message from WebSocket: %v", err)
			mutex.Lock()
			delete(clients, client.UserID)
			mutex.Unlock()
			break
		}

		// Send message to intended recipient via the broadcast channel
		broadcast <- msg
	}
}

func HandleMessages() {

	for {
		msg := <-broadcast
		/*
		   // Save message to the database
		   message := models.ChatMessage{
		       ConversationID: msg.ConversationID,
		       SenderID:       msg.SenderID,
		       Message:        msg.Message,
		   }*/

		_, err := db.DB.Exec(`
        INSERT INTO messages (conversation_id, sender_id, message, created_at)
        VALUES ($1, $2, $3, NOW())`,
			msg.ConversationID, msg.SenderID, msg.Message)

		if err != nil {
			log.Printf("Error saving message to database: %v", err)
			return
		}

		// Find the recipient client (ReceiverID)
		mutex.Lock()
		recipient, ok := clients[msg.ReceiverID]
		mutex.Unlock()

		if ok {
			// Send the message to the recipient if connected
			err := recipient.Conn.WriteJSON(msg)
			if err != nil {
				log.Printf("Error sending message to user %d: %v", msg.ReceiverID, err)
				recipient.Conn.Close()
				mutex.Lock()
				delete(clients, recipient.UserID)
				mutex.Unlock()
			}
		} else {
			log.Printf("User %d is not connected, message could not be delivered", msg.ReceiverID)
		}
	}
}
