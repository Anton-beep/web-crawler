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

// DataBase is a struct that contains the
// connection to the Postgres and Redis databases
type DataBase struct {
	postgres         *sqlx.DB
	redis            *redis.Client
	analyserQueueKey string
}

// GetProjectMaxDepth is a function that returns
// the maximum depth of the project
//
// params:
// - id: id of project ocf type uuid
//
// returns:
//   - int: the maximum depth of the project
//   - error: an error if the project with the
//     given id doesn't exist
func (d *DataBase) GetProjectMaxDepth(id string) (int, error) {
	err := d.checkIfIdExists(id)
	if err != nil {
		return 0, err
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("max_depth").
		From("projects").
		Where(sq.Eq{"id": id})

	ProjectQueryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query for project", err)
		return 0, fmt.Errorf("failed to build query for project: %w", err)
	}

	var depth int
	err = d.postgres.Get(&depth, ProjectQueryString, args...)
	if err != nil {
		return 0, err
	}
	return depth, nil
}

// CheckCollectorCounter is a function that
// checks if the collector counter is negative
//
// params:
// - id: id of project ocf type uuid
//
// returns:
//   - error: an error if the project with the
//     given id doesn't exist or the collector counter is negative
func (d *DataBase) CheckCollectorCounter(id string) error {
	err := d.checkIfIdExists(id)
	if err != nil {
		return err
	}
	ptd, err := d.GetProjectTemporaryData(id)
	if err != nil {
		return err
	}
	if ptd.TotalCollectorCounter <= 0 {
		return models.CollectorCounterIsNegative
	}
	err = d.SetProjectTemporaryData(id, ptd)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataBase) checkIfIdExists(id string) error {
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

// GetProject is a function that returns the project by id
//
// params:
// - id: id of project ocf type uuid
//
// returns:
// - *Project: the project with the given id
func (d *DataBase) GetProject(id string) (*models.Project, error) {
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
		"id", "owner_id", "name", "start_url", "processing", "web_graph", "max_depth", "max_number_of_links").
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

// GetProjectTemporaryData is a function that returns
// the temporary data of the project by id
//
// params:
// - id: id of project ocf type uuid
//
// returns:
//   - *ProjectTemporaryData: the temporary data of
//     the project with the given id
func (d *DataBase) GetProjectTemporaryData(id string) (*models.ProjectTemporaryData, error) {
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

// CreateProject is a function that creates a new project
//
// params:
// - project: the project to create
//
// returns:
// - string: the id of the created project
func (d *DataBase) CreateProject(project *models.Project) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Insert("projects").
		Columns("id", "owner_id", "name", "start_url", "processing", "web_graph", "dlq_sites", "max_depth", "max_number_of_links").
		Values(sq.Expr("gen_random_uuid()"), project.OwnerID, project.Name, project.StartUrl, project.Processing, project.WebGraph, pq.Array(project.DlqSites), project.MaxDepth, project.MaxNumberOfLinks).
		Suffix("RETURNING id")
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Debug("failed to build query", err)
		return "", fmt.Errorf("failed to build query: %w", err)
	}

	var generatedID string
	err = d.postgres.QueryRow(queryString, args...).Scan(&generatedID)
	if err != nil {
		zap.S().Debug("failed to create project", err)
		return "", fmt.Errorf("failed to create project: %w", err)
	}

	if generatedID == "" {
		return "", fmt.Errorf("failed to create project (id is \"\"): %w", models.DataBaseNotFound)
	}
	project.ID = generatedID

	return generatedID, nil
}

// SetProjectTemporaryData is a function that sets
// the temporary data of the project
//
// params:
// - id: id of project ocf type uuid
// - data: the temporary data of the project
//
// returns:
// - error: an error if the temporary data of the project
func (d *DataBase) SetProjectTemporaryData(id string, data *models.ProjectTemporaryData) error {
	val, err := json.Marshal(data)
	if err != nil {
		zap.S().Debug("failed to serialize project temporary data", err)
		return fmt.Errorf("failed to serialize project temporary data: %w", err)
	}

	d.redis.Set(context.Background(), id, val, 0)
	return nil
}

// UpdateProject is a function that updates the project
//
// params:
// - project: the project to update
//
// returns:
// - error: an error if the project with the given id doesn't exist
func (d *DataBase) UpdateProject(project *models.Project) error {
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

// DeleteProject is a function that deletes the project
//
// params:
// - id: id of project ocf type uuid
//
// returns:
// - error: an error if the project with the given id doesn't exist
func (d *DataBase) DeleteProject(id string) error {
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

// DeleteProjectTemporaryData is a function that deletes
// the temporary data of the project
//
// params:
// - id: id of project ocf type uuid
//
// returns:
// - error: an error if the temporary data of the project
func (d *DataBase) DeleteProjectTemporaryData(id string) error {
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

// GetProjectsByOwnerId is a function that returns
// the projects by owner id
//
// params:
// - ownerId: id of owner ocf type uuid
//
// returns:
// - []*ShortProject: the projects with the given owner id
func (d *DataBase) GetProjectsByOwnerId(ownerId string) ([]*models.ShortProject, error) {
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

// CheckSlug is a function that checks if the slug exists
//
// params:
// - slag: the slug to check
//
// returns:
// - bool: true if the slug exists, false otherwise
func (d *DataBase) CheckSlug(slag string) (bool, error) {
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

// UpdateSlug is a function that updates the slug
//
// params:
// - slag: the slug to update
// - status: new status of the slug
//
// returns:
// - error: an error if the slug wasn't updated
func (d *DataBase) UpdateSlug(slag string, status bool) error {
	err := d.redis.Set(context.Background(), slag, status, 0).Err()
	return err
}

// Push2Queue adds an arbitrary structure that supports
// being converted to json to the queue with key
//
// params:
// - key: key of queue
// - value: that supports json conversion
//
// returns:
// - error
func (d *DataBase) Push2Queue(key string, value interface{}) error {
	res, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = d.redis.LPush(context.Background(), key, res).Err()
	if err != nil {
		err = fmt.Errorf("failed to add value %v to Redis Queue: %s", value, err)
		return err
	}
	return nil
}

// PopFromQueue pops first value from queue by key
//
// params:
// - key: key of queue
//
// returns:
//   - string: first value from queue
//     if queue was empty, equals ""
//   - error: error received when retrieving value
func (d *DataBase) PopFromQueue(key string) (string, error) {
	res := d.redis.RPop(context.Background(), key)
	return res.Val(), res.Err()
}

// AddAnalyserTask adds task to analyser queue
//
// params:
// - projectId: uuid of Project
// - typeOfAnalysis: type of analysis
//
// returns:
// - error if adding fails
func (d *DataBase) AddAnalyserTask(projectId, typeOfAnalysis string) error {
	return d.Push2Queue(d.analyserQueueKey, models.AnalyserTask{
		ID:   projectId,
		Type: typeOfAnalysis,
	})
}

// GetAnalyserTask checks the task queue for analysers
// If the queue is empty, returns an empty structure and DataBaseQueueIsEmpty
func (d *DataBase) GetAnalyserTask() (models.AnalyserTask, error) {
	val, err := d.PopFromQueue(d.analyserQueueKey)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return models.AnalyserTask{}, models.DataBaseQueueIsEmpty
		}
		return models.AnalyserTask{}, err
	}
	if val == "" {
		return models.AnalyserTask{}, models.DataBaseQueueIsEmpty
	}
	var analTask models.AnalyserTask
	err = json.Unmarshal([]byte(val), &analTask)
	if err != nil {
		return models.AnalyserTask{}, err
	}
	return analTask, err
}

// NewDB is a function that creates a new DataBase struct
// DataBase implements the models.DataBase interface
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
		postgres:         postgresConnect,
		redis:            redisConnect,
		analyserQueueKey: cfg.Redis.AnalyserQueueKey,
	}
}
