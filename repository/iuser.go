package repository

import (
	"context"
	"program/model"
)

type IUser interface {
	IsUserExists(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	InsertNewUser(ctx context.Context, user *model.User) error
}

var UserRepo IUser
