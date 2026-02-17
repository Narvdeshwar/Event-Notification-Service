package main

import (
	"event-driven-notification-service/internal/config"
	"runtime/metrics"
)

func main() {
	cfg := config.Load()
	metrics.Register()
}
