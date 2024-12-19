package db

import (
	"context"
	"database/sql"
	"fmt"
	"program/model"
	"program/repository"
	"time"

	"github.com/uptrace/bun"
)

type Newsfeed struct{}

func NewNewsFeedRepo() repository.INewsfeed {
	return &Newsfeed{}
}

func (r *Newsfeed) PostNews(ctx context.Context, post *model.Post) error {
	_, err := repository.SqlClientConnection.GetDB().NewInsert().
		Model(post).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Newsfeed) GetNewsfeed(ctx context.Context, limit, offset int, user_id string, isFromFollowing bool) (*[]model.NewsFeed, error) {
	newsfeed := new([]model.NewsFeed)
	query := repository.SqlClientConnection.GetDB().NewSelect().Distinct().
		Column(
			"pf.avatarUrl",
			"pf.firstname",
			"pf.lastname",
			"p.content",
			"p.privacy",
			"p.likeCount",
			"p.commentCount",
			"p.shareCount",
			"p.createdAt",
			"p.updatedAt").
		TableExpr("follows as f").
		Join("JOIN profiles pf ON pf.userId = f.followingId").
		Join("JOIN posts p ON p.userId = f.followingId")
	if isFromFollowing {
		query.Where("f.followerId = ? OR f.followingId = ? AND deleted = 0 AND p.createdAt >= NOW() - INTERVAL 7 DAY", user_id, user_id)
	} else {
		query.Where("p.userId = ? AND deleted = 0", user_id)
	}
	query.Order("p.createdAt DESC")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	err := query.Scan(ctx, newsfeed)
	if err != nil {
		if err == sql.ErrNoRows {
			return newsfeed, nil
		}
		return nil, err
	}
	return newsfeed, nil
}

func (r *Newsfeed) CreateLike(ctx context.Context, tx *bun.Tx, like *model.Like) error {
	_, err := tx.NewInsert().
		Model(like).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Newsfeed) IncreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error {
	_, err := tx.NewUpdate().Model((*model.Post)(nil)).
		Set("likeCount = likeCount + 1").
		Where("postId = ?", postId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Newsfeed) DecreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error {
	_, err := tx.NewUpdate().Model((*model.Post)(nil)).
		Set("likeCount = likeCount - 1").
		Where("postId = ?", postId).
		Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (r *Newsfeed) IsLikeExisted(ctx context.Context, postId, userId string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.Like)(nil)).
		ColumnExpr("1").
		Where("postId = ?", postId).
		Where("userId = ?", userId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking like post: %w", err)
	}
	return exists, nil
}

func (r *Newsfeed) IsPostExisted(ctx context.Context, postId string) (bool, error) {
	exists, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.Post)(nil)).
		ColumnExpr("1").
		Where("postId = ?", postId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking exist post: %w", err)
	}
	return exists, nil
}

func (r *Newsfeed) IsActiveLike(ctx context.Context, post_id, user_id string) (bool, error) {
	isActive, err := repository.SqlClientConnection.GetDB().NewSelect().
		Model((*model.Like)(nil)).
		ColumnExpr("1").
		Where("userId = ? AND postId = ? AND isActive = ?", user_id, post_id, true).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking like status: %w", err)
	}
	return isActive, nil
}

func (r *Newsfeed) UpdateLikeTransaction(ctx context.Context, tx *bun.Tx, user_id, post_id string, status bool) error {
	_, err := tx.NewUpdate().
		Model((*model.Like)(nil)).
		Set("isActive = ?", status).
		Set("updatedAt = ?", time.Now()).
		Where("userId = ? AND postId = ?", user_id, post_id).
		Exec(ctx)
	return err
}

func (r *Newsfeed) GetLikers(ctx context.Context, limit, offset int, post_id string) (*[]model.LikerInfo, error) {
	likers := new([]model.LikerInfo)
	query := repository.SqlClientConnection.GetDB().NewSelect().
		Column("p.profileId", "p.firstname", "p.lastname", "p.avatarUrl").
		TableExpr("likes as l").
		Join("JOIN profiles p ON p.userId=l.userId").
		Where("postId = ? AND isActive = 1", post_id).
		Order("p.lastname ASC")

	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	err := query.Scan(ctx, likers)
	if err != nil {
		if err == sql.ErrNoRows {
			return likers, nil
		}
		return nil, err
	}
	return likers, nil
}
