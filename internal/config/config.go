package config

import (
	"os"
	"time"
)

type Config struct {
	HTTPPort    string
	DBUrl       string
	WorkerCount int
	QueueSize   int
	RequestTTL  time.Duration
}

func Load() Config {
	return Config{
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		DBUrl: getEnv(
			"DB_URL",
			"postgres://postgres:postgres@postgres:5432/notifications?sslmode=disable",
		),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
