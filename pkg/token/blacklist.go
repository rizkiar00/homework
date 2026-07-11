package token

import (
	"context"
	"time"

	"github.com/rizkiar00/homework/internal/model"
)

const blacklistKeyPrefix = "auth:blacklist:"

type Blacklist struct {
	redis model.Redis
}

func NewBlacklist(redis model.Redis) *Blacklist {
	return &Blacklist{redis: redis}
}

func (b *Blacklist) Revoke(ctx context.Context, claims Claims) error {
	if b.redis == nil || claims.ID == "" || claims.ExpiresAt == nil {
		return nil
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl <= 0 {
		return nil
	}

	return b.redis.Set(ctx, blacklistKeyPrefix+claims.ID, "revoked", ttl).Err()
}

func (b *Blacklist) IsRevoked(ctx context.Context, claims Claims) bool {
	if b.redis == nil || claims.ID == "" {
		return false
	}

	err := b.redis.Get(ctx, blacklistKeyPrefix+claims.ID).Err()
	return err == nil
}
