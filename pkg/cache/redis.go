package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrKeyNotFound = errors.New("redis: key not found")

type redisCacheService struct {
	rdb *redis.Client
}

func NewRedisCacheService(rdb *redis.Client) RedisCacheService {
	return &redisCacheService{
		rdb: rdb,
	}
}

func (rcs *redisCacheService) Get(ctx context.Context, key string, dest any) error {
	data, err := rcs.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return ErrKeyNotFound
	}
	if err != nil {
		return fmt.Errorf("redis get key %s: %w", key, err)
	}

	if strPtr, ok := dest.(*string); ok {
		*strPtr = data
		return nil
	}

	return json.Unmarshal([]byte(data), dest)
}

func (rcs *redisCacheService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	var data any

	switch v := value.(type) {
	case string, []byte:
		data = v
	default:
		marshaled, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("marshal cache value: %w", err)
		}
		data = marshaled
	}

	return rcs.rdb.Set(ctx, key, data, ttl).Err()
}

func (rcs *redisCacheService) Clear(ctx context.Context, pattern string) error {
	iter := rcs.rdb.Scan(ctx, 0, pattern, 0).Iterator()

	batchSize := 100
	keys := make([]string, 0, batchSize)

	for iter.Next(ctx) {
		keys = append(keys, iter.Val())

		if len(keys) >= batchSize {
			if err := rcs.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
			keys = keys[:0] // Reset slice
		}
	}

	if len(keys) > 0 {
		if err := rcs.rdb.Del(ctx, keys...).Err(); err != nil {
			return err
		}
	}

	return iter.Err()
}

func (rcs *redisCacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := rcs.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
