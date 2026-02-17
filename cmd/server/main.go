package main

import (
	"event-driven-notification-service/internal/config"
	"event-driven-notification-service/internal/metrics"
)

func main() {

	cfg := config.Load()

	metrics.Register()

	_ = cfg
}
