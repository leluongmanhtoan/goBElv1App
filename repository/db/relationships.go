package db

import (
	"context"
	"errors"
	"fmt"
	"program/model"
	"program/repository"
	"time"
)

type Relationships struct{}

func NewRelationshipsRepo() repository.IRelationships {
	return &Relationships{}
}

func (r *Relationships) AddFollow(ctx context.Context, postFollow *model.Follows) error {
	_, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(postFollow).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Relationships) IsFollowExists(ctx context.Context, followerId, followingId string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		ColumnExpr("1").
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking follow relationship: %w", err)
	}
	return exists, nil
}

func (r *Relationships) IsActiveFollow(ctx context.Context, followerId, followingId string) (bool, error) {
	isActive, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		ColumnExpr("1").
		Where("followerId = ? AND followingId = ? AND isActive = ?", followerId, followingId, true).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking follow status: %w", err)
	}
	return isActive, nil
}

func (r *Relationships) UpdateFollow(ctx context.Context, followerId, followingId string, status bool) error {
	_, err := repository.SqlClientConnection.GetDB().NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isActive = ?", status).
		Set("updatedAt = ?", time.Now()).
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exec(ctx)
	return err
}

func (r *Relationships) UpdateMutualFollowStatus(ctx context.Context, followerId, followingId string, status bool) error {
	_, err := repository.SqlClientConnection.GetDB().NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isMutual = ?", status).
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exec(ctx)
	if err != nil {
		return errors.New("can not update mutual status for follower")
	}
	_, err = repository.SqlClientConnection.GetDB().NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isMutual = ?", status).
		Where("followerId = ? AND followingId = ?", followingId, followerId).
		Exec(ctx)
	if err != nil {
		return errors.New("can not update mutual status for following")
	}
	return nil
}
