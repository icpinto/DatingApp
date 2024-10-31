package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/icpinto/dating-app/models"
)

var jwtSecret = []byte("secret") // Secret key used for signing tokens

func Authenticate(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		c.Abort()
		return
	}

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	// Set the username in the context so it can be accessed in the handler
	c.Set("username", claims.Username)

	c.Next() // Proceed to the next middleware or route handler
}

// Middleware to extract and validate JWT from query parameters
func AuthenticateWS(c *gin.Context) {
	// Retrieve token from query parameters

	tokenString := c.Param("token")
	if tokenString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JWT token is required"})
		c.Abort()
		return
	}

	// Parse and validate the JWT token
	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}
	// Set user information in context for further handling
	c.Set("username", claims.Username)
	c.Next()

}
