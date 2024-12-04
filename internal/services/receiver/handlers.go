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

func Pong(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

type inCreate struct {
	Name     string `json:"name"`
	StartUrl string `json:"start_url"`
}

type outCreate struct {
	Id string `json:"id"`
}

func (r *Service) CreateProject(c echo.Context) error {
	var in inCreate

	if err := c.Bind(&in); err != nil {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	if in.Name == "" || in.StartUrl == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid json"})
	}

	prj := models.Project{
		Name:     in.Name,
		StartUrl: in.StartUrl,
		OwnerID:  r.tempUUID,
	}

	id, err := r.db.CreateProject(&prj)
	if err != nil {
		zap.S().Errorf("error while creating project: %s", err)
		return echo.ErrInternalServerError
	}

	err = r.kafka.AddSiteToParse(in.StartUrl, id, r.depth)
	if err != nil {
		zap.S().Errorf("error while adding site to parse: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outCreate{Id: id})
}

func (r *Service) GetProject(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid id"})
	}

	prj, err := r.db.GetProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if errors.Is(err, models.DataBaseWrongID) {
		return c.JSON(http.StatusBadRequest, errMsg{Message: err.Error()})
	}

	if err != nil {
		zap.S().Errorf("error while getting project: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, prj)
}

func (r *Service) GetAllShort(c echo.Context) error {
	prjs, err := r.db.GetProjectsByOwnerId(r.tempUUID)
	if err != nil {
		zap.S().Errorf("error while getting projects: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, prjs)
}

type outDeleteProject struct {
	Message string `json:"message"`
}

func (r *Service) DeleteProject(c echo.Context) error {
	id := c.Param("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, errMsg{Message: "invalid id"})
	}

	err := r.db.DeleteProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outDeleteProject{Message: "ok"})
}
