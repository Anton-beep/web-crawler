/*
 * Receiver Service
 *
 * Description:
 * The Receiver service acts as an intermediary that accepts user requests (indirectly),
 * processes them, and redirects the tasks to the web crawlers (collectors).
 * This service is designed to handle high loads and supports horizontal scaling.
 *
 * Usage:
 * - Configure the port for the Receiver in the configuration file.
 * - Start the service to begin accepting requests and delegating tasks.
 *
 * This file is the entry point of the application and sets up the Receiver service.
 */

package main

import (
	"os"
	"os/signal"
	"syscall"
	"web-crawler/internal/config"
	"web-crawler/internal/services/receiver"
)

func main() {
	cfg := config.NewConfig()
	config.InitLogger(cfg.Debug)

	r := receiver.New(cfg.Receiver.Port)

	r.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	r.Stop()
}
