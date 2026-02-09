package store

import (
	"context"
	"event-driven-notification-service/internal/model"
)

type Notification_Repo interface {
	Insert(ctx context.Context,n *model.Notification) error
}
