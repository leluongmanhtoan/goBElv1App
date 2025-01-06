package middleware

import (
	"net/http"
	"program/internal/response"
	"program/internal/services"

	"github.com/gin-gonic/gin"
)

type IAuthor interface {
	RequestAuthorization() gin.HandlerFunc
	RequestNoRequiredAuthorization() gin.HandlerFunc
}

type AuthorMwd struct {
	authen services.IJwtAuthService
}

var AuthMdw IAuthor

func NewAuthorMdw(auth services.IJwtAuthService) IAuthor {
	return &AuthorMwd{
		authen: auth,
	}
}

func (m *AuthorMwd) RequestAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.ErrorResponse[string](c, http.StatusBadRequest, "authorization field can not be empty")
			c.Abort()
			return
		}
		tokenClaims, err := m.authen.ValidateToken(c, authHeader, false)
		if err != nil {
			response.ErrorResponse[string](c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("userId", tokenClaims.Subject)
	}
}

func (m *AuthorMwd) RequestNoRequiredAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Set("userId", "guest")
			return
		}
		tokenClaims, err := m.authen.ValidateToken(c, authHeader, false)
		if err != nil {
			c.Set("userId", "guest")
			return
		}
		c.Set("userId", tokenClaims.Subject)
	}
}
