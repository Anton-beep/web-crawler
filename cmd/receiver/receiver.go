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
	"web-crawler/internal/config"
	"web-crawler/internal/services/receiver"
)

func main() {
	cfg := config.NewConfig()
	config.InitLogger(cfg.Debug)

	receiver.New(cfg.Receiver.Port).Start()
}
