package middleware

import (
	"net/http"
	"program/service"

	"github.com/gin-gonic/gin"
)

type IAuthor interface {
	RequestAuthorization() gin.HandlerFunc
}

type AuthorMwd struct {
	authen service.IAuth
}

var AuthMdw IAuthor

func NewAuthorMdw(auth service.IAuth) IAuthor {
	return &AuthorMwd{
		authen: auth,
	}
}

func (m *AuthorMwd) RequestAuthorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Abort()
			return
		}
		_, err := m.authen.ValidateToken(authHeader, false)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				map[string]any{
					"error": http.StatusText(http.StatusUnauthorized),
				},
			)
			c.Abort()
			return
		}
	}
}
