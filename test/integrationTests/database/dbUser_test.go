package database_test

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"web-crawler/internal/models"
)

func TestAddUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		t.Skip("skipping integration test")
	}

	user1 := models.User{
		Username: "testAddUser1",
		Email:    "user1@bib.com",
		Password: "password",
	}

	id1, err := db.AddUser(&user1)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id1, "id should not be 0")

	usr, err := db.GetUserByUsername(user1.Username)
	assert.NoError(t, err)
	assert.Equal(t, user1.Username, usr.Username)
}

func TestGetUserByUsername(t *testing.T) {
	if !cfg.RunIntegrationTests {
		t.Skip("skipping integration test")
	}

	user1 := models.User{
		Username: "testGetUser1",
		Email:    "user2@bib.com",
		Password: "password",
	}

	_, err := db.AddUser(&user1)
	assert.NoError(t, err)

	usr, err := db.GetUserByUsername(user1.Username)
	assert.NoError(t, err)
	assert.Equal(t, user1.Username, usr.Username)

	_, err = db.GetUserByUsername("nonexistent")
	assert.Equal(t, models.DataBaseNotFound, err)
}

func TestGetUserByEmail(t *testing.T) {
	if !cfg.RunIntegrationTests {
		t.Skip("skipping integration test")
	}

	user1 := models.User{
		Username: "testGetUserByEmail1",
		Email:    "user3@bib.com",
		Password: "password",
	}

	_, err := db.AddUser(&user1)
	assert.NoError(t, err)

	usr, err := db.GetUserByEmail(user1.Email)
	assert.NoError(t, err)
	assert.Equal(t, user1.Email, usr.Email)

	_, err = db.GetUserByEmail("nonexistent@bib.com")
	assert.Equal(t, models.DataBaseNotFound, err)
}

func TestUpdateUser(t *testing.T) {
	if !cfg.RunIntegrationTests {
		t.Skip("skipping integration test")
	}

	// Add a new user
	user1 := models.User{
		Username: "testUpdateUser1",
		Email:    "user4@bib.com",
		Password: "password",
	}

	id1, err := db.AddUser(&user1)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, id1, "id should not be 0")

	// Update the user's details
	user1.Username = "updatedUsername"
	user1.Email = "updatedEmail@bib.com"
	user1.Password = "newpassword"

	err = db.UpdateUser(&user1)
	assert.NoError(t, err)

	// Retrieve the updated user
	updatedUser, err := db.GetUserByUsername(user1.Username)
	assert.NoError(t, err)
	assert.Equal(t, user1.Username, updatedUser.Username)
	assert.Equal(t, user1.Email, updatedUser.Email)
	assert.Equal(t, user1.Password, updatedUser.Password)
}
