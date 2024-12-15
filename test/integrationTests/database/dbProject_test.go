package database_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
	"web-crawler/internal/config"
	"web-crawler/internal/connection"
	"web-crawler/internal/models"
	"web-crawler/internal/repository"
)

var (
	db  models.DataBase
	cfg *config.Config
)

func TestMain(m *testing.M) {
	cfg = config.NewConfig("../../../configs/.env")
	config.InitLogger(true)
	if cfg.RunIntegrationTests {
		db = repository.NewDB(cfg)
	}
	code := m.Run()
	os.Exit(code)
}

func TestCreateProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	project1 := models.Project{
		OwnerID: uuid.New().String(),
	}
	id1, err := db.CreateProject(&project1)
	assert.NotEqual(t, id1, "", "id should not be empty")
	assert.Equal(t, id1, project1.ID, "id should be equal to project id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project2 := models.Project{
		ID:      id1,
		OwnerID: project1.OwnerID,
	}
	id2, err := db.CreateProject(&project2)
	assert.NotEqual(t, id2, "", "creating project should return an id")
	assert.Equal(t, id2, project2.ID, "id should be equal to project id")
	assert.NotEqual(t, id1, id2, "creating project should return a different id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project3 := models.Project{
		OwnerID:  uuid.New().String(),
		DlqSites: []string{"site1", "site2", "site3"},
	}
	id3, err := db.CreateProject(&project3)
	assert.NotEqual(t, id3, "", "creating project should return an id")
	assert.Equal(t, id3, project3.ID, "id should be equal to project id")
	assert.Equal(t, err, nil, "creating project should not return an error")
}

func TestGetProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	ownerId := uuid.New().String()
	project1 := models.Project{
		OwnerID: ownerId,
		Name:    "project1",
	}
	id1, err := db.CreateProject(&project1)
	assert.NotEqual(t, id1, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project2 := models.Project{
		OwnerID: ownerId,
		Name:    "project2",
	}
	id2, err := db.CreateProject(&project2)
	assert.NotEqual(t, id2, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project3 := models.Project{
		OwnerID: ownerId,
		Name:    "project3",
	}
	id3, err := db.CreateProject(&project3)
	assert.NotEqual(t, id3, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project1Cpy, err := db.GetProject(id1)
	assert.Equal(t, err, nil, "getting project should not return an error")
	assert.Equal(t, project1Cpy.ID, id1, "project id should be equal to id1")
	assert.Equal(t, project1Cpy.ID, project1.ID, "project id should be equal to project1 id")
	assert.Equal(t, project1Cpy.OwnerID, ownerId, "project owner id should be equal to ownerId")
	assert.Equal(t, project1Cpy.Name, "project1", "project name should be equal to project1")

	project2Cpy, err := db.GetProject(id2)
	assert.Equal(t, err, nil, "getting project should not return an error")
	assert.Equal(t, project2Cpy.ID, id2, "project id should be equal to id2")
	assert.Equal(t, project2Cpy.ID, project2.ID, "project id should be equal to project2 id")
	assert.Equal(t, project2Cpy.OwnerID, ownerId, "project owner id should be equal to ownerId")
	assert.Equal(t, project2Cpy.Name, "project2", "project name should be equal to project2")
	assert.NotEqual(t, project2Cpy.ID, project1Cpy.ID, "project id should not be equal to project1 id")

	project3Cpy, err := db.GetProject(id3)
	assert.Equal(t, err, nil, "getting project should not return an error")
	assert.Equal(t, project3Cpy.ID, id3, "project id should be equal to id3")
	assert.Equal(t, project3Cpy.ID, project3.ID, "project id should be equal to project3 id")
	assert.Equal(t, project3Cpy.OwnerID, ownerId, "project owner id should be equal to ownerId")
	assert.Equal(t, project3Cpy.Name, "project3", "project name should be equal to project3")
	assert.NotEqual(t, project3Cpy.ID, project1Cpy.ID, "project id should not be equal to project1 id")
	assert.NotEqual(t, project3Cpy.ID, project2Cpy.ID, "project id should not be equal to project2 id")

	project4, err := db.GetProject("non-existing-id")
	assert.NotEqual(t, err, nil, "getting project should return an error")
	assert.Truef(t, project4 == nil, "project should be nil")

	project5, err := db.GetProject("")
	assert.NotEqual(t, err, nil, "getting project should return an error")
	assert.Truef(t, project5 == nil, "project should be nil")

	project6, err := db.GetProject(uuid.New().String())
	assert.NotEqual(t, err, nil, "getting project should return an error")
	assert.Truef(t, project6 == nil, "project should be nil")
}

func TestUpdateProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}
	project1 := models.Project{
		OwnerID:    uuid.New().String(),
		Name:       "name",
		StartUrl:   "url",
		Processing: false,
		WebGraph:   "",
		DlqSites:   nil,
	}

	id1, err := db.CreateProject(&project1)
	assert.NotEqual(t, id1, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project1.Processing = true
	err = db.UpdateProject(&project1)
	assert.Equal(t, err, nil, "updating project should not return an error")

	project1Cpy, err := db.GetProject(id1)
	assert.Equal(t, err, nil, "getting project should not return an error")
	assert.Equal(t, project1Cpy.ID, id1, "project id should be equal to id1")
	assert.Equal(t, project1.Processing, true, "project processing should be true")

	project1.Processing = false
	project1.Name = "new name"
	project1.WebGraph = "graph"
	project1.DlqSites = []string{"site1", "site2", "site3"}
	err = db.UpdateProject(&project1)
	assert.Equal(t, err, nil, "updating project should not return an error")

	project1Cpy, err = db.GetProject(id1)
	assert.Equal(t, err, nil, "getting project should not return an error")
	assert.Equal(t, project1Cpy.ID, id1, "project id should be equal to id1")
	assert.Equal(t, project1.Processing, false, "project processing should be false")
	assert.Equal(t, project1.Name, "new name", "project name should be new name")
	assert.Equal(t, project1.WebGraph, "graph", "project web graph should be graph")
	assert.Equal(t, project1.DlqSites, []string{"site1", "site2", "site3"}, "project dlq sites should be site1, site2, site3")

}

func TestDeleteProject(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	project1 := models.Project{
		OwnerID: uuid.New().String(),
	}
	id1, err := db.CreateProject(&project1)
	assert.NotEqual(t, id1, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	err = db.DeleteProject(id1)
	assert.Equal(t, err, nil, "deleting project should not return an error")

	project1Cpy, err := db.GetProject(id1)
	assert.NotEqual(t, err, nil, "getting project should return an error")
	assert.Truef(t, project1Cpy == nil, "project should be nil")

	err = db.DeleteProject(id1)
	assert.NotEqual(t, err, nil, "deleting project should return an error")

	err = db.DeleteProject("non-existing-id")
	assert.NotEqual(t, err, nil, "deleting project should return an error")

	err = db.DeleteProject("")
	assert.NotEqual(t, err, nil, "deleting project should return an error")

	err = db.DeleteProject(uuid.New().String())
	assert.NotEqual(t, err, nil, "deleting project should return an error")
	assert.EqualError(t, err, models.DataBaseNotFound.Error())
}

func TestGetProjectsByOwnerId(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	ownerId1 := uuid.New().String()
	ownerId2 := uuid.New().String()

	project1 := models.Project{
		OwnerID: ownerId1,
	}
	id1, err := db.CreateProject(&project1)
	assert.NotEqual(t, id1, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project2 := models.Project{
		OwnerID: ownerId1,
	}
	id2, err := db.CreateProject(&project2)
	assert.NotEqual(t, id2, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	project3 := models.Project{
		OwnerID: ownerId2,
	}
	id3, err := db.CreateProject(&project3)
	assert.NotEqual(t, id3, "", "creating project should not return an empty id")
	assert.Equal(t, err, nil, "creating project should not return an error")

	projects1, err := db.GetProjectsByOwnerId(ownerId1)
	assert.Equal(t, err, nil, "getting projects should not return an error")
	assert.Equal(t, len(projects1), 2, "projects should contain 2 elements")

	projects2, err := db.GetProjectsByOwnerId(ownerId2)
	assert.Equal(t, err, nil, "getting projects should not return an error")
	assert.Equal(t, len(projects2), 1, "projects should contain 1 element")

	projects3, err := db.GetProjectsByOwnerId("non-existing-owner-id")
	assert.NotEqual(t, err, nil, "getting projects should not return an error")
	assert.Equal(t, len(projects3), 0, "projects should be empty")

	projects4, err := db.GetProjectsByOwnerId("")
	assert.NotEqual(t, err, nil, "getting projects should not return an error")
	assert.Equal(t, len(projects4), 0, "projects should be empty")

	projects5, err := db.GetProjectsByOwnerId(uuid.New().String())
	assert.NotEqual(t, err, nil, "getting projects should not return an error")
	assert.EqualError(t, err, models.DataBaseNotFound.Error(), "error should be DataBaseNotFound")
	assert.Equal(t, len(projects5), 0, "projects should be empty")
}

func TestSetProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	ptd := models.ProjectTemporaryData{
		Text:   []string{"text"},
		Titles: []string{"titles"},
		Nodes:  "nodes",
		Links:  "links",
	}
	id := uuid.New().String()

	err := db.SetProjectTemporaryData(id, &ptd)
	assert.Equal(t, err, nil, "setting project temporary data should not return an error")

	ptd.Text[0] = "new text"

	err = db.SetProjectTemporaryData(id, &ptd)
	assert.Equal(t, err, nil, "setting project temporary data should not return an error")
}

func TestGetProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	ptd := models.ProjectTemporaryData{
		Text:   []string{"text"},
		Titles: []string{"titles"},
		Nodes:  "nodes",
		Links:  "links",
	}
	id := uuid.New().String()

	err := db.SetProjectTemporaryData(id, &ptd)
	assert.Equal(t, err, nil, "setting project temporary data should not return an error")

	ptdCpy, err := db.GetProjectTemporaryData(id)
	assert.Equal(t, err, nil, "getting project temporary data should not return an error")
	assert.Equal(t, ptdCpy.Text, ptd.Text, "text should be equal to text")
	assert.Equal(t, ptdCpy.Titles, ptd.Titles, "titles should be equal to titles")
	assert.Equal(t, ptdCpy.Nodes, ptd.Nodes, "nodes should be equal to nodes")
	assert.Equal(t, ptdCpy.Links, ptd.Links, "links should be equal to links")

	ptdCpy2, err := db.GetProjectTemporaryData("non-existing-id")
	assert.NotEqual(t, err, nil, "getting project temporary data should return an error")
	assert.Truef(t, ptdCpy2 == nil, "project temporary data should be nil")
	assert.EqualError(t, err, models.DataBaseNotFound.Error(), "error should be DataBaseNotFound")

	ptd.Text = []string{"new text"}
	err = db.SetProjectTemporaryData(id, &ptd)
	assert.Equal(t, err, nil, "setting project temporary data should not return an error")

	ptdCpy3, err := db.GetProjectTemporaryData(id)
	assert.Equal(t, err, nil, "getting project temporary data should not return an error")
	assert.Equal(t, ptdCpy3.Text[0], "new text", "text should be equal to new text")
	assert.Equal(t, ptdCpy3.Titles, ptd.Titles, "titles should be equal to titles")
	assert.Equal(t, ptdCpy3.Nodes, ptd.Nodes, "nodes should be equal to nodes")
	assert.Equal(t, ptdCpy3.Links, ptd.Links, "links should be equal to links")
}

func TestDeleteProjectTemporaryData(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	ptd := models.ProjectTemporaryData{
		Text:   []string{"text"},
		Titles: []string{"titles"},
		Nodes:  "nodes",
		Links:  "links",
	}
	id := uuid.New().String()

	err := db.SetProjectTemporaryData(id, &ptd)
	assert.Equal(t, err, nil, "setting project temporary data should not return an error")

	err = db.DeleteProjectTemporaryData(id)
	assert.Equal(t, err, nil, "deleting project temporary data should not return an error")

	ptdCpy, err := db.GetProjectTemporaryData(id)
	assert.NotEqual(t, err, nil, "getting project temporary data should return an error")
	assert.Truef(t, ptdCpy == nil, "project temporary data should be nil")
	assert.EqualError(t, err, models.DataBaseNotFound.Error(), "error should be DataBaseNotFound")

	err = db.DeleteProjectTemporaryData(id)
	assert.NotEqual(t, err, nil, "deleting project temporary data should return an error")
	assert.EqualError(t, err, models.DataBaseNotFound.Error(), "error should be DataBaseNotFound")

	_, err = db.GetProjectTemporaryData(id)
	assert.Error(t, err, models.DataBaseNotFound)
}

func TestUpdateLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	link := "linkUpdateLink" + time.Now().String()

	visited, err := db.CheckSlug(link)
	assert.Equal(t, err, nil, "checking link should not return an error")
	assert.Equal(t, visited, false, "link should not be visited")

	err = db.UpdateSlug(link, true)
	assert.Equal(t, err, nil, "updating link should not return an error")

	visited, err = db.CheckSlug(link)
	assert.Equal(t, err, nil, "checking link should not return an error")
	assert.Equal(t, visited, true, "link should be visited")

	for i := 1; i < 10; i++ {
		err = db.UpdateSlug(link, i%2 == 0)
		assert.Equal(t, err, nil, "updating link should not return an error")
	}

	visited, err = db.CheckSlug(link)
	assert.Equal(t, err, nil, "checking link should not return an error")
	assert.Equal(t, visited, false, "link should not be visited")
}

func TestCheckLink(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	link := "linkCheckLink" + time.Now().String()

	visited, err := db.CheckSlug(link)
	assert.Equal(t, err, nil, "checking link should not return an error")
	assert.Equal(t, visited, false, "link should not be visited")

	err = db.UpdateSlug(link, true)
	assert.Equal(t, err, nil, "updating link should not return an error")

	visited, err = db.CheckSlug(link)
	assert.Equal(t, err, nil, "checking link should not return an error")
	assert.Equal(t, visited, true, "link should be visited")
}

func TestWrongConnection(t *testing.T) {
	cfg := config.Config{
		Postgres: connection.PostgresConfig{
			Host:     "non-existing-host",
			Port:     5432,
			User:     "",
			Password: "",
			DB:       "",
		},
		Redis: connection.RedisConfig{
			Host: "non-existing-host",
			Port: 6379,
		},
		RunIntegrationTests: false,
		Debug:               false,
	}
	_, err := connection.NewPostgresConnect(cfg.Postgres)
	assert.False(t, err == nil)
	_, err = connection.NewRedisConnect(cfg.Redis)
	assert.False(t, err == nil)
}

func TestPush2Queue(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	// Structure that supports conversion to json
	err := db.Push2Queue("key1", struct {
		Field1 string `json:"field_1"`
		Field2 string `json:"field_2"`
		Field3 string `json:"field_3"`
	}{"a", "b", "c"})
	assert.True(t, err == nil)
}

func TestPopFromQueue(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	err := db.Push2Queue("key2", struct {
		Field1 string `json:"field_1"`
		Field2 string `json:"field_2"`
		Field3 string `json:"field_3"`
	}{"a", "b", "c"})
	assert.True(t, err == nil)

	_, err = db.PopFromQueue("key2")
	assert.True(t, err == nil)
}

func TestAddAnalyserTask(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	err := db.AddAnalyserTask("id", "type")
	assert.True(t, err == nil)
	_, _ = db.GetAnalyserTask()
}

func TestGetAnalyserTask(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	_, err := db.GetAnalyserTask()
	assert.True(t, errors.Is(err, models.DataBaseQueueIsEmpty))
	err = db.AddAnalyserTask("id", "type")
	assert.True(t, err == nil)
	task, err := db.GetAnalyserTask()
	assert.True(t, err == nil)
	assert.Equal(t, task, models.AnalyserTask{
		ID:   "id",
		Type: "type",
	})
}

func TestUseWrongId(t *testing.T) {
	if !cfg.RunIntegrationTests {
		return
	}

	_, err := db.GetProjectMaxDepth("wrong-uuid")
	assert.Error(t, err, models.DataBaseWrongID)
	_, err = db.GetProjectMaxDepth(uuid.New().String())
	assert.Error(t, err, models.DataBaseNotFound)

	err = db.CheckCollectorCounter("wrong-uuid")
	assert.Error(t, err, models.DataBaseWrongID)
	err = db.CheckCollectorCounter(uuid.New().String())
	assert.Error(t, err, models.DataBaseNotFound)

	err = db.UpdateProject(&models.Project{ID: "wrong-uuid"})
	assert.Error(t, err, models.DataBaseWrongID)
	err = db.UpdateProject(&models.Project{ID: uuid.New().String()})
	assert.Error(t, err, models.DataBaseNotFound)
}
