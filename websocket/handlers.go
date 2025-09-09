package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/services"
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

	userService := ctx.MustGet("userService").(*services.UserService)
	userID, err := userService.GetUserIDByUsername(username.(string))
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

func HandleMessages(chatService *services.ChatService) {

	for {
		msg := <-broadcast

		if err := chatService.SaveMessage(msg); err != nil {
			log.Printf("Error saving message to database: %v", err)
			continue
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
