package authenticationRepo

import (
	"context"
	"time"
)

type IAuthenticationRepo interface {
	AddValidRefreshToken(ctx context.Context, userId, tokenId string, ttl time.Duration) error
}
