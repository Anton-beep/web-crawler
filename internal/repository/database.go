package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"web-crawler/internal/config"
	"web-crawler/internal/connection"
	"web-crawler/internal/models"
)

type DataBase struct {
	postgres *sqlx.DB
	redis    *redis.Client
}

func (d DataBase) checkIfIdExists(id string) error {
	if uuid.Validate(id) != nil {
		return models.DataBaseWrongID
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	checkQuery := psql.Select("id").From("projects").Where(sq.Eq{"id": id})
	checkQueryString, args, err := checkQuery.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build check query: %w", err)
	}
	var _id string
	err = d.postgres.Get(&_id, checkQueryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DataBaseNotFound
		}
		return fmt.Errorf("failed to execute check query: %w", err)
	}
	return nil
}

func (d DataBase) GetProject(id string) (*models.Project, error) {
	err := d.checkIfIdExists(id)
	if err != nil {
		return nil, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	DlqSitesQuery := psql.Select("dlq_sites").From("projects").Where(sq.Eq{"id": id})
	DlqSitesQueryString, args, err := DlqSitesQuery.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query for dlq sites", err)
		return nil, fmt.Errorf("failed to build query for dlq sites: %w", err)
	}

	var dlqSites []string
	err = d.postgres.QueryRow(DlqSitesQueryString, args...).Scan(pq.Array(&dlqSites))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.DataBaseNotFound
		}
		zap.S().Debug("failed to get project dlq sites", err)
		return nil, fmt.Errorf("failed to get project dlq sites: %w", err)
	}

	ProjectQuery := psql.Select(
		"id", "owner_id", "name", "start_url", "processing", "web_graph").
		From("projects").
		Where(sq.Eq{"id": id})

	ProjectQueryString, args, err := ProjectQuery.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query for project", err)
		return nil, fmt.Errorf("failed to build query for project: %w", err)
	}

	var project models.Project
	err = d.postgres.Get(&project, ProjectQueryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.DataBaseNotFound
		}
		zap.S().Debug("failed to get project", err)
		return nil, fmt.Errorf("failed to get project: %w", err)
	}
	project.DlqSites = dlqSites

	return &project, nil
}

func (d DataBase) GetProjectTemporaryData(id string) (*models.ProjectTemporaryData, error) {
	val, err := d.redis.Get(context.Background(), id).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, models.DataBaseNotFound
		}
		zap.S().Debug("failed to get project temporary data", err)
		return nil, fmt.Errorf("failed to get project temporary data: %w", err)
	}

	var ptd models.ProjectTemporaryData
	err = json.Unmarshal([]byte(val), &ptd)
	if err != nil {
		zap.S().Debug("failed to deserialize project temporary data", err)
		return nil, fmt.Errorf("failed to deserialize project temporary data: %w", err)
	}

	return &ptd, nil
}

func (d DataBase) CreateProject(project *models.Project) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	project.ID = uuid.New().String()
	query := psql.Insert("projects").Columns("id", "owner_id", "name", "start_url", "processing", "web_graph", "dlq_sites").
		Values(project.ID, project.OwnerID, project.Name, project.StartUrl, project.Processing, project.WebGraph, pq.Array(project.DlqSites))
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query", err)
		return "", fmt.Errorf("failed to build query: %w", err)
	}

	_, err = d.postgres.Exec(queryString, args...)

	if err != nil {
		zap.S().Debug("failed to create project", err)
		return "", fmt.Errorf("failed to create project: %w", err)
	}
	return project.ID, nil
}

func (d DataBase) SetProjectTemporaryData(id string, data *models.ProjectTemporaryData) error {
	val, err := json.Marshal(data)
	if err != nil {
		zap.S().Debug("failed to serialize project temporary data", err)
		return fmt.Errorf("failed to serialize project temporary data: %w", err)
	}

	d.redis.Set(context.Background(), id, val, 0)
	return nil
}

func (d DataBase) UpdateProject(project *models.Project) error {
	err := d.checkIfIdExists(project.ID)
	if err != nil {
		return err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Update("projects").
		SetMap(map[string]interface{}{
			"owner_id":   project.OwnerID,
			"name":       project.Name,
			"start_url":  project.StartUrl,
			"processing": project.Processing,
			"web_graph":  project.WebGraph,
			"dlq_sites":  pq.Array(project.DlqSites),
		}).
		Where(sq.Eq{"id": project.ID})

	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Debug("failed to build update query", err)
		return fmt.Errorf("failed to build update query: %w", err)
	}

	_, err = d.postgres.Exec(queryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.DataBaseNotFound
		}
		zap.S().Debug("failed to execute update query", err)
		return fmt.Errorf("failed to execute update query: %w", err)
	}

	return nil
}

func (d DataBase) DeleteProject(id string) error {
	err := d.checkIfIdExists(id)
	if err != nil {
		return err
	}
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	DeleteQuery := psql.Delete("projects").Where(sq.Eq{"id": id})
	DeleteQueryString, args, err := DeleteQuery.ToSql()
	if err != nil {
		zap.S().Debug("failed to build delete query", err)
		return fmt.Errorf("failed to build delete query: %w", err)
	}

	_, err = d.postgres.Exec(DeleteQueryString, args...)
	if err != nil {
		zap.S().Debug("failed to delete project", err)
		return fmt.Errorf("failed to delete project: %w", err)
	}

	return nil
}

func (d DataBase) DeleteProjectTemporaryData(id string) error {
	_, err := d.GetProjectTemporaryData(id)
	if err != nil {
		return err
	}
	err = d.redis.Del(context.Background(), id).Err()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.DataBaseNotFound
		}
		zap.S().Debug("failed to delete project temporary data", err)
		return fmt.Errorf("failed to delete project temporary data: %w", err)
	}

	return nil
}

func (d DataBase) GetProjectsByOwnerId(ownerId string) ([]*models.ShortProject, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("id", "name").From("projects").Where(sq.Eq{"owner_id": ownerId})
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := d.postgres.Query(queryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.DataBaseNotFound
		}
		zap.S().Debug("failed to execute query", err)
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	var projects []*models.ShortProject
	for rows.Next() {
		var project models.ShortProject
		err = rows.Scan(&project.ID, &project.Name)
		if err != nil {
			zap.S().Debug("failed to scan row", err)
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		projects = append(projects, &project)
	}

	if len(projects) == 0 {
		return nil, models.DataBaseNotFound
	}
	return projects, nil
}

func (d DataBase) CheckLink(slag string) (bool, error) {
	val, err := d.redis.Get(context.Background(), slag).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return false, nil
		}
		zap.S().Debug("failed to get link state", err)
		return false, fmt.Errorf("failed to get link state: %w", err)
	}
	return val == "1", nil
}

func (d DataBase) UpdateLink(slag string, status bool) error {
	err := d.redis.Set(context.Background(), slag, status, 0).Err()
	return err
}

func NewDB(cfg *config.Config) models.DataBase {
	// 🤓🤓🤓
	postgresConnect, err := connection.NewPostgresConnect(cfg.Postgres)
	if err != nil {
		zap.S().Fatal(fmt.Errorf("database wan't created due to error in postgres connect: %w", err))
	}
	redisConnect, err := connection.NewRedisConnect(cfg.Redis)
	if err != nil {
		zap.S().Fatal(fmt.Errorf("database wan't created due to error in redis connect: %w", err))
	}
	zap.S().Info("Database created")
	return &DataBase{
		postgres: postgresConnect,
		redis:    redisConnect,
	}
}
