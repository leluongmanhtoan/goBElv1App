package newsfeedRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"program/internal/database"
	"program/internal/model"
	"time"

	"github.com/uptrace/bun"
)

type NewsfeedRepo struct {
	db database.ISqlConnection
}

func NewNewsfeedRepo(db database.ISqlConnection) INewsfeedRepo {
	return &NewsfeedRepo{
		db: db,
	}
}

func (r *NewsfeedRepo) GetDBTx(ctx context.Context) (*bun.Tx, error) {
	tx, err := r.db.GetDB().BeginTx(ctx, nil)
	return &tx, err
}

func (r *NewsfeedRepo) CreatePost(ctx context.Context, post *model.Post) (*model.NewsFeed, error) {
	_, err := r.db.GetDB().NewInsert().
		Model(post).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	myPost := new(model.NewsFeed)
	myQuery := r.db.GetDB().NewSelect().
		Column(
			"p.postId",
			"pf.avatarUrl",
			"pf.firstname",
			"pf.lastname",
			"p.content",
			"p.privacy",
			"p.likeCount",
			"p.commentCount",
			"p.shareCount",
			"p.createdAt",
			"p.updatedAt").
		TableExpr("posts as p").
		Join("JOIN profiles pf ON pf.userId = p.userId").
		Where("p.postId = ?", post.PostId)
	err = myQuery.Scan(ctx, myPost)
	if err != nil {
		if err == sql.ErrNoRows {
			return myPost, nil
		}
		return nil, err
	}
	return myPost, nil
}

func (r *NewsfeedRepo) GetNewsfeed(ctx context.Context, limit, offset int, user_id string) (*[]model.NewsFeed, error) {
	newsfeed := new([]model.NewsFeed)
	othersQuery := r.db.GetDB().NewSelect().
		Column(
			"p.postId",
			"pf.avatarUrl",
			"pf.firstname",
			"pf.lastname",
			"p.content",
			"p.privacy",
			"p.likeCount",
			"p.commentCount",
			"p.shareCount",
			"p.createdAt",
			"p.updatedAt").
		ColumnExpr("IF(l.postId IS NOT NULL AND l.isActive = 1,TRUE,FALSE) AS liked").
		TableExpr("follows as f").
		Join("JOIN profiles pf ON pf.userId = f.followingId").
		Join("JOIN posts p ON p.userId = f.followingId").
		Join("LEFT JOIN likes l ON l.postId = p.postId AND l.userId = ?", user_id).
		Where("f.followerId = ? AND f.isActive = 1 AND p.deleted = 0 AND p.createdAt >= NOW() - INTERVAL 7 DAY AND (f.isMutual = 1 OR p.privacy = 'public')", user_id)

	myQuery := r.db.GetDB().NewSelect().
		Column(
			"p.postId",
			"pf.avatarUrl",
			"pf.firstname",
			"pf.lastname",
			"p.content",
			"p.privacy",
			"p.likeCount",
			"p.commentCount",
			"p.shareCount",
			"p.createdAt",
			"p.updatedAt").
		ColumnExpr("IF(l.postId IS NOT NULL AND l.isActive = 1,TRUE,FALSE) AS liked").
		TableExpr("posts as p").
		Join("JOIN profiles pf ON pf.userId = p.userId").
		Join("LEFT JOIN likes l ON l.postId = p.postId AND l.userId = ?", user_id).
		Where("p.userId = ? AND p.deleted = 0 AND p.createdAt >= NOW() - INTERVAL 7 DAY", user_id)

	unionQuery := r.db.GetDB().NewSelect().With("others", othersQuery).With("mine", myQuery).TableExpr("(SELECT * FROM others UNION ALL SELECT * FROM mine) AS newsfeed").
		OrderExpr("createdAt DESC")
	if limit > 0 {
		unionQuery.Limit(limit).Offset(offset)
	}
	err := unionQuery.Scan(ctx, newsfeed)
	if err != nil {
		if err == sql.ErrNoRows {
			return newsfeed, nil
		}
		return nil, err
	}
	return newsfeed, nil
}

func (r *NewsfeedRepo) CreateLike(ctx context.Context, tx *bun.Tx, like *model.Like) error {
	_, err := tx.NewInsert().
		Model(like).
		Exec(ctx)
	return err
}

func (r *NewsfeedRepo) IncreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error {
	_, err := tx.NewUpdate().Model((*model.Post)(nil)).
		Set("likeCount = likeCount + 1").
		Where("postId = ?", postId).
		Exec(ctx)
	return err
}

func (r *NewsfeedRepo) SetOwnerLikedStatus(ctx context.Context, tx *bun.Tx, postId string, status bool) error {
	_, err := tx.NewUpdate().Model((*model.Post)(nil)).Set("liked = ?", status).
		Where("postId = ?", postId).
		Exec(ctx)
	return err
}

func (r *NewsfeedRepo) DecreaseLikeCount(ctx context.Context, tx *bun.Tx, postId string) error {
	_, err := tx.NewUpdate().Model((*model.Post)(nil)).
		Set("likeCount = likeCount - 1").
		Where("postId = ?", postId).
		Exec(ctx)
	return err
}

func (r *NewsfeedRepo) IsLikeExisted(ctx context.Context, postId, userId string) (bool, error) {
	exists, err := r.db.GetDB().NewSelect().
		Model((*model.Like)(nil)).
		ColumnExpr("1").
		Where("postId = ?", postId).
		Where("userId = ?", userId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking like post: %w", err)
	}
	return exists, nil
}

func (r *NewsfeedRepo) IsPostExisted(ctx context.Context, postId string) (bool, error) {
	exists, err := r.db.GetDB().NewSelect().
		Model((*model.Post)(nil)).
		ColumnExpr("1").
		Where("postId = ?", postId).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking exist post: %w", err)
	}
	return exists, nil
}

func (r *NewsfeedRepo) IsActiveLike(ctx context.Context, post_id, user_id string) (bool, error) {
	isActive, err := r.db.GetDB().NewSelect().
		Model((*model.Like)(nil)).
		ColumnExpr("1").
		Where("userId = ? AND postId = ? AND isActive = ?", user_id, post_id, true).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking like status: %w", err)
	}
	return isActive, nil
}

func (r *NewsfeedRepo) IsOwnPost(ctx context.Context, post_id, user_id string) (bool, error) {
	isOwnPost, err := r.db.GetDB().NewSelect().
		Model((*model.Post)(nil)).
		ColumnExpr("1").
		Where("userId = ? AND postId = ?", user_id, post_id).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("error checking liked: %w", err)
	}
	return isOwnPost, nil
}

func (r *NewsfeedRepo) UpdateLikeTransaction(ctx context.Context, tx *bun.Tx, user_id, post_id string, status bool) error {
	_, err := tx.NewUpdate().
		Model((*model.Like)(nil)).
		Set("isActive = ?", status).
		Set("updatedAt = ?", time.Now()).
		Where("userId = ? AND postId = ?", user_id, post_id).
		Exec(ctx)
	return err
}

func (r *NewsfeedRepo) GetLikers(ctx context.Context, limit, offset int, post_id string) (*[]model.LikerInfo, error) {
	likers := new([]model.LikerInfo)
	query := r.db.GetDB().NewSelect().
		Column("p.profileId", "p.firstname", "p.lastname", "p.avatarUrl").
		TableExpr("likes as l").
		Join("JOIN profiles p ON p.userId=l.userId").
		Where("postId = ? AND isActive = 1", post_id).
		Order("p.lastname ASC")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}
	err := query.Scan(ctx, likers)
	if err != nil {
		if err == sql.ErrNoRows {
			return likers, nil
		}
		return nil, err
	}
	return likers, nil
}

func (r *NewsfeedRepo) CheckPublicPrivacyPermission(ctx context.Context, postId string) error {
	post := new(model.Post)
	query := r.db.GetDB().NewSelect().
		Model(post).Where("postId = ?", postId)
	if err := query.Scan(ctx, post); err != nil {
		return err
	}

	if post.Privacy != "public" {
		return errors.New("you don't have permission to get likers in this post")
	}
	return nil
}

func (r *NewsfeedRepo) CheckFriendPrivacyPermission(ctx context.Context, userId string, postId string) error {
	post := new(model.Post)
	query := r.db.GetDB().NewSelect().Model(post).Where("postId = ?", postId)
	if err := query.Scan(ctx, post); err != nil {
		return err
	}
	if post.Privacy == "public" {
		return nil
	}
	friendCheck, err := r.db.GetDB().NewSelect().
		TableExpr("follows as f").
		Where("f.followerId = ? AND f.followingId = ? AND isMutual = 1", post.UserId, userId).Exists(ctx)
	if err != nil {
		return err
	}
	if userId != post.UserId && friendCheck == false {
		return errors.New("you don't have permission to get likers in this post")
	}
	return nil
}
func (r *NewsfeedRepo) CreateComment(ctx context.Context, commentPost *model.Comment) (*model.CommentInfo, error) {
	_, err := r.db.GetDB().NewInsert().
		Model(commentPost).
		Exec(ctx)
	if err != nil {
		return nil, err
	}
	myComment := new(model.CommentInfo)
	myQuery := r.db.GetDB().NewSelect().
		Column(
			"c.commentId",
			"pf.profileId",
			"pf.firstname",
			"pf.lastname",
			"pf.avatarUrl",
			"c.content",
			"c.createdAt").
		TableExpr("comments as c").
		Join("JOIN profiles pf ON pf.userId = c.userId").
		Where("c.commentId = ?", commentPost.CommentId)
	err = myQuery.Scan(ctx, myComment)
	if err != nil {
		if err == sql.ErrNoRows {
			return myComment, nil
		}
		return nil, err
	}
	return myComment, nil
}

func (r *NewsfeedRepo) GetComments(ctx context.Context, limit, offset int, postId string) (*[]model.CommentInfo, error) {
	comments := new([]model.CommentInfo)
	query := r.db.GetDB().NewSelect().
		Column("c.commentId", "p.profileId", "p.firstname", "p.lastname", "p.avatarUrl", "c.createdAt", "c.content").
		TableExpr("comments as c").
		Join("JOIN profiles p ON p.userId = c.userId").
		Where("c.postId = ?", postId).
		OrderExpr("createdAt DESC")
	if limit > 0 {
		query.Limit(limit).Offset(offset)
	}

	err := query.Scan(ctx, comments)
	if err != nil {
		if err == sql.ErrNoRows {
			return comments, nil
		}
		return nil, err
	}
	return comments, nil
}

func (r *NewsfeedRepo) PutComment(ctx context.Context, commentId string, content string) error {
	commentPut := model.Comment{
		Content:   content,
		UpdatedAt: time.Now(),
	}
	resp, err := r.db.GetDB().NewUpdate().Model(commentPut).WherePK("commentId = ?", commentId).Exec(ctx)
	if err != nil {
		return err
	} else if affected, _ := resp.RowsAffected(); affected < 1 {
		return errors.New("update comment failed")
	}
	return nil
}
