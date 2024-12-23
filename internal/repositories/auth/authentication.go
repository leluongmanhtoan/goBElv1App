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
	_, err := r.rd.GetDB().Set(ctx, "refresh:"+tokenId, "userId:"+userId, ttl).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthenticationRepo) AddAccessToBlacklist(ctx context.Context, accessToken string, ttl time.Duration) error {
	_, err := r.rd.GetDB().Set(ctx, "blacklist:accessToken:"+accessToken, "revoked", ttl).Result()
	if err != nil {
		return err
	}
	return nil
}

func (r *AuthenticationRepo) IsExisted(ctx context.Context, key string) (bool, error) {
	res, err := r.rd.GetDB().Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return res > 0, err
}

func (r *AuthenticationRepo) DelRefreshToken(ctx context.Context, key string) error {
	_, err := r.rd.GetDB().Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
