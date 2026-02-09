package store

import (
	"context"
	"event-driven-notification-service/internal/model"
)

type NotificationRepository interface {
	Insert(ctx context.Context,n *model.Notification) error
}
