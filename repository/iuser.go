package repository

import (
	"context"
	"program/model"
)

type IUser interface {
	DoesUserExist(ctx context.Context, username string) (exists bool, err error)
	DoesUserProfileExist(ctx context.Context, userID string) (exists bool, err error)
	GetByUserName(ctx context.Context, username string) (*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	CreateUserProfle(ctx context.Context, userProfile *model.UserProfile) error
	RetrieveProfileForUser(ctx context.Context, user_id string) (*model.UserProfile, error)
	UpdateProfileForUser(ctx context.Context, user_id string, profile *model.UserProfilePut) (*model.UserProfile, error)
}

var UserRepo IUser
