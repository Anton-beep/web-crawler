package receiver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
	"web-crauler/internal/broker"
	"web-crauler/internal/models"
	"web-crauler/internal/repository"
)

type Service struct {
	port  int
	db    models.DataBase
	kafka broker.SitesKafka
}

func New(port int) *Service {
	return &Service{
		port:  port,
		db:    repository.NewDB(),
		kafka: broker.SitesKafka{},
	}
}

func (r *Service) Start() {
	e := echo.New()

	e.GET("/ping", Pong)
	e.POST("/project/create", r.CreateProject)
	e.POST("/project/get", r.GetProject)

	err := e.Start(":" + strconv.Itoa(r.port))
	if err != nil {
		zap.S().Fatalf("error while starting server: %s", err)
	}
}
