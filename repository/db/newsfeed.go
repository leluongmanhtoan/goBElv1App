package db

import (
	"context"
	"database/sql"
	"program/model"
	"program/repository"
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
		Join("JOIN userProfile pf ON pf.userId = f.followingId").
		Join("JOIN posts p ON p.userId = f.followingId")
	if isFromFollowing {
		query.Where("f.followerId = ? OR f.followingId = ? AND deleted = 0 AND p.createdAt >= NOW() - INTERVAL 7 DAY", user_id, user_id)
	} else {
		query.Where("p.userId = ? AND deleted = 0", user_id)
	}
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
