package services

import (
	"context"
	"errors"
	"fmt"
	"program/internal/model"
	newsfeedRepo "program/internal/repositories/newfeed"
	"time"

	"github.com/google/uuid"
)

type INewsfeedService interface {
	CreatePost(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error)
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error)
	ToggleLikePost(ctx context.Context, userId, postId string) error
	GetLikers(ctx context.Context, limit, offset int, userId, post_id string, isGuestUser bool) (any, error)
	PostComment(ctx context.Context, user_id, post_id string, comment *model.CommentPost) (any, error)
	GetComments(ctx context.Context, limit, offset int, post_id string) (any, error)
	PutComment(ctx context.Context, commentPut *model.CommentPut) (any, error)
}
type NewsfeedService struct {
	repo newsfeedRepo.INewsfeedRepo
}

func NewNewsFeedService(repo newsfeedRepo.INewsfeedRepo) INewsfeedService {
	return &NewsfeedService{
		repo: repo,
	}
}

func (s *NewsfeedService) CreatePost(ctx context.Context, userId string, post *model.NewsfeedPost) (any, error) {
	newpost := &model.Post{
		PostId:       uuid.NewString(),
		UserId:       userId,
		Content:      post.Content,
		Privacy:      post.Privacy,
		LikeCount:    0,
		CommentCount: 0,
		ShareCount:   0,
		Deleted:      0,
		CreatedAt:    time.Now(),
	}
	mypost, err := s.repo.CreatePost(ctx, newpost)
	return mypost, err
}

func (s *NewsfeedService) PostComment(ctx context.Context, user_id, post_id string, comment *model.CommentPost) (any, error) {
	newcomment := &model.Comment{
		CommentId:  uuid.NewString(),
		UserId:     user_id,
		PostId:     post_id,
		ParentId:   nil,
		LikeCount:  0,
		ReplyCount: 0,
		Content:    comment.Content,
		Status:     model.ActiveComment,
		CreatedAt:  time.Now(),
	}

	if comment.Parent != "" {
		//Check parent is existed
		newcomment.ParentId = &comment.Parent
	}
	mycomment, err := s.repo.CreateComment(ctx, newcomment)
	return mycomment, err
}

func (s *NewsfeedService) GetNewsfeed(ctx context.Context, limit, offset int, userId string) (any, error) {
	cacheKey := fmt.Sprintf("newsfeed:user:%s", userId)

	newsfeed, err := s.repo.GetNewsfeed(ctx, limit, offset, userId)
	if err != nil {
		return nil, err
	}
	s.repo.SaveNewsfeedCache(ctx, cacheKey, newsfeed)

	return newsfeed, nil
}

func (s *NewsfeedService) ToggleLikePost(ctx context.Context, userId, postId string) error {
	tx, err := s.repo.GetDBTx(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	postExisted, err := s.repo.IsPostExisted(ctx, postId)
	if err != nil {
		return err
	}
	if !postExisted {
		return errors.New("this post was not found")
	}
	likeExisted, err := s.repo.IsLikeExisted(ctx, postId, userId)
	if err != nil {
		return err
	}
	if !likeExisted {
		newlike := &model.Like{
			LikeId:    uuid.NewString(),
			PostId:    postId,
			UserId:    userId,
			Type:      model.LikePost,
			IsActive:  true,
			CreatedAt: time.Now(),
		}
		err := s.repo.CreateLike(ctx, tx, newlike)
		if err != nil {
			return err
		}
		err = s.repo.IncreaseLikeCount(ctx, tx, postId)
		if err != nil {
			return err
		}
	} else {
		isActive, err := s.repo.IsActiveLike(ctx, postId, userId)
		if err != nil {
			return err
		}
		if isActive {
			err = s.repo.UpdateLikeTransaction(ctx, tx, userId, postId, false)
			if err != nil {
				return err
			}
			err = s.repo.DecreaseLikeCount(ctx, tx, postId)
			if err != nil {
				return err
			}
		} else {
			err = s.repo.UpdateLikeTransaction(ctx, tx, userId, postId, true)
			if err != nil {
				return err
			}
			err = s.repo.IncreaseLikeCount(ctx, tx, postId)
			if err != nil {
				return err
			}
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (s *NewsfeedService) GetLikers(ctx context.Context, limit, offset int, userId, post_id string, isGuestUser bool) (any, error) {
	if isGuestUser {
		err := s.repo.CheckPublicPrivacyPermission(ctx, post_id)
		if err != nil {
			return nil, err
		}
	} else {
		err := s.repo.CheckFriendPrivacyPermission(ctx, userId, post_id)
		if err != nil {
			return nil, err
		}
	}
	likers, err := s.repo.GetLikers(ctx, limit, offset, post_id)
	if err != nil {
		return nil, err
	}
	return likers, nil
}

func (s *NewsfeedService) GetComments(ctx context.Context, limit, offset int, post_id string) (any, error) {
	comments, err := s.repo.GetComments(ctx, limit, offset, post_id)
	if err != nil {
		return nil, err
	}
	fmt.Println(comments)
	return &map[string]any{
		"post_id": post_id,
		"data":    comments,
		"limit":   limit,
		"offset":  offset,
	}, nil
}

func (s *NewsfeedService) PutComment(ctx context.Context, commentPut *model.CommentPut) (any, error) {
	err := s.repo.PutComment(ctx, commentPut.CommentId, commentPut.Content)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"comment_id": commentPut.CommentId,
		"message":    "modify comment successfully",
	}, nil
}
