package config

import "time"

type Config struct {
	HTTPPort    string
	DBUrl       string
	WorkerCount int
	QueueSize   int
	RequestTTL  time.Duration
}
