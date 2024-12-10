package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
	"web-crawler/internal/connection"
)

type KafkaConfig struct {
	SitesTopic     string `env:"SITES_TOPIC_NAME"`
	AnalyseTopic   string `env:"ANALYSE_TOPIC_NAME"`
	SitesGroupID   string `env:"SITES_GROUP_ID"`
	AnalyseGroupID string `env:"ANALYSE_GROUP_ID"`
	Address        string `env:"ADDRESS_KAFKA"`
	Partition      int    `env:"KAFKA_PARTITION"`
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

// Config is a struct that contains the configuration for the application
type Config struct {
	Postgres            connection.PostgresConfig
	Redis               connection.RedisConfig
	Kafka               KafkaConfig
	Receiver            ReceiverConfig
	Collector           CollectorConfig
	RunIntegrationTests bool `env:"RUN_INTEGRATION_TESTS" env-default:"false"`
	Debug               bool `env:"DEBUG" env-default:"true"`
}

// NewConfig is a function that creates a new Config struct
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
