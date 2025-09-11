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
	sqlDB, err := db.InitDB()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}
	defer sqlDB.Close()

	router, chatService := setupRouter(sqlDB)

	go websocket.HandleMessages(chatService)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func setupRouter(sqlDB *sql.DB) (*gin.Engine, *services.ChatService) {
	router := gin.Default()
	router.Static("/uploads", "./uploads")

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	userService := services.NewUserService(sqlDB)
	questionnaireService := services.NewQuestionnaireService(sqlDB)
	friendRequestService := services.NewFriendRequestService(sqlDB)
	profileService := services.NewProfileService(sqlDB)
	chatService := services.NewChatService(sqlDB)

	router.Use(middlewares.ServiceMiddleware(middlewares.Services{
		UserService:          userService,
		QuestionnaireService: questionnaireService,
		FriendRequestService: friendRequestService,
		ProfileService:       profileService,
		ChatService:          chatService,
	}))

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)

	protected := router.Group("/user")
	protected.Use(middlewares.Authenticate)

	protected.POST("/profile", controllers.CreateProfile)
	protected.GET("/profile", controllers.GetProfile)
	protected.GET("/profiles", controllers.GetProfiles)
	protected.GET("/profile/:user_id", controllers.GetUserProfile)

	protected.POST("/sendRequest", controllers.SendFriendRequest)
	protected.POST("/acceptRequest", controllers.AcceptFriendRequest)
	protected.POST("/rejectRequest", controllers.RejectFriendRequest)
	protected.GET("/requests", controllers.GetPendingRequests)
	protected.GET("/checkReqStatus/:reciver_id", controllers.CheckReqStatus)

	protected.GET("/questionnaire", controllers.GetQuestionnaire)
	protected.POST("/submitQuestionnaire", controllers.SubmitQuestionnaire)
	protected.GET("/questionnaireAnswers", controllers.GetUserAnswers)

	router.GET("/ws/:token", middlewares.AuthenticateWS, func(c *gin.Context) {
		websocket.HandleConnections(c)
	})

	protected.GET("/conversations", controllers.GetAllConversations)
	protected.POST("/conversations", controllers.CreateConversation)
	protected.GET("/conversations/:id", controllers.GetChatHistory)

	return router, chatService
}
