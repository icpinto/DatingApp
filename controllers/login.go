package controllers

import (
	"errors"
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
		utils.RespondError(ctx, http.StatusBadRequest, err, "Register bind error", "Invalid input")
		return
	}

	if err := userService.RegisterUser(user); err != nil {
		if errors.Is(err, repositories.ErrDuplicateUser) {
			logMsg := "Register duplicate user"
			utils.RespondError(ctx, http.StatusConflict, err, logMsg, "user already exists")
			return
		}
		logMsg := "Register service error"
		utils.RespondError(ctx, http.StatusInternalServerError, err, logMsg, "User creation failed")
		return
	}

	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(ctx *gin.Context) {
	userService := ctx.MustGet("userService").(*services.UserService)

	var user models.User
	var hashedPassword string
	var userId int

	if err := ctx.BindJSON(&user); err != nil {
		utils.RespondError(ctx, http.StatusBadRequest, err, "Login bind error", "Invalid input")
		return
	}

	hashedPassword, userId, err := userService.GetUsepwd(user.Username)

	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			utils.RespondError(ctx, http.StatusNotFound, err, "Login user not found", "user not found")
			return
		}
		utils.RespondError(ctx, http.StatusUnauthorized, err, "Login service error", "Invalid credentials")
		return
	}

	// Compare hashed password with user input
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		utils.RespondError(ctx, http.StatusUnauthorized, err, "Login password mismatch", "Invalid credentials")
		return
	}

	// Generate JWT token using user ID
	token, err := utils.GenerateToken(userId)
	if err != nil {
		utils.RespondError(ctx, http.StatusInternalServerError, err, "Login token generation error", "Token generation failed")
		return
	}

	// Return the token in the response
	utils.RespondSuccess(ctx, http.StatusOK, gin.H{"token": token, "user_id": userId})

}
