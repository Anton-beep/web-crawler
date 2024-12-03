package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"web-crawler/internal/connection"
)

type KafkaConfig struct {
	Listener string `env:"KAFKA_ADVERTISED_LISTENERS" env-default:"localhost"`
	Host     string `env:"HOST_KAFKA" env-default:"localhost:9092"`
}

type ReceiverConfig struct {
	Port int `env:"RECEIVER_PORT" env-default:"8080"`
}

type Config struct {
	Postgres            connection.PostgresConfig
	Redis               connection.RedisConfig
	Kafka               KafkaConfig
	Receiver            ReceiverConfig
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
		return nil
	}
	return &cfg
}
