package db

import (
	"context"
	"fmt"
	"program/model"
	"program/repository"
)

func NewUserRepo() repository.IUser {
	return &User{}
}

type User struct{}

func (r *User) DoesUserExist(ctx context.Context, username string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("1").
		Where("username = ?", username).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	return exists, nil
}

func (r *User) DoesUserProfileExist(ctx context.Context, userID string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.UserProfile)(nil)).
		ColumnExpr("1").
		Where("userId = ?", userID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking user profile: %w", err)
	}
	return exists, nil
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
	_, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *User) CreateUserProfle(ctx context.Context, userProfile *model.UserProfile) error {
	_, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(userProfile).
		Exec(ctx)
	if err != nil {
		return err
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

func (r *User) UpdateProfileForUser(ctx context.Context, user_id string, fields map[string]any) (*model.UserProfile, error) {
	query := repository.SqlClientConnection.GetDB().NewUpdate().
		Model(&model.UserProfile{}).
		Where("userId = ?", user_id)
	for field, value := range fields {
		query.Set(fmt.Sprintf("%s = ?", field), value)
	}
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
