package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/icpinto/dating-app/services"
)

type Services struct {
	UserService          *services.UserService
	QuestionnaireService *services.QuestionnaireService
	FriendRequestService *services.FriendRequestService
	ProfileService       *services.ProfileService
	MatchService         *services.MatchService
}

func ServiceMiddleware(s Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("userService", s.UserService)
		c.Set("questionnaireService", s.QuestionnaireService)
		c.Set("friendRequestService", s.FriendRequestService)
		c.Set("profileService", s.ProfileService)
		c.Set("matchService", s.MatchService)
		c.Next()
	}
}
