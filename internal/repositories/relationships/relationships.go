package relationshipsRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"program/internal/database"
	"program/internal/model"
	"time"

	"github.com/uptrace/bun"
)

type RelationshipsRepo struct {
	db database.ISqlConnection
}

func NewRelationshipsRepo(db database.ISqlConnection) IRelationshipsRepo {
	return &RelationshipsRepo{db: db}
}

func (r *RelationshipsRepo) GetDBTx(ctx context.Context) (*bun.Tx, error) {
	tx, err := r.db.GetDB().BeginTx(ctx, nil)
	return &tx, err
}

func (r *RelationshipsRepo) GetFollowList(ctx context.Context, limit, offset int, targetUserId string, isFollowingUser bool) (int, *[]model.FollowerInfo, error) {
	//var followers []model.FollowerInfo
	follow := new([]model.FollowerInfo)
	query := r.db.GetDB().NewSelect().
		Column("p.profileId", "p.firstname", "p.lastname", "p.avatarUrl").
		TableExpr("follows as f")

	if isFollowingUser {
		query.Join("JOIN profiles p ON p.userId = f.followerId")
		query.Where("f.followingId = ?", targetUserId)
	} else {
		query.Join("JOIN profiles p ON p.userId = f.followingId")
		query.Where("f.followerId = ?", targetUserId)
	}
	query.Where("isActive = 1").
		Order("p.lastname ASC")

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	total, err := query.ScanAndCount(ctx, follow)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, follow, nil
		}
		return 0, nil, err
	}
	return total, follow, nil

}

func (r *RelationshipsRepo) AddFollowTransaction(ctx context.Context, tx *bun.Tx, postFollow *model.Follows) error {
	_, err := tx.NewInsert().
		Model(postFollow).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *RelationshipsRepo) IsFollowExists(ctx context.Context, followerId, followingId string) (bool, error) {
	exists, err := r.db.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		ColumnExpr("1").
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking follow relationship: %w", err)
	}
	return exists, nil
}

func (r *RelationshipsRepo) IsActiveFollow(ctx context.Context, followerId, followingId string) (bool, error) {
	isActive, err := r.db.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		ColumnExpr("1").
		Where("followerId = ? AND followingId = ? AND isActive = ?", followerId, followingId, true).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking follow status: %w", err)
	}
	return isActive, nil
}

func (r *RelationshipsRepo) UpdateFollowTransaction(ctx context.Context, tx *bun.Tx, followerId, followingId string, status bool) error {
	_, err := tx.NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isActive = ?", status).
		Set("updatedAt = ?", time.Now()).
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exec(ctx)
	return err
}

func (r *RelationshipsRepo) UpdateMutualFollowStatusTransaction(ctx context.Context, tx *bun.Tx, followerId, followingId string, status bool) error {
	_, err := tx.NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isMutual = ?", status).
		Where("followerId = ? AND followingId = ?", followerId, followingId).
		Exec(ctx)
	if err != nil {
		return errors.New("can not update mutual status for follower")
	}
	_, err = tx.NewUpdate().
		Model((*model.Follows)(nil)).
		Set("isMutual = ?", status).
		Where("followerId = ? AND followingId = ?", followingId, followerId).
		Exec(ctx)
	if err != nil {
		return errors.New("can not update mutual status for following")
	}
	return nil
}

func (r *RelationshipsRepo) NumOfFollowRelationship(ctx context.Context, targetUserId string) (int, int, error) {
	var followerCount = 0
	var followingCount = 0
	err := r.db.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		Where("followingId = ?", targetUserId).
		Where("isActive = 1").
		ColumnExpr("COUNT(*)").
		Scan(ctx, &followerCount)
	if err != nil {
		return 0, 0, err
	}

	err = r.db.GetDB().NewSelect().
		Model((*model.Follows)(nil)).
		Where("followerId = ?", targetUserId).
		ColumnExpr("COUNT(*)").
		Scan(ctx, &followingCount)
	if err != nil {
		return 0, 0, err
	}
	return followerCount, followingCount, nil
}
