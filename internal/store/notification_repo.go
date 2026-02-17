package store

import (
	"context"
	"event-driven-notification-service/internal/model"
	"time"
)

type NotificationRepository interface {
	Insert(ctx context.Context, n *model.Notification) error
	FetchAndMarkProcessing(ctx context.Context, limit int) ([]model.Notification, error)
	MarkSent(ctx context.Context, id string) error
	ScheduleRetry(ctx context.Context, id string, nextRetry time.Time) error
	MarkFailed(ctx context.Context,id string) error
	RecoverStuckJob(ctx context.Context,timeout time.Duration) error
}
