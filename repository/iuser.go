package repository

import (
	"context"
	"program/model"
)

type IUser interface {
	DoesUserExist(ctx context.Context, username string) (bool, error)
	DoesUserProfileExist(ctx context.Context, userID string) (bool, error)
	GetByUserName(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	CreateUserProfle(ctx context.Context, userProfile *model.UserProfile) error
	RetrieveProfileForUser(ctx context.Context, user_id string) (*model.UserProfile, error)
	UpdateProfileForUser(ctx context.Context, user_id string, fields map[string]any) (*model.UserProfile, error)
}

var UserRepo IUser
