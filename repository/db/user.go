package db

import (
	"context"
	"program/model"
	"program/repository"
)

func NewUser() repository.IUser {
	return &User{}
}

type User struct{}

func (r *User) IsUserExists(ctx context.Context, username string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.User)(nil)).
		ColumnExpr("1").
		Where("username = ?", username).
		Exists(ctx)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (r *User) GetUserByUsername(ctx context.Context, username string) (*model.User, error) {
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

func (r *User) InsertNewUser(ctx context.Context, user *model.User) error {
	res, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(user).
		Exec(ctx)
	if err != nil {

	} else if affected, _ := res.RowsAffected(); affected != 1 {

	}
	return nil
}
