package service

import (
	"context"
	"errors"
	"fmt"
	"program/model"
	"program/repository"
	"time"

	"github.com/google/uuid"
)

type INewsfeed interface {
	PostNewsfeed(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error)
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error)
	ToggleLikePost(ctx context.Context, user_id, post_id string) (any, error)
	GetLikers(ctx context.Context, limit, offset int, post_id string) (any, error)
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
	fmt.Println(newsfeed)
	return &map[string]any{
		"userId": user_id,
		"data":   newsfeed,
		"limit":  limit,
		"offset": offset,
	}, nil
}

func (s *Newsfeed) ToggleLikePost(ctx context.Context, user_id, post_id string) (any, error) {
	tx, err := repository.SqlClientConnection.GetDB().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	postExisted, err := repository.NewsfeedRepo.IsPostExisted(ctx, post_id)
	if err != nil {
		return nil, err
	}
	if !postExisted {
		return nil, errors.New("postId not found")
	}
	exists, err := repository.NewsfeedRepo.IsLikeExisted(ctx, post_id, user_id)
	if err != nil {
		return nil, err
	}
	if !exists {
		newlike := &model.Like{
			LikeId:    uuid.NewString(),
			PostId:    post_id,
			UserId:    user_id,
			Type:      model.LikePost,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		err := repository.NewsfeedRepo.CreateLike(ctx, &tx, newlike)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = repository.NewsfeedRepo.IncreaseLikeCount(ctx, &tx, post_id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		isActive, err := repository.NewsfeedRepo.IsActiveLike(ctx, post_id, user_id)
		if err != nil {
			return nil, err
		}
		if isActive {
			err = repository.NewsfeedRepo.UpdateLikeTransaction(ctx, &tx, user_id, post_id, false)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = repository.NewsfeedRepo.DecreaseLikeCount(ctx, &tx, post_id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			err = repository.NewsfeedRepo.UpdateLikeTransaction(ctx, &tx, user_id, post_id, true)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = repository.NewsfeedRepo.IncreaseLikeCount(ctx, &tx, post_id)
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
		"status": "successful",
	}, nil

}

func (s *Newsfeed) GetLikers(ctx context.Context, limit, offset int, post_id string) (any, error) {
	likers, err := repository.NewsfeedRepo.GetLikers(ctx, limit, offset, post_id)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"data":   likers,
		"limit":  limit,
		"offset": offset,
	}, nil
}
