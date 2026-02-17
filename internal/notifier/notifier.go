package notifier

import (
	"context"
	"event-driven-notification-service/internal/model"
)

type Notifier interface {
	Send(ctx context.Context, n model.Notification) error
}