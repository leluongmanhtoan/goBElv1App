package service

import (
	"context"
	"errors"
	"program/model"
	"program/repository"
	"time"
)

type IRelationships interface {
	GetFollowers(ctx context.Context, limit, offset int, userId string) (any, error)
	GetFollowing(ctx context.Context, limit, offset int, userId string) (any, error)
	GetFollowRelationshipCount(ctx context.Context, userId string) (any, error)
	ToggleFollow(ctx context.Context, followerId, followingId string) (any, error)
}

type Relationships struct {
}

func NewRelationships() IRelationships {
	return &Relationships{}
}

func (s *Relationships) ToggleFollow(ctx context.Context, followerId, followingId string) (any, error) {

	var isActive = false
	var message = ""
	var isMutual = false
	if followerId == followingId {
		return nil, errors.New("can not follow yourself")
	}
	tx, err := repository.SqlClientConnection.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	exists, err := repository.RelationshipsRepo.IsFollowExists(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}

	if exists {
		isActive, err = repository.RelationshipsRepo.IsActiveFollow(ctx, followerId, followingId)
		if err != nil {
			return nil, err
		}
		if !isActive {
			err = repository.RelationshipsRepo.UpdateFollowTransaction(ctx, &tx, followerId, followingId, true)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			message = "refollow " + followingId + " successful"
		} else {
			err = repository.RelationshipsRepo.UpdateFollowTransaction(ctx, &tx, followerId, followingId, false)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			err = repository.RelationshipsRepo.UpdateMutualFollowStatusTransaction(ctx, &tx, followerId, followingId, false)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			message = "unfollow " + followingId + " successful"
		}
	} else {
		followRelationship := &model.Follows{
			FollowerId:  followerId,
			FollowingId: followingId,
			IsActive:    true,
			IsMutual:    false,
			CreatedAt:   time.Now(),
		}
		err := repository.RelationshipsRepo.AddFollowTransaction(ctx, &tx, followRelationship)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		message = "follow " + followingId + " successful"
	}
	if exists && !isActive || !exists {
		isMutual, err = repository.RelationshipsRepo.IsActiveFollow(ctx, followingId, followerId)
		if err != nil {
			return nil, err
		}
		if isMutual {
			err := repository.RelationshipsRepo.UpdateMutualFollowStatusTransaction(ctx, &tx, followerId, followingId, true)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"status":   "successful",
		"message":  message,
		"isMutual": isMutual,
	}, nil
}

func (s *Relationships) GetFollowers(ctx context.Context, limit, offset int, userId string) (any, error) {
	total, followers, err := repository.RelationshipsRepo.GetFollowList(ctx, limit, offset, userId, true)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"data":   followers,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}, nil
}

func (s *Relationships) GetFollowing(ctx context.Context, limit, offset int, userId string) (any, error) {
	total, following, err := repository.RelationshipsRepo.GetFollowList(ctx, limit, offset, userId, false)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"data":   following,
		"limit":  limit,
		"offset": offset,
		"total":  total,
	}, nil
}

func (s *Relationships) GetFollowRelationshipCount(ctx context.Context, userId string) (any, error) {
	totalFollower, totalFollowing, err := repository.RelationshipsRepo.NumOfFollowRelationship(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"num_of_followers": totalFollower,
		"num_of_following": totalFollowing,
	}, nil
}
