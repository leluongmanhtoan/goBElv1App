package middleware

import (
	"net/http"
	"program/internal/services"

	"github.com/gin-gonic/gin"
)

type IAuthor interface {
	RequestAuthorization() gin.HandlerFunc
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
			c.JSON(
				http.StatusUnauthorized,
				map[string]any{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
		tokenClaims, err := m.authen.ValidateToken(authHeader, false)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				map[string]any{
					"error":   http.StatusText(http.StatusUnauthorized),
					"message": "Invalid or expired token",
				},
			)
			c.Abort()
			return
		}
		c.Set("user_id", tokenClaims.Subject)
	}
}
