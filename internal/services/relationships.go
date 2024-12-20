package services

import (
	"context"
	"errors"
	"program/internal/model"
	relationshipsRepo "program/internal/repositories/relationships"

	"time"
)

type IRelationshipsService interface {
	GetFollowers(ctx context.Context, limit, offset int, userId string) (any, error)
	GetFollowing(ctx context.Context, limit, offset int, userId string) (any, error)
	GetFollowRelationshipCount(ctx context.Context, userId string) (any, error)
	ToggleFollow(ctx context.Context, followerId, followingId string) (any, error)
}

type RelationshipsService struct {
	repo relationshipsRepo.IRelationshipsRepo
}

func NewRelationshipsService(repo relationshipsRepo.IRelationshipsRepo) IRelationshipsService {
	return &RelationshipsService{
		repo: repo,
	}
}

func (s *RelationshipsService) ToggleFollow(ctx context.Context, followerId, followingId string) (any, error) {
	var isActive = false
	var message = ""
	var isMutual = false
	if followerId == followingId {
		return nil, errors.New("can not follow yourself")
	}

	tx, err := s.repo.GetDBTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	exists, err := s.repo.IsFollowExists(ctx, followerId, followingId)
	if err != nil {
		return nil, err
	}

	if exists {
		isActive, err = s.repo.IsActiveFollow(ctx, followerId, followingId)
		if err != nil {
			return nil, err
		}
		if !isActive {
			err = s.repo.UpdateFollowTransaction(ctx, tx, followerId, followingId, true)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			message = "refollow " + followingId + " successful"
		} else {
			err = s.repo.UpdateFollowTransaction(ctx, tx, followerId, followingId, false)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			err = s.repo.UpdateMutualFollowStatusTransaction(ctx, tx, followerId, followingId, false)
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
		err := s.repo.AddFollowTransaction(ctx, tx, followRelationship)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		message = "follow " + followingId + " successful"
	}
	if exists && !isActive || !exists {
		isMutual, err = s.repo.IsActiveFollow(ctx, followingId, followerId)
		if err != nil {
			return nil, err
		}
		if isMutual {
			err := s.repo.UpdateMutualFollowStatusTransaction(ctx, tx, followerId, followingId, true)
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

func (s *RelationshipsService) GetFollowers(ctx context.Context, limit, offset int, userId string) (any, error) {
	total, followers, err := s.repo.GetFollowList(ctx, limit, offset, userId, true)
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

func (s *RelationshipsService) GetFollowing(ctx context.Context, limit, offset int, userId string) (any, error) {
	total, following, err := s.repo.GetFollowList(ctx, limit, offset, userId, false)
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

func (s *RelationshipsService) GetFollowRelationshipCount(ctx context.Context, userId string) (any, error) {
	totalFollower, totalFollowing, err := s.repo.NumOfFollowRelationship(ctx, userId)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"num_of_followers": totalFollower,
		"num_of_following": totalFollowing,
	}, nil
}
