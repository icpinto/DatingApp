package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/websocket"
	_ "github.com/lib/pq"
)

func main() {

	// Initialize the database
	db.InitDB()
	// Connect to the database
	db, err := sql.Open("postgres", "user=icpinto dbname=datingapp sslmode=disable password=yourpassword")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer db.Close()

	// Ping the database to check if it's reachable
	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	router := gin.Default()
	// Apply the middleware to inject the db connection into the context
	router.Use(middlewares.DBMiddleware(db))

	router.POST("/register", controllers.Register) // User registration route
	router.POST("/login", controllers.Login)

	// Start handling WebSocket connections in the background
	//go websocket.HandleMessages()

	// Group of routes that require authentication
	protected := router.Group("/user")
	protected.Use(middlewares.Authenticate) // Apply the JWT authentication middleware

	protected.POST("/profile", controllers.CreateProfile) // Update or create profile
	protected.GET("/profile", controllers.GetProfile)     // Get profile

	// WebSocket routes
	protected.GET("/ws", func(c *gin.Context) {
		websocket.HandleConnections(c)
	})

	// Start handling WebSocket connections in the background
	go websocket.HandleMessages()

	// Chat API routes
	protected.POST("/conversations", controllers.CreateConversation)
	protected.GET("/conversations/:id/messages", controllers.GetChatHistory)

	router.Run(":8080")

}
