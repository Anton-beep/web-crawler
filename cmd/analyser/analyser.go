package main

import (
	"web-crawler/internal/config"
	"web-crawler/internal/services/analyser"
)

func main() {
	cfg := config.NewConfig()
	config.InitLogger(cfg.Debug)

	server := analyser.NewServer(cfg)

	server.Start()
}
