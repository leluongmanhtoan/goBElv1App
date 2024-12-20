package userRepo

import (
	"context"
	"fmt"
	"program/internal/database"
	"program/internal/model"
)

type UserRepo struct {
	db database.ISqlConnection
}

func NewUserRepo(db database.ISqlConnection) IUserRepo {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) DoesUserExist(ctx context.Context, username string) (bool, error) {
	exists, err := r.db.GetDB().NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("1").
		Where("username = ?", username).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking user: %w", err)
	}
	return exists, nil
}

func (r *UserRepo) DoesUserProfileExist(ctx context.Context, userID string) (bool, error) {
	exists, err := r.db.GetDB().NewSelect().
		Model((*model.UserProfile)(nil)).
		ColumnExpr("1").
		Where("userId = ?", userID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking user profile: %w", err)
	}
	return exists, nil
}

func (r *UserRepo) GetByUserName(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := r.db.GetDB().NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	_, err := r.db.GetDB().NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) CreateUserProfle(ctx context.Context, userProfile *model.UserProfile) error {
	_, err := r.db.GetDB().NewInsert().
		Model(userProfile).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepo) RetrieveProfileForUser(ctx context.Context, user_id string) (*model.UserProfile, error) {
	profile := new(model.UserProfile)
	err := r.db.GetDB().NewSelect().
		Model(profile).
		Where("userId = ?", user_id).
		Limit(1).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *UserRepo) UpdateProfileForUser(ctx context.Context, user_id string, fields map[string]any) (*model.UserProfile, error) {
	query := r.db.GetDB().NewUpdate().
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
