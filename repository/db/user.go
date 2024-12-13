package db

import (
	"context"
	"errors"
	"fmt"
	"program/model"
	"program/repository"
	"time"
)

func NewUser() repository.IUser {
	return &User{}
}

type User struct{}

func (r *User) DoesUserExist(ctx context.Context, username string) (exists bool, err error) {
	exists, err = repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("1").
		Where("username = ?", username).
		Exists(ctx)
	if err != nil {
		return
	}
	return
}

func (r *User) DoesUserProfileExist(ctx context.Context, userID string) (exists bool, err error) {
	exists, err = repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.UserProfile)(nil)).
		ColumnExpr("1").
		Where("userId = ?", userID).
		Exists(ctx)
	if err != nil {
		return
	}
	return
}

func (r *User) GetByUserName(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := repository.SqlClientConnection.GetDB().NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *User) CreateUser(ctx context.Context, user *model.User) error {
	res, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected != 1 {
		return errors.New("insert new user failed")
	}
	return nil
}

func (r *User) CreateUserProfle(ctx context.Context, userProfile *model.UserProfile) error {
	res, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(userProfile).
		Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := res.RowsAffected(); affected != 1 {
		return errors.New("insert new user profile failed")
	}
	return nil
}

func (r *User) RetrieveProfileForUser(ctx context.Context, user_id string) (*model.UserProfile, error) {
	profile := new(model.UserProfile)
	err := repository.SqlClientConnection.GetDB().NewSelect().
		Model(profile).
		Where("userId = ?", user_id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *User) UpdateProfileForUser(ctx context.Context, user_id string, profile *model.UserProfilePut) (*model.UserProfile, error) {
	query := repository.SqlClientConnection.GetDB().NewUpdate().
		Model(&model.UserProfile{}).
		Where("userId = ?", user_id)
	if profile.FirstName != "" {
		query.Set("firstname = ?", profile.FirstName)
	}
	if profile.LastName != "" {
		query.Set("lastname = ?", profile.LastName)
	}
	if profile.Gender != nil {
		query.Set("gender = ?", profile.Gender)
	}
	if profile.Avatar != "" {
		query.Set("avatarUrl = ?", profile.Avatar)
	}
	if profile.Address != "" {
		query.Set("address = ?", profile.Address)
	}
	if profile.Email != "" {
		query.Set("email = ?", profile.Email)
	}
	if profile.PhoneNumber != "" {
		query.Set("phoneNumber = ?", profile.PhoneNumber)
	}
	query.Set("updatedAt = ?", time.Now())
	_, err := query.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %v", err)
	}
	profileUpdated, err := r.RetrieveProfileForUser(ctx, user_id)
	if err != nil {
		return nil, err
	}
	return profileUpdated, nil
}
