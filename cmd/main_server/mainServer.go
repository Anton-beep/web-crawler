package main

import (
	"context"
	"experiments/internal/config"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"net/http"
)

func main() {
	cfg := config.NewConfig()
	fmt.Println(cfg)

	// Try to connect to Postgres
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.DB)
	_, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to connect to Postgres: %w", err))
		panic(err)
	} else {
		fmt.Println("Successfully connected to Postgres")
	}

	// Try to connect to Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: "",
		DB:       0,
	})
	// Ping Redis to test the connection
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println(fmt.Errorf("failed to connect to Redis: %w", err))
		panic(err)
	} else {
		fmt.Println("Successfully connected to Redis:", pong)
	}

	// Try to connect to Kafka
	// to produce messages
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", cfg.Kafka.Host, topic, partition)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to dial leader: %w", err))
		panic(err)
	}

	if err := conn.Close(); err != nil {
		fmt.Println(fmt.Errorf("failed to close writer: %w", err))
		panic(err)
	}

	fmt.Println("Successfully connect to Kafka and closed connection")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello World!")
	})
	http.ListenAndServe(":8080", nil)
}
