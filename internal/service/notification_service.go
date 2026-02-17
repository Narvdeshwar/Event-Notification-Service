package service

import (
	"context"
	"event-driven-notification-service/internal/model"
	"event-driven-notification-service/internal/store"
	"time"

	"github.com/google/uuid"
)

type NotificationService struct {
	repo store.NotificationRepository
}

func New(repo store.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (s *NotificationService) Enqueue(
	ctx context.Context,
	typ string,
	recipient string,
	payload []byte,
) error {

	n := &model.Notification{
		Id:        uuid.NewString(),
		Type:      typ,
		Recipient: recipient,
		Payload:   payload,
		Status:    model.StatusPending,
		Attempts:  0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repo.Insert(ctx, n)
}
