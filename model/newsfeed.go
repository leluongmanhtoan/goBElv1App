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
}
