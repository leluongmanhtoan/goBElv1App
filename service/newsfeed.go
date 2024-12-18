package service

import (
	"context"
	"program/model"
	"program/repository"
	"time"

	"github.com/google/uuid"
)

type INewsfeed interface {
	PostNewsfeed(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error)
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error)
}

type Newsfeed struct {
}

func NewNewsFeed() INewsfeed {
	return &Newsfeed{}
}

func (s *Newsfeed) PostNewsfeed(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error) {
	newpost := &model.Post{
		PostId:       uuid.NewString(),
		UserId:       user_id,
		Content:      post.Content,
		Privacy:      post.Privacy,
		LikeCount:    0,
		CommentCount: 0,
		ShareCount:   0,
		Deleted:      0,
		CreatedAt:    time.Now(),
	}
	if err := repository.NewsfeedRepo.PostNews(ctx, newpost); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"postId":  newpost.PostId,
		"message": "your post is created",
	}, nil
}

func (s *Newsfeed) GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error) {
	newsfeed, err := repository.NewsfeedRepo.GetNewsfeed(ctx, limit, offset, user_id, true)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"userId": user_id,
		"data":   newsfeed,
		"limit":  limit,
		"offset": offset,
	}, nil
}
