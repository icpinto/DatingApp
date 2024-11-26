package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/models"
	"github.com/lib/pq"
)

func CreateProfile(ctx *gin.Context) {
	// Get the authenticated user's username (from JWT)

	username, exists := ctx.Get("username")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	log.Println("Name:", username)

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

	// Bind the incoming JSON to the Profile struct
	var profile models.Profile
	if err := ctx.BindJSON(&profile); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	profile.UserID = userID

	// Insert or update the profile
	_, err = db.(*sql.DB).Exec(`
        INSERT INTO profiles (user_id, bio, gender, date_of_birth, location, interests)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (user_id) 
        DO UPDATE SET bio = $2, gender = $3, date_of_birth = $4, location = $5, interests = $6, updated_at = NOW()`,
		userID, profile.Bio, profile.Gender, profile.DateOfBirth, profile.Location, pq.Array(profile.Interests))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})

}

func GetProfile(ctx *gin.Context) {
	// Get the authenticated user's username (from JWT)
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

	// Retrieve the user's profile
	var profile models.Profile
	err = db.(*sql.DB).QueryRow(`
        SELECT id, user_id, bio, gender, date_of_birth, location, interests, created_at, updated_at
        FROM profiles WHERE user_id = $1`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Bio, &profile.Gender,
		&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests),
		&profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}

func GetProfiles(ctx *gin.Context) {
	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	// Query to get all profiles
	rows, err := db.(*sql.DB).Query(`
		SELECT id, user_id, bio, gender, date_of_birth, location, interests, created_at, updated_at
		FROM profiles
	`)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profiles"})
		return
	}
	defer rows.Close()

	// Slice to hold all profiles
	var profiles []models.Profile

	// Iterate through the rows and scan each row into a Profile struct
	for rows.Next() {
		var profile models.Profile
		err := rows.Scan(
			&profile.ID, &profile.UserID, &profile.Bio, &profile.Gender,
			&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests),
			&profile.CreatedAt, &profile.UpdatedAt,
		)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan profile"})
			return
		}
		profiles = append(profiles, profile)
	}

	// Check for any errors encountered during iteration
	if err = rows.Err(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error during rows iteration"})
		return
	}

	// Return the list of profiles as JSON
	ctx.JSON(http.StatusOK, profiles)
}

func GetUserProfile(ctx *gin.Context) {
	userID := ctx.Param("user_id")

	db, exists := ctx.Get("db")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "database not found"})
		return
	}

	var profile models.Profile
	err := db.(*sql.DB).QueryRow(`
		SELECT id, user_id, bio, gender, date_of_birth, location, interests, created_at, updated_at
		FROM profiles WHERE user_id = $1`, userID).Scan(
		&profile.ID, &profile.UserID, &profile.Bio, &profile.Gender,
		&profile.DateOfBirth, &profile.Location, pq.Array(&profile.Interests),
		&profile.CreatedAt, &profile.UpdatedAt)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve profile"})
		return
	}

	ctx.JSON(http.StatusOK, profile)
}
