package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"web-crawler/internal/connection"
)

type KafkaConfig struct {
	Topic        string `env:"TOPIC_NAME"`
	SitesGroupID string `env:"SITES_GROUP_ID"`
	Address      string `env:"ADDRESS_KAFKA"`
	Partition    int    `env:"KAFKA_PARTITION"`
}

type ReceiverConfig struct {
	Port             int    `env:"RECEIVER_PORT" env-default:"8080"`
	Depth            int    `env:"DEFAULT_DEPTH" env-default:"20"`
	MaxNumberOfLinks int    `env:"DEFAULT_MAX_NUMBER_OF_LINKS" env-default:"1000"`
	TempUUID         string `env:"TEMP_UUID" env-default:"00000000-0000-0000-0000-000000000000"`
}

type CollectorConfig struct {
	Tags string `env:"TEXT_TAGS" env-default:""`
}

type Config struct {
	Postgres            connection.PostgresConfig
	Redis               connection.RedisConfig
	Kafka               KafkaConfig
	Receiver            ReceiverConfig
	Collector           CollectorConfig
	RunIntegrationTests bool `env:"RUN_INTEGRATION_TESTS" env-default:"false"`
	Debug               bool `env:"DEBUG" env-default:"true"`
}

func NewConfig(args ...string) *Config {
	path := "configs/.env"
	if len(args) > 0 {
		path = args[0]
	}

	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		errString := fmt.Sprintf("failed to read config: %e", err)
		zap.S().Error(errString)
		panic(errString)
	}
	return &cfg
}
