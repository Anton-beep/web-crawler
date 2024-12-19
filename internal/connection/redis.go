package connection

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"time"
	"web-crawler/internal/config"
	"web-crawler/internal/utils"
)

// NewRedisConnect is a function that creates a new connection to a Redis database
func NewRedisConnect(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password, // no password set
		DB:       0,                  // use default DB
	})

	err := utils.RetryTimeout(time.Millisecond*time.Duration(cfg.RetryTimeout), time.Millisecond*time.Duration(cfg.RetryPause), nil, func() error {
		return rdb.Ping(context.Background()).Err()
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	zap.S().Info("Connected to Redis")
	return rdb, nil
}
