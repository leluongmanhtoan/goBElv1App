package service

import (
	"context"
	"errors"
	"fmt"
	"program/model"
	"program/repository"
	"time"
)

type IRelationships interface {
	FollowUser(ctx context.Context, followerId, followingId string) (any, error)
	UnFollowUser(ctx context.Context, followerId, followingId string) (any, error)
}

type Relationships struct {
}

func NewRelationships() IRelationships {
	return &Relationships{}
}

func (s *Relationships) FollowUser(ctx context.Context, followerId, followingId string) (any, error) {
	if followerId == followingId {
		return nil, errors.New("can not follow yourself")
	}
	followRelationship := &model.Follows{
		FollowerId:  followerId,
		FollowingId: followingId,
		IsActive:    true,
		IsMutual:    false,
		CreatedAt:   time.Now(),
	}
	exists, err := repository.RelationshipsRepo.IsFollowExists(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}
	if exists {
		isActive, err := repository.RelationshipsRepo.IsActiveFollow(ctx, followerId, followingId)
		if err != nil {
			return nil, err
		}
		if isActive == true {
			return nil, fmt.Errorf("%s is followed", followingId)
		}
		err = repository.RelationshipsRepo.UpdateFollow(ctx, followerId, followingId, true)
		if err != nil {
			return nil, err
		}
	} else {
		err := repository.RelationshipsRepo.AddFollow(ctx, followRelationship)
		if err != nil {
			return nil, err
		}
	}
	isMutual, err := repository.RelationshipsRepo.IsActiveFollow(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}
	if isMutual {
		err := repository.RelationshipsRepo.UpdateMutualFollowStatus(ctx, followerId, followingId, true)
		if err != nil {
			return nil, err
		}
	}
	return &map[string]string{
		"status":  "successful",
		"message": followingId + " is followed ",
	}, nil
}

func (s *Relationships) UnFollowUser(ctx context.Context, followerId, followingId string) (any, error) {
	exists, err := repository.RelationshipsRepo.IsFollowExists(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.New("not follow that user")
	}
	isActive, err := repository.RelationshipsRepo.IsActiveFollow(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}
	if !isActive {
		return nil, errors.New("that user is unfollowed")
	}
	err = repository.RelationshipsRepo.UpdateFollow(ctx, followerId, followingId, false)
	if err != nil {
		return nil, err
	}
	err = repository.RelationshipsRepo.UpdateMutualFollowStatus(ctx, followerId, followingId, false)
	if err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"message": followingId + " is unfollowed ",
	}, nil
}
