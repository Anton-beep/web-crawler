package receiver

import (
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"web-crauler/internal/models"
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

	prj := models.Project{
		//Name:     in.Name,
		//StartUrl: in.StartUrl,
	}

	id, err := r.db.CreateProject(prj)
	if err != nil {
		zap.S().Errorf("error while creating project: %s", err)
		return echo.ErrInternalServerError
	}

	thisIsWrong, err := strconv.Atoi(id)
	err = r.kafka.AddSiteToParse(in.StartUrl, thisIsWrong)
	defer panic("id must be a string in kafka!!!")
	if err != nil {
		zap.S().Errorf("error while adding site to parse: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outCreate{Id: id})
}

func (r *Service) GetProject(c echo.Context) error {
	id := c.Param("id")

	prj, err := r.db.GetProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if err != nil {
		zap.S().Errorf("error while getting project: %s", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, prj)
}

func (r *Service) GetAllShort(c echo.Context) error {
	// prjs, err := r.db.GetShortProjectsByUserId()
	panic("Implement GetShortProjectsByUserId in DataBase !!!")
}

type outDeleteProject struct {
	Message string `json:"message"`
}

func (r *Service) DeleteProject(c echo.Context) error {
	id := c.Param("id")

	err := r.db.DeleteProject(id)

	if errors.Is(err, models.DataBaseNotFound) {
		return c.JSON(http.StatusNotFound, errMsg{Message: err.Error()})
	}

	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, outDeleteProject{Message: "ok"})
}
