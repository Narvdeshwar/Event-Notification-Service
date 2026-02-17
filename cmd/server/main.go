package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/gin-gonic/gin"

	"event-driven-notification-service/internal/api"
	"event-driven-notification-service/internal/config"
	"event-driven-notification-service/internal/metrics"
	"event-driven-notification-service/internal/service"
	"event-driven-notification-service/internal/store"
)

func main() {

	// 1️⃣ Load config
	cfg := config.Load()

	// 2️⃣ Register metrics
	metrics.Register()

	// 3️⃣ Connect to database
	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("failed to open db:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to db:", err)
	}

	log.Println("Database connected successfully")

	// 4️⃣ Create repository
	repo := store.NewPostgresRepo(db)

	// 5️⃣ Create service (inject repo)
	svc := service.New(repo)

	// 6️⃣ Create handler (inject service)
	handler := api.New(svc)

	// 7️⃣ Setup router
	router := gin.Default()
	api.RegisterRoutes(router, handler)

	// 8️⃣ Start HTTP server
	log.Println("Server running on port", cfg.HTTPPort)
	if err := router.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatal("server failed:", err)
	}
}
