package receiver

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"web-crawler/internal/models"
)

type errMsg struct {
	Message string `json:"message"`
}

// Pong is a simple health check handler
func Pong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

type inCreate struct {
	Name          string `json:"name"`
	StartUrl      string `json:"start_url"`
	NumberOfLinks int    `json:"number_of_links"`
	Depth         int    `json:"depth"`
}

type outCreate struct {
	Id string `json:"id"`
}

// CreateProject creates a new project
func (r *Service) CreateProject(c echo.Context) error {
	var in inCreate

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	if in.Name == "" || in.StartUrl == "" || in.NumberOfLinks < 1 || in.Depth < 1 {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid json"})
	}

	user := c.Get("user").(*models.User)

	prj := models.Project{
		Name:             in.Name,
		StartUrl:         in.StartUrl,
		OwnerID:          user.ID,
		MaxNumberOfLinks: min(in.NumberOfLinks, r.maxNumberOfLinks),
		MaxDepth:         min(in.Depth, r.depth),
		Processing:       true,
	}

	zap.S().Debug("Creating project in db: ", prj)

	id, err := r.db.CreateProject(&prj)
	if err != nil {
		zap.S().Errorf("error while creating project: %s", err)
		return echo.ErrInternalServerError
	}

	err = r.db.SetProjectTemporaryData(id, &models.ProjectTemporaryData{
		Nodes:                 "",
		Links:                 "",
		TotalCollectorCounter: prj.MaxNumberOfLinks,
		CollectorCounterQueue: 1,
	})
	if err != nil {
		zap.S().Errorf("error while setting project temporary data: %s", err)
		return echo.ErrInternalServerError
	}

	zap.S().Debug("Adding site to parse to kafka: ", in.StartUrl)

	err = r.kafka.ProduceSite(in.StartUrl, id, 0)
	if err != nil {
		zap.S().Errorf("error while adding site to parse: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outCreate{Id: id})
}

// GetProject returns project by id
func (r *Service) GetProject(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid id"})
	}

	zap.S().Debug("Getting project from db: ", id)

	prj, err := r.db.GetProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if errors.Is(err, models.DataBaseWrongID) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	user := c.Get("user").(*models.User)

	if prj.OwnerID != user.ID {
		return c.JSON(http.StatusForbidden, errMsg{Message: "forbidden"})
	}

	if err != nil {
		zap.S().Errorf("error while getting project: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, prj)
}

// GetAllShort returns list of projects that belong to the user
func (r *Service) GetAllShort(c echo.Context) error {
	zap.S().Debug("Getting all projects from db")

	user := c.Get("user").(*models.User)

	prjs, err := r.db.GetProjectsByOwnerId(user.ID)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusOK, []models.ShortProject{})
	}

	if err != nil {
		zap.S().Errorf("error while getting projects: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, prjs)
}

type outDeleteProject struct {
	Message string `json:"message"`
}

// DeleteProject deletes project by id
func (r *Service) DeleteProject(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid id"})
	}

	zap.S().Debug("Deleting project from db: ", id)

	prj, err := r.db.GetProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if errors.Is(err, models.DataBaseWrongID) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	user := c.Get("user").(*models.User)

	if prj.OwnerID != user.ID {
		return c.JSON(http.StatusForbidden, errMsg{Message: "forbidden"})
	}

	err = r.db.DeleteProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if err != nil {
		zap.S().Errorf("error while deleting project: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outDeleteProject{Message: "ok"})
}
