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
	debug            bool
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
		debug:            cfg.Debug,
	}
}

// Start starts the receiver server
func (r *Service) Start() {
	e := echo.New()

	e.Debug = r.debug

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
	notAuthorized.POST("/user/login", r.Login)

	authorized := e.Group("/api")
	authorized.Use(r.AuthMiddleware)

	authorized.POST("/project/create", r.CreateProject)
	authorized.GET("/project/get/:id", r.GetProject)
	authorized.GET("/project/getAllShort", r.GetAllShort)
	authorized.DELETE("/project/delete/:id", r.DeleteProject)

	err := e.Start(":" + strconv.Itoa(r.port))
	if err != nil {
		zap.S().Fatalf("error while starting server: %s", err)
	}
}
