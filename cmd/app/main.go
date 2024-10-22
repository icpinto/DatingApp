package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/controllers"
	"github.com/icpinto/dating-app/middlewares"
	_ "github.com/lib/pq"
)

func main() {
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

	fmt.Printf("%+v\n", db)

	router := gin.Default()
	// Apply the middleware to inject the db connection into the context
	router.Use(middlewares.DBMiddleware(db))

	router.POST("/register", controllers.Register) // User registration route
	router.POST("/login", controllers.Login)       // User login route

	// Group of routes that require authentication
	protected := router.Group("/user")
	protected.Use(middlewares.Authenticate) // Apply the JWT authentication middleware

	protected.POST("/profile", controllers.CreateProfile) // Update or create profile
	protected.GET("/profile", controllers.GetProfile)     // Get profile
	router.Run(":8080")

}
