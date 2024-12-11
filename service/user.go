package service

import (
	"context"
	"program/model"
	"program/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func NewUser(passHandler *PasswordHandler, auth IAuth) IUser {
	return &User{
		PassHandler: passHandler,
		Authen:      auth,
	}
}

type IUser interface {
	Login(ctx context.Context, loginForm model.Login) (*model.LoginResponse, error)
	Register(ctx context.Context, registerForm model.Register) (*model.RegisterResponse, error)
	RefreshToken(token string) (*model.RefreshToken, error)
	Logout(accessToken, refreshToken string) (*map[string]string, error)
}

type User struct {
	PassHandler *PasswordHandler
	Authen      IAuth
}

func (s *User) Register(ctx context.Context, registerForm model.Register) (*model.RegisterResponse, error) {
	userExisted, err := repository.UserRepo.IsUserExists(ctx, registerForm.Username)
	if !userExisted {

	}
	salt, err := s.PassHandler.GenerateSalt()

	hash, err := s.PassHandler.HashPassword(registerForm.Password, salt)
	user := &model.User{
		UserUuid:  uuid.NewString(),
		Username:  registerForm.Username,
		Salt:      salt,
		Hash:      hash,
		CreatedAt: time.Now(),
		Deleted:   0,
	}
	if err = repository.UserRepo.InsertNewUser(ctx, user); err != nil {

	}
	return &model.RegisterResponse{
		Message:  "User registration successful",
		UserUuid: user.UserUuid,
	}, nil
}

func (s *User) Login(ctx context.Context, loginForm model.Login) (*model.LoginResponse, error) {
	userExisted, err := repository.UserRepo.GetUserByUsername(ctx, loginForm.Username)
	if err != nil || userExisted == nil {

	}
	err = s.PassHandler.ValidatePassword(userExisted.Hash, loginForm.Password, userExisted.Salt)
	if err != nil {

	}
	newAccessToken, err := s.Authen.GenerateToken(userExisted.UserUuid, false)
	if err != nil || newAccessToken == "" {

	}
	newRefreshToken, err := s.Authen.GenerateToken(userExisted.UserUuid, true)
	if err != nil || newRefreshToken == "" {

	}
	//s.Authen.RevokeSession(refreshToken)
	return &model.LoginResponse{
		UserID:       userExisted.UserUuid,
		Username:     userExisted.Username,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *User) Logout(accessToken, refreshToken string) (*map[string]string, error) {
	if err := s.Authen.RevokeSession(accessToken, refreshToken); err != nil {
		return nil, err
	}
	return &map[string]string{
		"message": "logged out successfully",
	}, nil
}

func (s *User) RefreshToken(token string) (*model.RefreshToken, error) {
	refreshToken, err := s.Authen.ValidateToken(token, true)
	if err != nil {

	}
	tokenClaims := refreshToken.Claims.(*jwt.StandardClaims)
	newAccessToken, err := s.Authen.GenerateToken(tokenClaims.Subject, false)
	if err != nil {

	}
	return &model.RefreshToken{
		UserId:         tokenClaims.Subject,
		NewAccessToken: newAccessToken,
	}, nil
}
