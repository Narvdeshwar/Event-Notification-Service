package store

import (
	"context"
	"notification-service/internal/model"
)

type NotificationRepository interface {
	Insert(ctx Context.Context, n *model.Notification) error
	FetchPending(ctx context.Context, int limit) ([]model.Notification, error)
	UpdateStatus(ctx context.Context, id string, status model.Status) error
	IncrementAttempts(ctx context.Context, id string) error
}
