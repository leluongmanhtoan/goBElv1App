package model

import (
	"time"

	"github.com/uptrace/bun"
)

type Follows struct {
	bun.BaseModel `bun:"follows"`
	FollowId      int64     `json:"followId" bun:"id,pk,autoincrement"`
	FollowerId    string    `json:"followerId" bun:"followerId,type:varchar(36),notnull"`
	FollowingId   string    `json:"followingId" bun:"followingId,type:varchar(36),notnull"`
	IsActive      bool      `json:"isActive" bun:"isActive,default:1"`
	IsMutual      bool      `json:"isMutual" bun:"isMutual,default:0"`
	CreatedAt     time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
}

type FollowerInfo struct {
	ProfileId string `json:"profileId" bun:"profileId"`
	FirstName string `json:"firstname" bun:"firstname"`
	Lastname  string `json:"lastname" bun:"lastname"`
	Avatar    string `json:"avatar" bun:"avatarUrl"`
}
