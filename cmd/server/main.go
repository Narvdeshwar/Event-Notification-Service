package main

import (
	"notification-service/internal/api"
	"notification-service/internal/config"
	"notification-service/internal/service"
	"notification-service/internal/store"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db := store.NewPostgres(cfg.DBUrl)
	repo := store.NewNotificationRepo(db)

	svc := service.New(repo)
	handler := api.New(svc)

	r := gin.New()
	api.RegisterRoutes(r, handler)
	r.Run(":" + cfg.HTTPPort)

}
