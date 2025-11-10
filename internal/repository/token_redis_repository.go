package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRedisRepository interface {
	SetToken(ctx context.Context, key, value string, ttl time.Duration) error
	GetToken(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type tokenRedisRepository struct {
	client *redis.Client
}

func NewTokenRedisRepository(client *redis.Client) TokenRedisRepository {
	return &tokenRedisRepository{client}
}

func (t *tokenRedisRepository) SetToken(ctx context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("invalid ttl value: %v", ttl)
	}
	if err := t.client.Set(ctx, key, value, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set redis token: %w", err)
	}
	return nil
}

func (t *tokenRedisRepository) GetToken(ctx context.Context, key string) (string, error) {
	value, err := t.client.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", errors.New("access token not found in redis")
	}

	if err != nil {
		return "", err
	}

	return value, nil
}

func (t *tokenRedisRepository) Exists(ctx context.Context, key string) (bool, error) {
	count, err := t.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check redis key: %w", err)
	}
	return count > 0, nil
}
