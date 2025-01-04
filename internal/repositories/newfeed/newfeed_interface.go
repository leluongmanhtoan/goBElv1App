package newsfeedRepo

import (
	"context"
	"program/internal/model"

	"github.com/uptrace/bun"
)

type INewsfeedRepo interface {
	GetDBTx(ctx context.Context) (*bun.Tx, error)
	CreatePost(ctx context.Context, post *model.Post) (*model.NewsFeed, error)
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (*[]model.NewsFeed, error)
	CreateLike(ctx context.Context, tx *bun.Tx, like *model.Like) error
	IncreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error
	DecreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error
	IsLikeExisted(ctx context.Context, postId, userId string) (bool, error)
	IsActiveLike(ctx context.Context, post_id, user_id string) (bool, error)
	UpdateLikeTransaction(ctx context.Context, tx *bun.Tx, user_id, post_id string, status bool) error
	IsPostExisted(ctx context.Context, postId string) (bool, error)
	GetLikers(ctx context.Context, limit, offset int, post_id string) (*[]model.LikerInfo, error)
	CheckPublicPrivacyPermission(ctx context.Context, postId string) error
	CheckFriendPrivacyPermission(ctx context.Context, userId string, postId string) error
	CreateComment(ctx context.Context, commentPost *model.Comment) (*model.CommentInfo, error)
	GetComments(ctx context.Context, limit, offset int, postId string) (*[]model.CommentInfo, error)
	IsOwnPost(ctx context.Context, post_id, user_id string) (bool, error)
	SetOwnerLikedStatus(ctx context.Context, tx *bun.Tx, postId string, status bool) error
	PutComment(ctx context.Context, commentId string, content string) error
}
