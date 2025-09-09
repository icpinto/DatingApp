package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/services"
	"github.com/icpinto/dating-app/utils"
	"golang.org/x/crypto/bcrypt"
)

/*
var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("postgres", "user=icpinto dbname=datingapp sslmode=disable password=yourpassword")
	if err != nil {
		log.Fatal(err)
	}
}*/

func Register(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)

	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		log.Printf("Register bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := userService.RegisterUser(user); err != nil {
		if errors.Is(err, repositories.ErrDuplicateUser) {
			log.Printf("Register duplicate user %s: %v", user.Username, err)
			ctx.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
			return
		}
		log.Printf("Register service error for %s: %v", user.Username, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)

	var user models.User
	var hashedPassword string
	var userId int

	if err := ctx.BindJSON(&user); err != nil {
		log.Printf("Login bind error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, userId, err := userService.GetUsepwd(user.Username)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			log.Printf("Login user not found %s: %v", user.Username, err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		log.Printf("Login service error for %s: %v", user.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare hashed password with user input
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		log.Printf("Login password mismatch for %s: %v", user.Username, err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		log.Printf("Login token generation error for %s: %v", user.Username, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Return the token in the response
	ctx.JSON(http.StatusOK, gin.H{"token": token, "user_id": userId})

}
