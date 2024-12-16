package service

import (
	"context"
	"errors"
	"program/model"
	"program/repository"
	"time"

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
	CreateUserProfile(ctx context.Context, user_id string, userProfilePost *model.UserProfilePost) (any, error)
	GetUserProfile(ctx context.Context, user_id string) (any, error)
	UpdateUserProfile(ctx context.Context, user_id string, profilePut *model.UserProfilePut) (any, error)
}

type User struct {
	PassHandler *PasswordHandler
	Authen      IAuth
}

func (s *User) Register(ctx context.Context, registerForm model.Register) (*model.RegisterResponse, error) {
	userExisted, err := repository.UserRepo.DoesUserExist(ctx, registerForm.Username)
	if err != nil {
		return nil, errors.New("can not check user existed")
	}
	if userExisted {
		return nil, errors.New("user existed")
	}
	salt, err := s.PassHandler.GenerateSalt()
	if err != nil {
		return nil, errors.New("can not generate salt")
	}
	hash, err := s.PassHandler.HashPassword(registerForm.Password, salt)
	if err != nil {
		return nil, errors.New("can not hash password")
	}
	user := &model.User{
		UserUuid:  uuid.NewString(),
		Username:  registerForm.Username,
		Salt:      salt,
		Hash:      hash,
		CreatedAt: time.Now(),
		Deleted:   0,
	}
	if err = repository.UserRepo.CreateUser(ctx, user); err != nil {
		return nil, errors.New("insert new user failed")
	}
	newAccessToken, err := s.Authen.GenerateToken(user.UserUuid, false)
	if err != nil || newAccessToken == "" {
		return nil, err
	}
	newRefreshToken, err := s.Authen.GenerateToken(user.UserUuid, true)
	if err != nil || newRefreshToken == "" {
		return nil, err
	}
	return &model.RegisterResponse{
		Message:      "User registration successful",
		UserUuid:     user.UserUuid,
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *User) Login(ctx context.Context, loginForm model.Login) (*model.LoginResponse, error) {
	userExisted, err := repository.UserRepo.GetByUserName(ctx, loginForm.Username)
	if err != nil {
		return nil, errors.New("can not get user by username")
	}
	if userExisted == nil {
		return nil, errors.New("username is invalid")
	}
	err = s.PassHandler.ValidatePassword(userExisted.Hash, loginForm.Password, userExisted.Salt)
	if err != nil {
		return nil, errors.New("wrong password")
	}
	newAccessToken, err := s.Authen.GenerateToken(userExisted.UserUuid, false)
	if err != nil || newAccessToken == "" {
		return nil, err
	}
	newRefreshToken, err := s.Authen.GenerateToken(userExisted.UserUuid, true)
	if err != nil || newRefreshToken == "" {
		return nil, err
	}
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
		"status":  "successful",
		"message": "logged out successfully",
	}, nil
}

func (s *User) RefreshToken(token string) (*model.RefreshToken, error) {
	refreshToken, err := s.Authen.ValidateToken(token, true)
	if err != nil {
		return nil, err
	}
	newAccessToken, err := s.Authen.GenerateToken(refreshToken.Subject, false)
	if err != nil {
		return nil, err
	}
	return &model.RefreshToken{
		UserId:         refreshToken.Subject,
		NewAccessToken: newAccessToken,
	}, nil
}

func (s *User) CreateUserProfile(ctx context.Context, user_id string, userProfilePost *model.UserProfilePost) (any, error) {
	existed, err := repository.UserRepo.DoesUserProfileExist(ctx, user_id)
	if err != nil {
		return nil, errors.New("can not check user profile exists")
	}
	if existed {
		return nil, errors.New("user profile existed")
	}
	userProfile := model.UserProfile{
		ProfileId:   uuid.NewString(),
		UserId:      user_id,
		FirstName:   userProfilePost.FirstName,
		LastName:    userProfilePost.LastName,
		Gender:      userProfilePost.Gender,
		Avatar:      userProfilePost.Avatar,
		Address:     userProfilePost.Address,
		Email:       userProfilePost.Email,
		PhoneNumber: userProfilePost.PhoneNumber,
		CreatedAt:   time.Now(),
	}

	if err := repository.UserRepo.CreateUserProfle(ctx, &userProfile); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"message": userProfile.ProfileId,
	}, nil
}

func (s *User) GetUserProfile(ctx context.Context, user_id string) (any, error) {
	profile, err := repository.UserRepo.RetrieveProfileForUser(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *User) UpdateUserProfile(ctx context.Context, user_id string, profilePut *model.UserProfilePut) (any, error) {
	fields := make(map[string]interface{})
	if profilePut.FirstName != "" {
		fields["firstname"] = profilePut.FirstName
	}
	if profilePut.LastName != "" {
		fields["lastname"] = profilePut.LastName
	}
	if profilePut.Gender != nil {
		fields["gender"] = profilePut.Gender
	}
	if profilePut.Avatar != "" {
		fields["avatarUrl"] = profilePut.Avatar
	}
	if profilePut.Address != "" {
		fields["address"] = profilePut.Address
	}
	if profilePut.Email != "" {
		fields["email"] = profilePut.Email
	}
	if profilePut.PhoneNumber != "" {
		fields["phoneNumber"] = profilePut.PhoneNumber
	}
	fields["updatedAt"] = time.Now()
	if len(fields) == 0 {
		return nil, errors.New("no fields to update")
	}
	profile, err := repository.UserRepo.UpdateProfileForUser(ctx, user_id, fields)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
