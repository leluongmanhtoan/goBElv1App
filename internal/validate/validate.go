package validate

import (
	"program/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func ValidateRequest(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(response.BadRequest(err))
		return false
	}
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		c.JSON(response.BadRequest(err))
		return false
	}
	return true
}
