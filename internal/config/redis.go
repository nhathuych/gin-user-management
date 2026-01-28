package config

import (
	"context"
	"gin-user-management/internal/util"
	"gin-user-management/pkg/logger"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr     string
	Username string
	Password string
	DB       int
}

func NewRedisClient() *redis.Client {
	cfg := RedisConfig{
		Addr:     util.GetEnv("REDIS_ADDR", "localhost:6379"),
		Username: util.GetEnv("REDIS_USERNAME", ""),
		Password: util.GetEnv("REDIS_PASSWORD", ""),
		DB:       util.GetEnvInt("REDIS_DB", 0),
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Username:     cfg.Username,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     30,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolTimeout:  4 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		logger.AppLogger.Fatal().Err(err).Msg("🔴 Redis connection failed")
	}

	logger.AppLogger.Info().Msg("📦 Redis connected.")

	return client
}
