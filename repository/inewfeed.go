package repository

import (
	"context"
	"program/model"
)

type INewsfeed interface {
	PostNews(ctx context.Context, post *model.Post) error
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string, isFromFollowing bool) (*[]model.NewsFeed, error)
}

var NewsfeedRepo INewsfeed
