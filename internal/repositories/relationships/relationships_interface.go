package relationshipsRepo

import (
	"context"
	"program/internal/model"

	"github.com/uptrace/bun"
)

type IRelationshipsRepo interface {
	GetDBTx(ctx context.Context) (*bun.Tx, error)
	AddFollowTransaction(ctx context.Context, tx *bun.Tx, postFollow *model.Follows) error
	IsFollowExists(ctx context.Context, followerId, followingId string) (bool, error)
	IsActiveFollow(ctx context.Context, followerId, followingId string) (bool, error)
	UpdateFollowTransaction(ctx context.Context, tx *bun.Tx, followerId, followingId string, status bool) error
	UpdateMutualFollowStatusTransaction(ctx context.Context, tx *bun.Tx, followerId, followingId string, status bool) error
	GetFollowList(ctx context.Context, limit, offset int, targetUserId string, isFollowingUser bool) (int, *[]model.FollowerInfo, error)
	NumOfFollowRelationship(ctx context.Context, targetUserId string) (int, int, error)
}
