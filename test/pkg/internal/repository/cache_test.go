package repository_test

import (
	"experiments/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetNonExistentProjectTemporaryData(t *testing.T) {
	_, err := cache.GetProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "ProjectTemporaryData with id 123456 doesn't exist")
}

func TestDeleteNonExistentProjectTemporaryData(t *testing.T) {
	err := cache.DeleteProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "Deleting non-existent ProjectTemporaryData should return an error")
}

func TestCheckLink(t *testing.T) {
	status, err := cache.CheckLink("non-existent-link")
	assert.Equal(t, status, false, "Link shouldn't be checked")
	assert.Equal(t, err, nil, "Checking non-existent link shouldn't return an error")
}

func TestUpdateLink(t *testing.T) {
	cache.UpdateLink("test-link", true)
	exists, err := cache.CheckLink("test-link")
	assert.Equal(t, exists, true, "Link should exist after update")
	assert.Equal(t, err, nil, "Checking link shouldn't return an error")
}

func TestCreateAndRetrieveProjectTemporaryData(t *testing.T) {
	data := &models.ProjectTemporaryData{}
	err := cache.CreateProjectTemporaryData("123456", data)
	assert.Equal(t, err, nil, "Updating ProjectTemporaryData should not return an error")

	retrievedData, err := cache.GetProjectTemporaryData("123456")
	assert.Equal(t, err, nil, "Retrieving ProjectTemporaryData should not return an error")
	assert.Equal(t, retrievedData, data, "Retrieved ProjectTemporaryData should match the original")
}

func TestDeleteProjectTemporaryData(t *testing.T) {
	err := cache.DeleteProjectTemporaryData("123456")
	assert.Equal(t, err, nil, "Deleting ProjectTemporaryData should not return an error")

	_, err = cache.GetProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "ProjectTemporaryData with id 123456 should not exist after deletion")
}
