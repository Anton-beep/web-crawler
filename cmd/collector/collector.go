package main

import (
	"web-crawler/internal/config"
	"web-crawler/internal/services/collector"
)

func main() {
	cfg := config.NewConfig()
	config.InitLogger(cfg.Debug)

	server := collector.NewServer(cfg)

	server.Start()
}
