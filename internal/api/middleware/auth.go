package middleware

import (
	"coretrix/internal/auth"
	"coretrix/internal/user"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	MessageJwtNotFound          = "JWT Token not found"
	MessageJwtInvalidToken      = "invalid token"
	MessageJwtInvalidCredential = "invalid credential"
	UserKey                     = "user"
)

type authHeader struct {
	Token string `header:"Authorization"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func AuthMiddleware(s auth.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := authHeader{}
		err := c.ShouldBindHeader(&h)
		if err != nil {
			abort(c, MessageJwtNotFound)
			return
		}

		parts := strings.Split(h.Token, " ")
		if len(parts) < 2 {
			abort(c, MessageJwtInvalidToken)
			return
		}

		token := parts[1]
		loggedInUser, err := s.GetUser(token)
		if err != nil {
			abort(c, MessageJwtInvalidToken)
			return

		}

		if loggedInUser != nil {
			next(c, loggedInUser)
			return
		}

		abort(c, MessageJwtInvalidCredential)
		return
	}
}

func abort(c *gin.Context, message string) {
	response := Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	}
	c.AbortWithStatusJSON(http.StatusUnauthorized, response)
}

func next(c *gin.Context, u *user.AutheticatedUser) {
	c.Set(UserKey, u)
	c.Next()
}
