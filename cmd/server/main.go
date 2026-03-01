package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"event-driven-notification-service/internal/api"
	"event-driven-notification-service/internal/config"
	"event-driven-notification-service/internal/metrics"
	"event-driven-notification-service/internal/service"
	"event-driven-notification-service/internal/store"
	"event-driven-notification-service/internal/migrations"
)

func main() {
	// Load config
	cfg := config.Load()

	// Register metrics
	metrics.Register()

	// Connect to database
	db, err := store.Connect(cfg.DBUrl)
	if err != nil {
		log.Fatal("Could not connect to database after retries:", err)
	}
	defer db.Close()


	// run migrations after db is connected
	migrations.Run(db)

	// Create repository
	repo := store.NewPostgresRepo(db)

	// Create service (inject repo)
	svc := service.New(repo)

	// Create handler (inject service)
	handler := api.New(svc)

	// Setup router
	router := gin.Default()
	api.RegisterRoutes(router, handler)

	// Start HTTP server
	log.Println("Server running on port", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("server failed:", err)
	}
}
