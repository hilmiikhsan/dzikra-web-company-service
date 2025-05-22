package ports

import (
	"context"
	"fmt"
	"time"

	"github.com/Digitalkeun-Creative/be-dzikra-web-company-service/internal/infrastructure/redis/ports"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

var _ ports.RedisRepository = &redisRepository{}

type redisRepository struct {
	db *redis.Client
}

func NewRedisRepository(db *redis.Client) *redisRepository {
	return &redisRepository{
		db: db,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	// Check if the key already exists
	exists, err := r.db.Exists(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to check key existence")
		return fmt.Errorf("failed to check key existence: %w", err)
	}

	// If the key does not exist, set the value
	if exists == 0 {
		err := r.db.Set(ctx, key, value, expiration).Err()
		if err != nil {
			log.Error().Err(err).Msg("failed to set key")
			return fmt.Errorf("failed to set key: %w", err)
		}
	}

	return nil
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	data, err := r.db.Get(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to get key")
		return "", err
	}

	return data, nil
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	_, err := r.db.Del(ctx, key).Result()
	if err != nil {
		log.Error().Err(err).Msg("failed to delete key")
		return err
	}

	return nil
}

func (r *redisRepository) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.db.TTL(ctx, key).Result()
}

func (r *redisRepository) Incr(ctx context.Context, key string) (int64, error) {
	return r.db.Incr(ctx, key).Result()
}

func (r *redisRepository) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.db.Expire(ctx, key, expiration).Err()
}
