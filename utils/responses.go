package utils

import (
	"github.com/gin-gonic/gin"
	"log"
)

// RespondSuccess sends a JSON response with the provided status and payload.
func RespondSuccess(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}

// RespondError logs the error with the provided log message and sends a JSON response with a client-facing message.
// An optional details parameter can be included for additional client context.
func RespondError(ctx *gin.Context, status int, err error, logMsg, clientMsg string, details ...string) {
	if err != nil {
		log.Printf("%s: %v", logMsg, err)
	} else {
		log.Println(logMsg)
	}

	response := gin.H{"error": clientMsg}
	if len(details) > 0 {
		response["details"] = details[0]
	}

	ctx.JSON(status, response)
}
