package integrationTests_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"web-crauler/internal/config"
	"web-crauler/internal/models"
	"web-crauler/internal/repository"
)

var (
	db  models.DataBase
	cfg *config.Config
)

func TestMain(m *testing.M) {
	db = repository.NewDB()
	cfg = config.NewConfig()
	code := m.Run()
	os.Exit(code)
}

func TestCreateAndRetrieveProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	proj1 := models.Project{
		OwnerID:  "1",
		WebGraph: "",
		DlqSites: nil,
	}
	proj1Id, err := db.CreateProject(proj1)
	assert.Equal(t, err, nil, "Creating a new project does not imply an error")
	proj1.ID = proj1Id

	proj1Copy, err := db.GetProject(proj1Id)
	assert.Equal(t, err, nil, fmt.Sprintf("Project with id %s exists", proj1Id))
	assert.Equal(t, *proj1Copy, proj1, fmt.Sprintf("Projects don't match %v != %v", *proj1Copy, proj1))
}

func TestGetProjects(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	expectedProjects := make([]models.Project, 0)
	for i := 0; i < 10; i++ {
		proj := models.Project{
			OwnerID:  fmt.Sprintf("0"),
			WebGraph: "",
			DlqSites: nil,
		}
		projId, err := db.CreateProject(proj)
		assert.Equal(t, err, nil, "Creating a new project does not imply an error")
		proj.ID = projId
		expectedProjects = append(expectedProjects, proj)
	}
	receivedProjects, err := db.GetProjectsByOwnerId("0")
	assert.Equal(t, err, nil, "Getting projects does not imply an error")
	assert.Equal(t, len(receivedProjects), len(expectedProjects), "Number of projects should match")
	for i := 0; i < len(receivedProjects); i++ {
		found := false
		for j := 0; j < len(expectedProjects); j++ {
			if receivedProjects[i].Equals(expectedProjects[j]) {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("Project %v not found in expected projects", receivedProjects[i]))
	}
}

func TestDeleteProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	proj1 := models.Project{
		OwnerID: "1",
	}
	proj1Id, err := db.CreateProject(proj1)
	assert.Equal(t, err, nil, "Creating a new project does not imply an error")

	err = db.DeleteProject(proj1Id)
	assert.Equal(t, err, nil, "Deleting a project does not imply an error")

	_, err = db.GetProject(proj1Id)
	assert.NotEqual(t, err, nil, "Project with id %s doesn't exist", proj1Id)
}

func TestDeleteNonExistentProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	err := db.DeleteProject("123456")
	assert.NotEqual(t, err, nil, "Project with id 123456 doesn't exist")
}

func TestGetNonExistentProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	_, err := db.GetProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "ProjectTemporaryData with id 123456 doesn't exist")
}

func TestDeleteNonExistentProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	err := db.DeleteProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "Deleting non-existent ProjectTemporaryData should return an error")
}

func TestCheckLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	status, err := db.CheckLink("non-existent-link")
	assert.Equal(t, status, false, "Link shouldn't be checked")
	assert.Equal(t, err, nil, "Checking non-existent link shouldn't return an error")
}

func TestUpdateLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	db.UpdateLink("test-link", true)
	exists, err := db.CheckLink("test-link")
	assert.Equal(t, exists, true, "Link should exist after update")
	assert.Equal(t, err, nil, "Checking link shouldn't return an error")
}

func TestCreateAndRetrieveProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	data := &models.ProjectTemporaryData{}
	err := db.CreateProjectTemporaryData("123456", data)
	assert.Equal(t, err, nil, "Updating ProjectTemporaryData should not return an error")

	retrievedData, err := db.GetProjectTemporaryData("123456")
	assert.Equal(t, err, nil, "Retrieving ProjectTemporaryData should not return an error")
	assert.Equal(t, retrievedData, data, "Retrieved ProjectTemporaryData should match the original")
}

func TestDeleteProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	err := db.DeleteProjectTemporaryData("123456")
	assert.Equal(t, err, nil, "Deleting ProjectTemporaryData should not return an error")

	_, err = db.GetProjectTemporaryData("123456")
	assert.NotEqual(t, err, nil, "ProjectTemporaryData with id 123456 should not exist after deletion")
}
