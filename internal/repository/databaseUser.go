package repository

import (
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
	"web-crawler/internal/models"
)

// GetUserByUsername is a function that retrieves a user by their username
//
// params:
// - username: the username of the user
//
// returns:
// - *models.User: the user with the given username
// - error: an error if the user with the given username doesn't exist
func (d DataBase) GetUserByUsername(username string) (*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("id", "username", "email", "password").From("users").Where(sq.Eq{"username": username})
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("failed to build query: ", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user models.User
	err = d.postgres.Get(&user, queryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.DataBaseNotFound
		}
		return nil, fmt.Errorf("failed to execute check query: %w", err)
	}

	return &user, nil
}

// GetUserByEmail is a function that retrieves a user by their email
//
// params:
// - email: the email of the user
//
// returns:
// - *models.User: the user with the given email
// - error: an error if the user with the given email doesn't exist
func (d DataBase) GetUserByEmail(email string) (*models.User, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("id", "username", "email", "password").From("users").Where(sq.Eq{"email": email})
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("failed to build query: ", err)
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	var user models.User
	err = d.postgres.Get(&user, queryString, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.DataBaseNotFound
		}
		return nil, fmt.Errorf("failed to execute check query: %w", err)
	}

	return &user, nil
}

// AddUser is a function that adds a new user to the database
//
// params:
// - user: the user to add
//
// returns:
// - string: the ID of the newly created user
// - error: an error if the user wasn't added
func (d DataBase) AddUser(user *models.User) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Insert("users").
		Columns("id", "username", "email", "password").
		Values(sq.Expr("gen_random_uuid()"), user.Username, user.Email, user.Password).
		Suffix("RETURNING id")
	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("failed to build query: ", err)
		return "", fmt.Errorf("failed to build query: %w", err)
	}

	var generatedID string
	err = d.postgres.QueryRow(queryString, args...).Scan(&generatedID)
	if err != nil {
		zap.S().Error("failed to query: ", err)
		return "", fmt.Errorf("failed to query: %w", err)
	}

	if generatedID == "" {
		return "", fmt.Errorf("failed to create user (id is \"\"): %w", models.DataBaseNotFound)
	}
	user.ID = generatedID

	return generatedID, nil
}

// UpdateUser is a function that updates an existing user in the database
//
// params:
// - user: the user to update
//
// returns:
// - error: an error if the user wasn't updated
func (d DataBase) UpdateUser(user *models.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Update("users").
		Set("username", user.Username).
		Set("email", user.Email).
		Set("password", user.Password).
		Where(sq.Eq{"id": user.ID})

	queryString, args, err := query.ToSql()
	if err != nil {
		zap.S().Error("failed to build query: ", err)
		return fmt.Errorf("failed to build query: %w", err)
	}

	_, err = d.postgres.Exec(queryString, args...)
	if err != nil {
		zap.S().Error("failed to exec: ", err)
		return fmt.Errorf("failed to exec: %w", err)
	}

	return nil
}
