package receiver

import (
	"context"
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
	kafka            *broker.Kafka
	depth            int // depth of how many links to parse
	maxNumberOfLinks int
	tempUUID         string // temporary uuid while we don't have auth
	secretSignature  []byte
	debug            bool
	echo             *echo.Echo
}

func New(port int, cfgPath ...string) *Service {
	cfg := config.NewConfig(cfgPath...)
	return &Service{
		port:             port,
		db:               repository.NewDB(cfg),
		kafka:            broker.New(cfg, true, false, "sites"),
		depth:            cfg.Receiver.Depth,
		maxNumberOfLinks: cfg.Receiver.MaxNumberOfLinks,
		tempUUID:         cfg.Receiver.TempUUID,
		secretSignature:  []byte(cfg.Receiver.SecretSignature),
		debug:            cfg.Debug,
	}
}

// Start starts the receiver server
func (r *Service) Start() {
	r.echo = echo.New()

	r.echo.Debug = r.debug

	r.echo.Use(middleware.CORS())
	r.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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

	notAuthorized := r.echo.Group("/api")

	notAuthorized.GET("/ping", Pong)

	notAuthorized.POST("/user/register", r.Register)
	notAuthorized.POST("/user/login", r.Login)

	authorized := r.echo.Group("/api")
	authorized.Use(r.AuthMiddleware)

	authorized.POST("/project/create", r.CreateProject)
	authorized.GET("/project/get/:id", r.GetProject)
	authorized.GET("/project/getAllShort", r.GetAllShort)
	authorized.DELETE("/project/delete/:id", r.DeleteProject)

	authorized.GET("/user/get", r.GetUser)
	authorized.PUT("/user/update", r.UpdateUser)

	go func() {
		err := r.echo.Start(":" + strconv.Itoa(r.port))
		if err != nil {
			zap.S().Fatalf("error while starting server: %s", err)
		}
	}()
}

func (r *Service) Stop() {
	zap.S().Info("gracefully stopping server")
	if err := r.echo.Shutdown(context.Background()); err != nil {
		zap.S().Errorf("error while stopping server: %s", err)
	}
}
