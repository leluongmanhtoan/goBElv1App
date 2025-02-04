package api

import (
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"program/internal/middleware"
	"program/internal/model"
	"program/internal/response"
	"program/internal/services"
	"program/internal/validate"
	"time"

	"github.com/gin-gonic/gin"
)

// Create user struct -> container include all func for handler
type User struct {
	userService services.IUserService
}

// Func create api collection
func NewUserAPI(engine *gin.Engine, userSerivce services.IUserService) {
	handler := &User{
		userService: userSerivce,
	}
	Group := engine.Group("api/v1")
	{
		//Authen
		Group.POST("auth/register", handler.Register)
		Group.POST("auth/login", handler.Login)
		Group.POST("auth/logout", handler.Logout)
		Group.POST("auth/refresh", handler.RefeshToken)
		Group.POST("auth/validate", middleware.AuthMdw.RequestAuthorization(), func(c *gin.Context) {
			user_id, existed := c.Get("userId")
			if !existed {
				c.JSON(response.BadRequest(errors.New("user_id not found")))
				return
			}
			c.JSON(http.StatusOK, user_id.(string))
		})

		//Profile
		Group.GET("user/profile", middleware.AuthMdw.RequestAuthorization(), handler.GetUserProfile)
		Group.POST("user/profile", middleware.AuthMdw.RequestAuthorization(), handler.NewUserProfile)
		Group.PATCH("user/profile", middleware.AuthMdw.RequestAuthorization(), handler.EditUserProfile)
		Group.POST("user/profile/avatar", middleware.AuthMdw.RequestAuthorization(), handler.UploadAvatar)
	}
}

// Method implement struct User
func (h *User) Register(c *gin.Context) {
	registerForm := model.Register{}
	if !validate.ValidateRequest(c, &registerForm) {
		return
	}
	result, err := h.userService.Register(c, registerForm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *User) Login(c *gin.Context) {
	var loginInfo model.Login
	if !validate.ValidateRequest(c, &loginInfo) {
		return
	}
	loginResponse, err := h.userService.Login(c, loginInfo)
	if err != nil {
		c.JSON(response.Unauthorized(err))
		return
	}
	c.JSON(http.StatusOK, loginResponse)
}

func (h *User) Logout(c *gin.Context) {
	var request struct {
		AccessToken  string `json:"accessToken" validate:"required"`
		RefreshToken string `json:"refreshToken"`
	}
	if !validate.ValidateRequest(c, &request) {
		return
	}
	logoutResponse, err := h.userService.Logout(c, request.AccessToken, request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.SetCookie("refresh_token", "", -1, "/", "", true, true)
	c.JSON(http.StatusOK, logoutResponse)
}

func (h *User) RefeshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}
	if !validate.ValidateRequest(c, &request) {
		return
	}
	refreshResponse, err := h.userService.RefreshToken(c, request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, refreshResponse)
}

func (h *User) NewUserProfile(c *gin.Context) {
	var userProfilePost model.UserProfilePost
	if !validate.ValidateRequest(c, &userProfilePost) {
		return
	}
	user_id, existed := c.Get("userId")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user_id not found")))
		return
	}
	insertResponse, err := h.userService.CreateUserProfile(c, user_id.(string), &userProfilePost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]any{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, insertResponse)
}

func (h *User) EditUserProfile(c *gin.Context) {
	user_id, existed := c.Get("userId")
	if !existed {
		c.JSON(response.BadRequest(errors.New("user id not found")))
		return
	}

	var userProfilePut model.UserProfilePut
	if !validate.ValidateRequest(c, &userProfilePut) {
		return
	}

	profileResponse, err := h.userService.UpdateUserProfile(c, user_id.(string), &userProfilePut)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, profileResponse)

}

func (h *User) GetUserProfile(c *gin.Context) {
	user_id, existed := c.Get("userId")
	if !existed {
		c.JSON(response.BadRequest(errors.New("userid not found")))
		return
	}
	userProfile, err := h.userService.GetUserProfile(c, user_id.(string))
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, userProfile)

}

func (h *User) UploadAvatar(c *gin.Context) {
	user_id, existed := c.Get("userId")
	if !existed {
		c.JSON(response.BadRequest(errors.New("userid not found")))
		return
	}
	file, err := c.FormFile("userAvatar")

	if err != nil {
		c.JSON(response.BadRequest(errors.New("can not get file from request")))
		return
	}

	if file != nil {
		mimeType := file.Header.Get("Content-Type")
		if mimeType != "image/jpeg" && mimeType != "image/png" {
			c.JSON(response.BadRequest(errors.New("invalid file type. Only JPEG and PNG allowed.")))
			return
		}
	}
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("avt_%s_%s%s", user_id, time.Now().Format("2006-01-02_15-04-05"), ext)
	avtPath, err := h.userService.UploadAvatar(c, file, filename)
	if err != nil {
		c.JSON(response.ServiceUnavailableMsg(err.Error()))
		return
	}

	response.SuccessResponse(c, "upload avatar successfully", avtPath)
}
