package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	tokenBlacklistKey = "token_blacklist"
)

type TokenRedisRepository struct {
	redis *redis.Client
}

func NewTokenRedisRepository(redis *redis.Client) *TokenRedisRepository {
	return &TokenRedisRepository{redis: redis}
}

func (r *TokenRedisRepository) Set(ctx context.Context, token string, expiration time.Duration) error {
	return r.redis.Set(ctx, fmt.Sprintf("%s:%s", tokenBlacklistKey, token), 1, expiration).Err()
}

func (r *TokenRedisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	result, err := r.redis.Get(ctx, fmt.Sprintf("%s:%s", tokenBlacklistKey, token)).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return result == "1", nil
}
