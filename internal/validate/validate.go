package validate

import (
	"net/http"
	"program/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func ValidateRequest(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, "can not parse request body")
		return false
	}
	validate := validator.New()
	if err := validate.Struct(obj); err != nil {
		response.ErrorResponse[string](c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}
