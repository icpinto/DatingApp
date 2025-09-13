package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/icpinto/dating-app/models"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("secret") // Secret key used for signing tokens

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

// Function to generate a JWT token
func GenerateToken(userID int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expiration time (24 hours)

	claims := &models.Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	return tokenString, err
}
