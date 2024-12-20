package services

import (
	"context"
	"errors"
	"program/internal/model"
	userRepo "program/internal/repositories/user"
	"time"

	"github.com/google/uuid"
)

func NewUserService(repo userRepo.IUserRepo, passHandler *PasswordHandler, auth IJwtAuthService) IUserService {
	return &UserService{
		PassHandler: passHandler,
		Authen:      auth,
		repo:        repo,
	}
}

type IUserService interface {
	Login(ctx context.Context, loginForm model.Login) (*model.LoginResponse, error)
	Register(ctx context.Context, registerForm model.Register) (*model.RegisterResponse, error)
	RefreshToken(ctx context.Context, token string) (*model.RefreshToken, error)
	Logout(accessToken, refreshToken string) (*map[string]string, error)
	CreateUserProfile(ctx context.Context, user_id string, userProfilePost *model.UserProfilePost) (any, error)
	GetUserProfile(ctx context.Context, user_id string) (any, error)
	UpdateUserProfile(ctx context.Context, user_id string, profilePut *model.UserProfilePut) (any, error)
}

type UserService struct {
	PassHandler *PasswordHandler
	Authen      IJwtAuthService
	repo        userRepo.IUserRepo
}

func (s *UserService) Register(ctx context.Context, registerForm model.Register) (*model.RegisterResponse, error) {
	userExisted, err := s.repo.DoesUserExist(ctx, registerForm.Username)
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
	if err = s.repo.CreateUser(ctx, user); err != nil {
		return nil, errors.New("insert new user failed")
	}
	newAccessToken, err := s.Authen.GenerateToken(ctx, user.UserUuid, false)
	if err != nil || newAccessToken == "" {
		return nil, err
	}
	newRefreshToken, err := s.Authen.GenerateToken(ctx, user.UserUuid, true)
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

func (s *UserService) Login(ctx context.Context, loginForm model.Login) (*model.LoginResponse, error) {
	userExisted, err := s.repo.GetByUserName(ctx, loginForm.Username)
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
	newAccessToken, err := s.Authen.GenerateToken(ctx, userExisted.UserUuid, false)
	if err != nil || newAccessToken == "" {
		return nil, err
	}
	newRefreshToken, err := s.Authen.GenerateToken(ctx, userExisted.UserUuid, true)
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

func (s *UserService) Logout(accessToken, refreshToken string) (*map[string]string, error) {
	if err := s.Authen.RevokeSession(accessToken, refreshToken); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"message": "logged out successfully",
	}, nil
}

func (s *UserService) RefreshToken(ctx context.Context, token string) (*model.RefreshToken, error) {
	refreshToken, err := s.Authen.ValidateToken(token, true)
	if err != nil {
		return nil, err
	}
	newAccessToken, err := s.Authen.GenerateToken(ctx, refreshToken.Subject, false)
	if err != nil {
		return nil, err
	}
	return &model.RefreshToken{
		UserId:         refreshToken.Subject,
		NewAccessToken: newAccessToken,
	}, nil
}

func (s *UserService) CreateUserProfile(ctx context.Context, user_id string, userProfilePost *model.UserProfilePost) (any, error) {
	existed, err := s.repo.DoesUserProfileExist(ctx, user_id)
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

	if err := s.repo.CreateUserProfle(ctx, &userProfile); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"message": userProfile.ProfileId,
	}, nil
}

func (s *UserService) GetUserProfile(ctx context.Context, user_id string) (any, error) {
	profile, err := s.repo.RetrieveProfileForUser(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, user_id string, profilePut *model.UserProfilePut) (any, error) {
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
	profile, err := s.repo.UpdateProfileForUser(ctx, user_id, fields)
	if err != nil {
		return nil, err
	}
	return profile, nil
}
