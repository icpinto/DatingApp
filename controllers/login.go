package controllers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
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
	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Hash the password before storing
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Insert the new user into the database
	_, err = db.(*sql.DB).Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", user.Username, user.Email, hashedPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(ctx *gin.Context) {
	var user models.User
	var hashedPassword string

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := services.GetUsepwd(user.Username, db.(*sql.DB))

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Compare hashed password with user input
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.Password))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	// Return the token in the response
	ctx.JSON(http.StatusOK, gin.H{"token": token})

}
