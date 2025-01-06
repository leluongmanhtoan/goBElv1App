package model

import (
	"time"

	"github.com/uptrace/bun"
)

type UserProfile struct {
	bun.BaseModel `bun:"profiles"`
	ProfileId     string    `json:"id" bun:"profileId,type:varchar(36),pk,notnull"`
	UserId        string    `json:"userId" bun:"userId,type:varchar(36),notnull"`
	FirstName     string    `json:"firstname" bun:"firstname,type:varchar(255),notnull"`
	LastName      string    `json:"lastname" bun:"lastname,type:varchar(255),notnull"`
	Gender        int       `json:"gender" bun:"gender,type:tinyint,notnull"`
	Avatar        string    `json:"avatarUrl,omitempty" bun:"avatarUrl,type:varchar(255)"`
	Address       string    `json:"address,omitempty" bun:"address,type:varchar(255)"`
	Email         string    `json:"email" bun:"email,type:varchar(150),notnull"`
	PhoneNumber   string    `json:"phone,omitempty" bun:"phoneNumber,type:varchar(20)"`
	CreatedAt     time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
	UserAuth      *User     `json:"accounts,omitempty" bun:"rel:belongs-to,join:userId=id"`
}

type UserProfilePost struct {
	FirstName   string `json:"firstname" validate:"required"`
	LastName    string `json:"lastname" validate:"required"`
	Gender      int    `json:"gender" validate:"required,oneof=0 1 2"`
	Avatar      string `json:"avatarUrl"`
	Address     string `json:"address"`
	Email       string `json:"email" validate:"required"`
	PhoneNumber string `json:"phone"`
}

type UserProfilePut struct {
	FirstName   string `json:"firstname"`
	LastName    string `json:"lastname"`
	Gender      *int   `json:"gender"`
	Avatar      string `json:"avatarUrl"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone"`
}
