package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"alloy/internal/shared/config"

	"go.uber.org/zap"
)

var ctx = context.Background()

func GetRedisClient(cfg *config.Config, zapLogger *zap.Logger, db int) (*redis.Client, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       db,
	})

	if cfg.RedisUrl != "" {
		opt, _ := redis.ParseURL(cfg.RedisUrl)
		rdb = redis.NewClient(opt)
	}

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	zapLogger.Info("Connected to Redis!")
	return rdb, err
}
