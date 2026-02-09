package service

import "event-driven-notification-service/internal/store"

type NotificationService struct {
	repo store.NotificationRepository
}
