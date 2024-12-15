package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
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
	SecretSignature  string `env:"SECRET_SIGNATURE" env-default:"thisIsADefaultSignatureIfYouSeeItInYourCodeYouBetterChangeIt"`
}

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"root"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	DB       string `env:"POSTGRES_DB" env-default:"root"`
}

type RedisConfig struct {
	Host             string `env:"REDIS_HOST" env-default:"localhost"`
	Port             int    `env:"REDIS_PORT" env-default:"6379"`
	AnalyserQueueKey string `env:"ANALYSER_QUEUE_KEY" env-default:"analyser-queue"`
}

type CollectorConfig struct {
	Tags string `env:"TEXT_TAGS" env-default:""`
}

// Config is a struct that contains the configuration for the application
type Config struct {
	Postgres            PostgresConfig
	Redis               RedisConfig
	Kafka               KafkaConfig
	Receiver            ReceiverConfig
	Collector           CollectorConfig
	RunIntegrationTests bool `env:"RUN_INTEGRATION_TESTS" env-default:"false"`
	Debug               bool `env:"DEBUG" env-default:"true"`
	RetryPause          int  `env:"RETRY_PAUSE" env-default:"1000"`
	RetryAttempts       int  `env:"RETRY_COUNT" env-default:"3"`
	RetryTimeout        int  `env:"RETRY_TIMEOUT" env-default:"10000"`
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
