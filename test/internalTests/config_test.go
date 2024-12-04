package internalTests

import (
	"testing"
	"web-crawler/internal/config"

	"github.com/stretchr/testify/assert"
)

func TestNewConfigWithDevConfig(t *testing.T) {
	path := "../../configs/.env"

	cfg := config.NewConfig(path)
	assert.Truef(t, cfg != nil, "Config shouldn't be nil")
	assert.Equal(t, cfg.Postgres.Host, "localhost", "The host of Postgres must be localhost, since the config is intended to run with exposed ports")
	assert.Equal(t, cfg.Kafka.Address, "localhost:9092", "The host of Kafka must be localhost, since the config is intended to run with exposed ports")
	assert.Equal(t, cfg.Redis.Host, "localhost", "The host of Redis must be localhost, since the config is intended to run with exposed ports")
}

func TestNewConfigWithDockerConfig(t *testing.T) {
	path := "../../configs/Docker.env"

	cfg := config.NewConfig(path)
	assert.Truef(t, cfg != nil, "Config shouldn't be nil")
	assert.Equal(t, cfg.Postgres.Host, "postgres", "The host of Postgres must be localhost, since the config is intended to run with exposed ports")
	assert.Equal(t, cfg.Kafka.Address, "kafka:9092", "The host of Kafka must be localhost, since the config is intended to run with exposed ports")
	assert.Equal(t, cfg.Redis.Host, "redis", "The host of Redis must be localhost, since the config is intended to run with exposed ports")
}

func TestNewConfigPanics(t *testing.T) {
	invalidPath := "invalid_path.env"

	defer func() {
		if r := recover(); r != nil {
			assert.Contains(t, r, "failed to read config", "The panic should contain a message about the configuration being unable to be read")
		} else {
			t.Errorf("Panic was expected, but it didn’t happen")
		}
	}()

	config.NewConfig(invalidPath)
}
