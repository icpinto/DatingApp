package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/websocket"
	_ "github.com/lib/pq"
)

func main() {

	// Initialize the database
	db.InitDB()
	// Connect to the database
	sqlDB, err := sql.Open("postgres", "user=icpinto dbname=datingapp sslmode=disable password=yourpassword")
	if err != nil {
		log.Fatal("Error opening database:", err)
	}
	defer sqlDB.Close()

	// Ping the database to check if it's reachable
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}

	router := gin.Default()

	// Use the CORS middleware with specific configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// Instantiate services
	userService := services.NewUserService(sqlDB)
	questionnaireService := services.NewQuestionnaireService(sqlDB)
	friendRequestService := services.NewFriendRequestService(sqlDB)
	profileService := services.NewProfileService(sqlDB)
	chatService := services.NewChatService(sqlDB)

	// Apply the middleware to inject services into the context
	router.Use(middlewares.ServiceMiddleware(middlewares.Services{
		UserService:          userService,
		QuestionnaireService: questionnaireService,
		FriendRequestService: friendRequestService,
		ProfileService:       profileService,
		ChatService:          chatService,
	}))

	router.POST("/register", controllers.Register) // User registration route
	router.POST("/login", controllers.Login)

	// Group of routes that require authentication
	protected := router.Group("/user")
	protected.Use(middlewares.Authenticate) // Apply the JWT authentication middleware

	//APIs for user profile
	protected.POST("/profile", controllers.CreateProfile)          // Update or create profile
	protected.GET("/profile", controllers.GetProfile)              // Get profile
	protected.GET("/profiles", controllers.GetProfiles)            // Get Active profiles
	protected.GET("/profile/:user_id", controllers.GetUserProfile) // Get a USer profile

	//APIs for requests
	protected.POST("/sendRequest", controllers.SendFriendRequest)
	protected.POST("/acceptRequest", controllers.AcceptFriendRequest)
	protected.POST("/rejectRequest", controllers.RejectFriendRequest)
	protected.GET("/requests", controllers.GetPendingRequests)
	protected.GET("/checkReqStatus/:reciver_id", controllers.CheckReqStatus)

	//APIs for the Questionnaire
	protected.GET("/questionnaire", controllers.GetQuestionnaire)
	protected.POST("/submitQuestionnaire", controllers.SubmitQuestionnaire)
	protected.GET("/questionnaireAnswers", controllers.GetUserAnswers)

	// WebSocket routes
	router.GET("/ws/:token", middlewares.AuthenticateWS, func(c *gin.Context) {
		websocket.HandleConnections(c)
	})

	// Start handling WebSocket connections in the background
	go websocket.HandleMessages(chatService)

	//APIs for the chat
	protected.GET("/conversations", controllers.GetAllConversations)
	protected.POST("/conversations", controllers.CreateConversation)
	protected.GET("/conversations/:id", controllers.GetChatHistory)

	router.Run(":8080")

}
