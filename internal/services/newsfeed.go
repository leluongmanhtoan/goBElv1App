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
	PostNewsfeed(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error)
	GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error)
	ToggleLikePost(ctx context.Context, user_id, post_id string) (any, error)
	GetLikers(ctx context.Context, limit, offset int, post_id string) (any, error)
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

func (s *NewsfeedService) PostNewsfeed(ctx context.Context, user_id string, post *model.NewsfeedPost) (any, error) {
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
		Liked:        false,
	}
	if err := s.repo.PostNews(ctx, newpost); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":  "successful",
		"postId":  newpost.PostId,
		"message": "your post is created",
	}, nil
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
	if err := s.repo.CreateComment(ctx, newcomment); err != nil {
		return nil, err
	}
	return &map[string]string{
		"status":    "successful",
		"commentId": newcomment.CommentId,
		"message":   "your comment is created",
	}, nil
}

func (s *NewsfeedService) GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (any, error) {
	newsfeed, err := s.repo.GetNewsfeed(ctx, limit, offset, user_id)
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

func (s *NewsfeedService) ToggleLikePost(ctx context.Context, user_id, post_id string) (any, error) {
	tx, err := s.repo.GetDBTx(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	postExisted, err := s.repo.IsPostExisted(ctx, post_id)
	if err != nil {
		return nil, err
	}
	if !postExisted {
		return nil, errors.New("postId not found")
	}
	exists, err := s.repo.IsLikeExisted(ctx, post_id, user_id)
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
		err := s.repo.CreateLike(ctx, tx, newlike)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		err = s.repo.IncreaseLikeCount(ctx, tx, post_id)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		isActive, err := s.repo.IsActiveLike(ctx, post_id, user_id)
		if err != nil {
			return nil, err
		}
		if isActive {
			err = s.repo.UpdateLikeTransaction(ctx, tx, user_id, post_id, false)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = s.repo.DecreaseLikeCount(ctx, tx, post_id)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
		} else {
			err = s.repo.UpdateLikeTransaction(ctx, tx, user_id, post_id, true)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			err = s.repo.IncreaseLikeCount(ctx, tx, post_id)
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

func (s *NewsfeedService) GetLikers(ctx context.Context, limit, offset int, post_id string) (any, error) {
	likers, err := s.repo.GetLikers(ctx, limit, offset, post_id)
	if err != nil {
		return nil, err
	}
	return &map[string]any{
		"data":   likers,
		"limit":  limit,
		"offset": offset,
	}, nil
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
