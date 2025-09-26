// @title           Dating App API
// @version         1.0
// @description     API documentation for the Dating App backend.
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	docs "github.com/icpinto/dating-app/docs"
	"github.com/icpinto/dating-app/internals/db"
	"github.com/icpinto/dating-app/middlewares"
	"github.com/icpinto/dating-app/services"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	sqlDB, err := db.InitDB()
	if err != nil {
		log.Fatal("Cannot connect to the database:", err)
	}
	defer sqlDB.Close()

	matchServiceURL := os.Getenv("MATCH_SERVICE_URL")

	router := setupRouter(sqlDB, matchServiceURL)

	messagingURL := os.Getenv("MESSAGING_SERVICE_URL")
	if messagingURL == "" {
		messagingURL = "http://localhost:8082"
	}
	worker := services.NewOutboxWorker(sqlDB, messagingURL)
	go worker.Start()

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func setupRouter(sqlDB *sql.DB, matchServiceURL string) *gin.Engine {
	docs.SwaggerInfo.BasePath = "/"

	router := gin.Default()
	router.Static("/uploads", "./uploads")

	// Swagger documentation endpoint.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
	matchService := services.NewMatchService(matchServiceURL)

	router.Use(middlewares.ServiceMiddleware(middlewares.Services{
		UserService:          userService,
		QuestionnaireService: questionnaireService,
		FriendRequestService: friendRequestService,
		ProfileService:       profileService,
		MatchService:         matchService,
	}))

	router.POST("/register", controllers.Register)
	router.POST("/login", controllers.Login)
	router.POST("/signout", middlewares.Authenticate, controllers.SignOut)

	protected := router.Group("/user")
	protected.Use(middlewares.Authenticate)

	protected.POST("/profile", controllers.CreateProfile)
	protected.GET("/profile", controllers.GetProfile)
	protected.GET("/profiles", controllers.GetProfiles)
	protected.GET("/profile/:user_id", controllers.GetUserProfile)
	protected.GET("/matches/:user_id", controllers.GetUserMatches)
	protected.POST("/core-preferences", controllers.SaveCorePreferences)
	protected.PUT("/core-preferences", controllers.UpdateCorePreferences)

	// Allow authenticated users to retrieve profile enumerations via /user/profile/enums
	protected.GET("/profile/enums", controllers.GetProfileEnums)

	router.GET("/profile/enums", controllers.GetProfileEnums)

	protected.POST("/sendRequest", controllers.SendFriendRequest)
	protected.POST("/acceptRequest", controllers.AcceptFriendRequest)
	protected.POST("/rejectRequest", controllers.RejectFriendRequest)
	protected.GET("/requests", controllers.GetPendingRequests)
	protected.GET("/sentRequests", controllers.GetSentRequests)
	protected.GET("/checkReqStatus/:reciver_id", controllers.CheckReqStatus)

	protected.GET("/questionnaire", controllers.GetQuestionnaire)
	protected.POST("/submitQuestionnaire", controllers.SubmitQuestionnaire)
	protected.GET("/questionnaireAnswers", controllers.GetUserAnswers)

	return router
}
