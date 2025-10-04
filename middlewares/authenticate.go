package middlewares

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/icpinto/dating-app/models"
	"github.com/icpinto/dating-app/repositories"
	"github.com/icpinto/dating-app/services"
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

	// Make the user ID available to downstream handlers regardless of account status.
	c.Set("userID", claims.UserID)

	// Retrieve the username based on user ID and set both in the context
	userService := c.MustGet("userService").(*services.UserService)
	username, err := userService.GetUsernameByID(claims.UserID)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) && allowsInactiveAccess(c) {
			// Allow specific routes (such as reactivation) to proceed even when the
			// account is currently inactive. These handlers rely on the user ID from
			// the token and do not require an active username lookup.
			c.Set("username", "")
			c.Next()
			return
		}

		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return
	}

	c.Set("username", username)

	c.Next() // Proceed to the next middleware or route handler
}

func allowsInactiveAccess(c *gin.Context) bool {
	switch c.FullPath() {
	case "/user/reactivate":
		return true
	case "/user/status":
		return true
	default:
		return false
	}
}
