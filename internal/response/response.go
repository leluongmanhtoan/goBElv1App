package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequest(err error) (int, any) {
	return http.StatusBadRequest, map[string]any{
		"status":  "error -" + http.StatusText(http.StatusBadRequest),
		"message": err.Error(),
	}
}

func Unauthorized(err error) (int, any) {
	return http.StatusUnauthorized, map[string]any{
		"status":  "error -" + http.StatusText(http.StatusUnauthorized),
		"message": err.Error(),
	}
}

func ServiceUnavailableMsg(msg any) (int, any) {
	return http.StatusServiceUnavailable, map[string]any{
		"status":  "error -" + http.StatusText(http.StatusServiceUnavailable),
		"message": msg,
	}
}

type GenericResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type GenericResponsePagination[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data"`
	Limit   int    `json:"limit"`
	Offset  int    `json:"offset"`
}

func SuccessResponse[T any](c *gin.Context, message string, data T) {
	c.JSON(http.StatusOK, GenericResponse[T]{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
	})
}

func SuccessResponseWithPagination[T any](c *gin.Context, limit, offset int, message string, data T) {
	c.JSON(http.StatusOK, GenericResponsePagination[T]{
		Code:    http.StatusOK,
		Message: message,
		Data:    data,
		Limit:   limit,
		Offset:  offset,
	})
}

func ErrorResponse[T any](c *gin.Context, code int, message string) {
	c.JSON(code, GenericResponse[T]{
		Code:    code,
		Message: message,
		Data:    *new(T),
	})
}
