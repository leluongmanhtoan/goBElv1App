package authenticationRepo

import (
	"context"
	"time"
)

type IAuthenticationRepo interface {
	AddValidRefreshToken(ctx context.Context, userId, tokenId string, ttl time.Duration) error
	AddAccessToBlacklist(ctx context.Context, accessToken string, ttl time.Duration) error
	IsExisted(ctx context.Context, key string) (bool, error)
	DelRefreshToken(ctx context.Context, key string) error
}
