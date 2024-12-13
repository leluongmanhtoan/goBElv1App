package response

import "net/http"

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
