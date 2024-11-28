package config

import "github.com/ilyakaznacheev/cleanenv"

// PostgresConfig TODO move the configuration structure from this package to the postgres connection package
type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"root"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	DB       string `env:"POSTGRES_DB" env-default:"root"`
}

// RedisConfig TODO move the configuration structure from this package to the redis connection package
type RedisConfig struct {
	Host string `env:"REDIS_HOST" env-default:"localhost"`
	Port int    `env:"REDIS_PORT" env-default:"6379"`
}

type KafkaConfig struct {
	Listener string `env:"KAFKA_ADVERTISED_LISTENERS" env-default:"localhost"`
	Host     string `env:"HOST_KAFKA" env-default:"localhost:9092"`
}

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
}

func NewConfig() *Config {
	var cfg Config
	err := cleanenv.ReadConfig("configs/.env", &cfg)
	if err != nil {
		return nil
	}
	return &cfg
}
