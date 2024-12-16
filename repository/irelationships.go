package repository

import (
	"context"
	"program/model"
)

type IRelationships interface {
	AddFollow(ctx context.Context, postFollow *model.Follows) error
	IsFollowExists(ctx context.Context, followerId, followingId string) (bool, error)
	IsActiveFollow(ctx context.Context, followerId, followingId string) (bool, error)
	UpdateFollow(ctx context.Context, followerId, followingId string, status bool) error
	UpdateMutualFollowStatus(ctx context.Context, followerId, followingId string, status bool) error
}

var RelationshipsRepo IRelationships
