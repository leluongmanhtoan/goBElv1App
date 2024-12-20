package authenticationRepo

import (
	"context"
	"program/internal/database"
	"time"
)

type AuthenticationRepo struct {
	rd database.IRedisConnection
}

func NewAuthenticationRepo(rd database.IRedisConnection) IAuthenticationRepo {
	return &AuthenticationRepo{
		rd: rd,
	}
}

func (r *AuthenticationRepo) AddValidRefreshToken(ctx context.Context, userId, tokenId string, ttl time.Duration) error {
	_, err := r.rd.GetDB().Set(ctx, "userId:"+userId, tokenId, ttl).Result()
	if err != nil {
		return err
	}
	return nil
}
