package notifier

import (
	"context"
	"errors"
	"event-driven-notification-service/internal/model"
	"math/rand"
	"time"
)

type EmailNotifier struct{}

func NewEmailNotifier() *EmailNotifier {
	return &EmailNotifier{}
}

func (e *EmailNotifier) Send(ctx context.Context,n model.Notification) error {
	time.Sleep(500*time.Millisecond)
    if rand.Intn(10)<3{
        return errors.New("Email Provider Failed")
    }
    return nil
}