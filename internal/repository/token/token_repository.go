package tokenrepo

import (
	"context"
	"fmt"
	"time"

	"auth_service/internal/config"

	"github.com/redis/go-redis/v9"
)

type Token_Repository interface {
	StoreRefreshToken(ctx context.Context, userID int64, token string) error
	GetRefreshToken(ctx context.Context, userID int64) (string, error)
	DeleteRefreshToken(ctx context.Context, userID int64) error
	StoreBlacklistedToken(ctx context.Context, token string, ttl time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type TokenRepository struct {
	redisClient *redis.Client
}

func NewTokenRepository(redisClient *redis.Client) *TokenRepository {
	return &TokenRepository{redisClient: redisClient}
}

func (r *TokenRepository) StoreRefreshToken(ctx context.Context, userID int64, token string) error {
	key := fmt.Sprintf("refresh_token:%d", userID)
	ttl := parseDuration(config.App.JWT.RefreshTTL)

	err := r.redisClient.Set(ctx, key, token, ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

func (r *TokenRepository) GetRefreshToken(ctx context.Context, userID int64) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userID)

	token, err := r.redisClient.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return "", fmt.Errorf("refresh token not found")
		}
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}

	return token, nil
}

func (r *TokenRepository) DeleteRefreshToken(ctx context.Context, userID int64) error {
	key := fmt.Sprintf("refresh_token:%d", userID)

	err := r.redisClient.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (r *TokenRepository) StoreBlacklistedToken(ctx context.Context, token string, ttl time.Duration) error {
	key := fmt.Sprintf("blacklisted_token:%s", token)

	err := r.redisClient.Set(ctx, key, "1", ttl).Err()
	if err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func (r *TokenRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	key := fmt.Sprintf("blacklisted_token:%s", token)

	exists, err := r.redisClient.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check blacklisted token: %w", err)
	}

	return exists == 1, nil
}

func parseDuration(durationStr string) time.Duration {
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return 720 * time.Hour
	}
	return dur
}
