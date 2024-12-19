package model

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"accounts"`
	UserUuid      string    `json:"id" bun:"id,type:varchar(36),pk,notnull"`
	Username      string    `json:"username" bun:"username,type:varchar(50),notnull"`
	Salt          string    `json:"salt" bun:"salt,type:varchar(64),notnull"`
	Hash          string    `json:"hashPassword" bun:"hashPassword,type:varchar(255),notnull"`
	CreatedAt     time.Time `json:"createdAt" bun:"createdAt,type:timestamp,notnull,nullzero"`
	UpdatedAt     time.Time `json:"updatedAt" bun:"updatedAt,type:timestamp,nullzero"`
	Deleted       int       `json:"deleted" bun:"deleted,type:tinyint,notnull"`
}

type (
	Login struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	LoginResponse struct {
		UserID       string `json:"userID"`
		Username     string `json:"username"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)

type (
	Register struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	RegisterResponse struct {
		Message      string `json:"message"`
		UserUuid     string `json:"userUUID"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}
)

type RefreshToken struct {
	UserId         string `json:"userId"`
	NewAccessToken string `json:"accessToken"`
}

type (
	Logout struct {
		RefreshToken string `json:"refreshToken"`
	}
	LogoutResponse struct {
		Message string `json:"message"`
	}
)
