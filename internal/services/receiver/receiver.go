package receiver

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"strconv"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/repository"
)

type Service struct {
	port     int
	db       models.DataBase
	kafka    broker.SitesKafka
	depth    int    // depth of how many links to parse
	tempUUID string // temporary uuid while we don't have auth
}

func New(port int, cfgPath ...string) *Service {
	cfg := config.NewConfig(cfgPath...)
	return &Service{
		port:     port,
		db:       repository.NewDB(cfg),
		kafka:    broker.SitesKafka{},
		depth:    cfg.Receiver.Depth,
		tempUUID: cfg.Receiver.TempUUID,
	}
}

func (r *Service) Start() {
	e := echo.New()

	e.GET("/ping", Pong)
	e.POST("/project/create", r.CreateProject)
	e.GET("/project/get/:id", r.GetProject)
	e.GET("/project/getAllShort", r.GetAllShort)
	e.DELETE("/project/delete/:id", r.DeleteProject)

	err := e.Start(":" + strconv.Itoa(r.port))
	if err != nil {
		zap.S().Fatalf("error while starting server: %s", err)
	}
}
