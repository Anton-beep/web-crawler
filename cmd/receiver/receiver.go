package main

import (
	"web-crawler/internal/config"
	"web-crawler/internal/services/receiver"
)

func main() {
	config.InitLogger(true)
	cfg := config.NewConfig()

	receiver.New(cfg.Receiver.Port).Start()
}
