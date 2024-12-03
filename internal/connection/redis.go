package connection

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisConfig struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port int    `env:"REDIS_PORT" env-default:"6379"`
}

func NewRedisConnect(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	pong := rdb.Ping(context.Background())
	if pong.Err() != nil {
		zap.S().Fatal(fmt.Errorf("failed to connect to redis: %w", pong.Err()))
	}
	zap.S().Info("Connected to Redis")
	return rdb, nil
}
