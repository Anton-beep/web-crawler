package connection

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type PostgresConfig struct {
	Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
	Port     int    `env:"POSTGRES_PORT" env-default:"5432"`
	User     string `env:"POSTGRES_USER" env-default:"root"`
	Password string `env:"POSTGRES_PASSWORD" env-default:"123"`
	DB       string `env:"POSTGRES_DB" env-default:"root"`
}

// NewPostgresConnect is a function that creates a new connection to a Postgres database
func NewPostgresConnect(config PostgresConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.DB,
	)

	connect, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to connect to Postgres: %w", err)
	}

	_, err = connect.Conn(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to connect to Postgres: %w", err)
	}

	zap.S().Info("Connected to Postgres")

	err = CreateProjectTable(connect)
	if err != nil {
		return nil, fmt.Errorf("failed to create table in postgres: %w", err)
	}
	return connect, nil
}

// CreateProjectTable is a function that creates the table 'projects' in the Postgres database
func CreateProjectTable(db *sqlx.DB) error {
	queryProjects := `
	CREATE TABLE IF NOT EXISTS projects (
		id UUID PRIMARY KEY,
		owner_id UUID NOT NULL,
		name TEXT NOT NULL,
		start_url TEXT NOT NULL,
		processing BOOLEAN NOT NULL,
		web_graph TEXT,
		dlq_sites TEXT[],
		max_depth INT,
		max_number_of_links INT
	);
	`

	zap.S().Info("Table 'projects' created or already exists.")

	_, err := db.ExecContext(context.Background(), queryProjects)
	if err != nil {
		return err
	}

	queryUsers := `
	CREATE TABLE IF NOT EXISTS users (
	    		id UUID PRIMARY KEY,
	    		username TEXT NOT NULL,
	    		email TEXT NOT NULL,
	    		password TEXT NOT NULL
	);
	`

	_, err = db.ExecContext(context.Background(), queryUsers)
	if err != nil {
		return err
	}

	zap.S().Info("Table 'users' created or already exists.")

	return nil
}
