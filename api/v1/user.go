package api

import (
	"program/model"
	"program/service"

	"github.com/gin-gonic/gin"
)

// Create user struct -> container include all func for handler
type User struct {
	userService service.IUser
}

// Func create api collection
func NewUser(r *gin.Engine, userSerivce service.IUser) {
	handler := &User{
		userService: userSerivce,
	}
	Group := r.Group("api/v1")
	{
		Group.POST("auth/register", handler.Register)
		Group.POST("auth/login", handler.Login)
		Group.POST("auth/logout", handler.Logout)
		Group.POST("auth/refresh", handler.RefeshToken)
	}
}

// Method implement struct User
func (h *User) Register(c *gin.Context) {
	registerForm := model.Register{}
	if err := c.BindJSON(&registerForm); err != nil {

	}
	result, _ := h.userService.Register(c, registerForm)
	c.JSON(200, result)
}

func (h *User) Login(c *gin.Context) {
	var loginInfo model.Login
	err := c.BindJSON(&loginInfo)
	if err != nil {

	}
	loginResponse, err := h.userService.Login(c, loginInfo)
	if err != nil {

	}
	c.JSON(200, loginResponse)
}

func (h *User) Logout(c *gin.Context) {
	var request struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BindJSON(&request); err != nil {

	}
	logoutResponse, err := h.userService.Logout(request.AccessToken, request.RefreshToken)
	if err != nil {

	}
	c.JSON(200, logoutResponse)
}

func (h *User) RefeshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refreshToken"`
	}
	if err := c.BindJSON(&request); err != nil {

	}
	refreshResponse, err := h.userService.RefreshToken(request.RefreshToken)
	if err != nil {

	}
	c.JSON(200, refreshResponse)
}
