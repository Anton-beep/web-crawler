package receiver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"strconv"
	"web-crawler/internal/broker"
	"web-crawler/internal/config"
	"web-crawler/internal/models"
	"web-crawler/internal/repository"
)

type Service struct {
	port             int
	db               models.DataBase
	kafka            *broker.SitesKafka
	depth            int // depth of how many links to parse
	maxNumberOfLinks int
	tempUUID         string // temporary uuid while we don't have auth
	secretSignature  []byte
}

func New(port int, cfgPath ...string) *Service {
	cfg := config.NewConfig(cfgPath...)
	return &Service{
		port:             port,
		db:               repository.NewDB(cfg),
		kafka:            broker.New(cfg, true, false),
		depth:            cfg.Receiver.Depth,
		maxNumberOfLinks: cfg.Receiver.MaxNumberOfLinks,
		tempUUID:         cfg.Receiver.TempUUID,
		secretSignature:  []byte(cfg.Receiver.SecretSignature),
	}
}

// Start starts the receiver server
func (r *Service) Start() {
	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			zap.L().Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))

	notAuthorized := e.Group("/api")

	notAuthorized.GET("/ping", Pong)

	notAuthorized.POST("/user/register", r.Register)

	authorized := e.Group("/api")
	authorized.Use(r.AuthMiddleware)

	authorized.POST("/create", r.CreateProject)
	authorized.GET("/get/:id", r.GetProject)
	authorized.GET("/getAllShort", r.GetAllShort)
	authorized.DELETE("/delete/:id", r.DeleteProject)

	authorized.POST("/user/login/:id", r.Login)

	err := e.Start(":" + strconv.Itoa(r.port))
	if err != nil {
		zap.S().Fatalf("error while starting server: %s", err)
	}
}
