package main

import (
	"go.uber.org/zap"
	"web-crawler/internal/config"
	"web-crawler/internal/repository"
)

func main() {
	cfg := config.NewConfig()
	config.InitLogger(true)
	zap.S().Debug(cfg)
	_ = repository.NewDB(cfg)
}
