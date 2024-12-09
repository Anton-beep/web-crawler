/*
 * Collector service
 *
 * Description:
 * The Collector service is responsible for fetching the
 * HTML content of the target web pages. It is designed to
 * extract titles, links and text from the HTML content.
 *
 * Usage:
 * - Configure the port for the Collector in the configuration file.
 * - Start the service to begin fetching the HTML content.
 */

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
