package middlewares

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

// Middleware to inject the database into the context
func DBMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set the database in the context
		c.Set("db", db)
		c.Next()
	}
}
