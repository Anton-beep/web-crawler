package main

import (
	"web-crauler/internal/config"
	"web-crauler/internal/services/receiver"
)

func main() {
	config.InitLogger(true)
	cfg := config.NewConfig()

	receiver.New(cfg.Receiver.Port).Start()
}
