package worker

import (
	"context"
	"event-driven-notification-service/internal/model"
	"event-driven-notification-service/internal/store"
	"time"
)

type Worker struct {
	id    int
	queue <-chan model.Notification
	repo  store.NotificationRepository
}

func (w *Worker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case job := <-w.queue:
			w.process(ctx, job)
		}
	}
}

func (w *Worker) process(ctx context.Context, n model.Notification) {
	err := sendEmail(n) // stimulate notification
	if err != nil {
		w.handleFailure(ctx, n)
		return
	}
	w.repo.MarkSent(ctx, n.Id)

}

func (w *Worker) handleFailure(ctx context.Context, n model.Notification) {
	attempts := n.Attempts + 1
	if attempts >= 5 {
		w.repo.MarkFailed(ctx, n.Id)
		return
	}
	nextRetry := time.Now().Add(time.Duration(attempts*2) * time.Second)
	w.repo.ScheduleRetry(ctx, n.Id, nextRetry)
}
