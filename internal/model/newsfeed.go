package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Privacy string

const (
	Public  Privacy = "public"
	Private Privacy = "private"
	Friends Privacy = "friends"
)

type LikeType string

const (
	LikePost    LikeType = "post"
	LikeComment LikeType = "comment"
)

type CommentStatus string

const (
	ActiveComment  CommentStatus = "active"
	HiddenComment  CommentStatus = "hidden"
	DeletedComment CommentStatus = "deleted"
)

type Post struct {
	bun.BaseModel `bun:"posts"`
	PostId        string    `json:"id" bun:"postId,type:varchar(36),pk,notnull"`
	UserId        string    `json:"userId" bun:"userId,type:varchar(36),notnull"`
	Content       string    `json:"content" bun:"content,type:text,notnull"`
	Privacy       Privacy   `json:"privacy" bun:"privacy,type:enum"`
	LikeCount     int64     `json:"likeCount" bun:"likeCount,type:int"`
	CommentCount  int64     `json:"commentCount" bun:"commentCount,type:int"`
	ShareCount    int64     `json:"shareCount" bun:"shareCount,type:int"`
	Deleted       int       `json:"deleted" bun:"deleted,type:tinyint,notnull"`
	CreatedAt     time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
}

type NewsfeedPost struct {
	Content string  `json:"content" validate:"required"`
	Privacy Privacy `json:"privacy" validate:"required"`
}

type NewsFeed struct {
	PostId       string    `json:"postId" bun:"postId"`
	AvatarUrl    string    `json:"avatarUrl" bun:"avatarUrl"`
	FirstName    string    `json:"firstname" bun:"firstname"`
	LastName     string    `json:"lastname" bun:"lastname"`
	Content      string    `json:"content" bun:"content"`
	Privacy      Privacy   `json:"privacy" bun:"privacy"`
	LikeCount    int       `json:"likeCount" bun:"likeCount"`
	CommentCount int       `json:"commentCount" bun:"commentCount"`
	ShareCount   int       `json:"shareCount" bun:"shareCount"`
	CreatedAt    time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt    time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
	Liked        bool      `json:"liked" bun:"liked"`
}

type Like struct {
	bun.BaseModel `bun:"likes"`
	LikeId        string    `json:"likeId" bun:"likeId,type:varchar(36),pk,notnull"`
	PostId        string    `json:"postId" bun:"postId,type:varchar(36),notnull"`
	UserId        string    `json:"userId" bun:"userId,type:varchar(36),notnull"`
	Type          LikeType  `json:"type" bun:"type"`
	IsActive      bool      `json:"isActive" bun:"isActive,default:1"`
	CreatedAt     time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
}

type LikerInfo struct {
	ProfileId string `json:"profileId" bun:"profileId"`
	FirstName string `json:"firstname" bun:"firstname"`
	Lastname  string `json:"lastname" bun:"lastname"`
	Avatar    string `json:"avatar" bun:"avatarUrl"`
}

type Comment struct {
	bun.BaseModel `bun:"comments"`
	CommentId     string        `json:"commentId" bun:"commentId,type:varchar(36),pk,notnull"`
	PostId        string        `json:"postId" bun:"postId,type:varchar(36),notnull"`
	UserId        string        `json:"userId" bun:"userId,type:varchar(36),notnull"`
	ParentId      *string       `json:"parentId" bun:"parentId,type:varchar(36)"`
	LikeCount     int           `json:"likeCount" bun:"likeCount"`
	ReplyCount    int           `json:"repliesCount" bun:"repliesCount"`
	Content       string        `json:"content" bun:"content,type:text,notnull"`
	Status        CommentStatus `json:"status" bun:"status"`
	CreatedAt     time.Time     `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time     `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
}

type CommentPost struct {
	Content string `json:"content"`
	Parent  string `json:"parent"`
}

type CommentInfo struct {
	ProfileId string    `json:"profileId" bun:"profileId"`
	FirstName string    `json:"firstname" bun:"firstname"`
	Lastname  string    `json:"lastname" bun:"lastname"`
	Avatar    string    `json:"avatar" bun:"avatarUrl"`
	Content   string    `json:"content" bun:"content"`
	CreatedAt time.Time `json:"createdAt" bun:"createdAt"`
}

type CommentPut struct {
	CommentId string `json:"commentId"`
	Content   string `json:"content"`
}
