package repository_test

import (
	"experiments/internal/models"
	"experiments/internal/repository"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var (
	db    models.LongTermDataBase
	cache models.ShortTermDataBase
)

func TestMain(m *testing.M) {
	db = repository.NewDB()
	cache = repository.NewCache()
	code := m.Run()
	os.Exit(code)
}

func TestCreateAndRetrieveProject(t *testing.T) {
	proj1 := models.Project{
		OwnerID: "1",
	}
	proj1Id, err := db.CreateProject(proj1)
	assert.Equal(t, err, nil, "Creating a new project does not imply an error")
	proj1.ID = proj1Id

	proj1Copy, err := db.GetProject(proj1Id)
	assert.Equal(t, err, nil, fmt.Sprintf("Project with id %s exists", proj1Id))
	assert.Equal(t, *proj1Copy, proj1, fmt.Sprintf("Projects don't match %v != %v", *proj1Copy, proj1))
}

func TestGetProjects(t *testing.T) {
	expectedProjects := make([]models.Project, 0)
	for i := 0; i < 10; i++ {
		proj := models.Project{
			OwnerID: fmt.Sprintf("0"),
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
			if *receivedProjects[i] == expectedProjects[j] {
				found = true
				break
			}
		}
		assert.True(t, found, fmt.Sprintf("Project %v not found in expected projects", receivedProjects[i]))
	}
}

func TestCreateAndRetrieveDashboard(t *testing.T) {
	proj1 := models.Project{
		OwnerID: "1",
	}
	proj1Id, err := db.CreateProject(proj1)
	assert.Equal(t, err, nil, "Creating a new project does not imply an error")
	proj1.ID = proj1Id

	dash1 := models.Dashboard{
		ProjectID: proj1Id,
		WebGraph:  "",
		DlqSites:  nil,
	}
	dash1Id, err := db.CreateDashboard(dash1)
	assert.Equal(t, err, nil, fmt.Sprintf("Creating a new dashboard does not imply an error"))
	dash1.ID = dash1Id

	dash1Copy, err := db.GetDashboard(dash1Id)
	assert.Equal(t, err, nil, fmt.Sprintf("Dashboard with id %s exists", dash1Id))
	assert.Equal(t, *dash1Copy, dash1, fmt.Sprintf("Dashboards don't match %v != %v", *dash1Copy, dash1))
}

func TestDeleteProject(t *testing.T) {
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

func TestDeleteDashboard(t *testing.T) {
	proj1 := models.Project{
		OwnerID: "1",
	}
	proj1Id, err := db.CreateProject(proj1)
	assert.Equal(t, err, nil, "Creating a new project does not imply an error")

	dash1 := models.Dashboard{
		ProjectID: proj1Id,
		WebGraph:  "",
		DlqSites:  nil,
	}
	dash1Id, err := db.CreateDashboard(dash1)
	assert.Equal(t, err, nil, fmt.Sprintf("Creating a new dashboard does not imply an error"))

	err = db.DeleteDashBoard(dash1Id)
	assert.Equal(t, err, nil, "Deleting a dashboard does not imply an error")

	_, err = db.GetDashboard(dash1Id)
	assert.NotEqual(t, err, nil, "Dashboard with id %s doesn't exist", dash1)
}

func TestDeleteNonExistentProject(t *testing.T) {
	err := db.DeleteProject("123456")
	assert.NotEqual(t, err, nil, "Project with id 123456 doesn't exist")
}

func TestDeleteNonExistentDashboard(t *testing.T) {
	err := db.DeleteDashBoard("123456")
	assert.NotEqual(t, err, nil, "Dashboard with id 123456 doesn't exist")
}
