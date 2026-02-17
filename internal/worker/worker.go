package worker

import (
	"context"
	"event-driven-notification-service/internal/model"
	"event-driven-notification-service/internal/notifier"
	"event-driven-notification-service/internal/store"
	"time"
)

type Worker struct {
	id       int
	queue    <-chan model.Notification
	repo     store.NotificationRepository
	notifier notifier.Notifier
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
	err := w.notifier.Send(ctx, n)
	if err != nil {
		w.handleFailure(ctx, n, err)
		return
	}
	w.repo.MarkSent(ctx, n.Id)

}

func (w *Worker) handleFailure(
	ctx context.Context,
	n model.Notification,
	err error,
) {

	maxRetries := 5
	attempts := n.Attempts + 1

	if attempts >= maxRetries {
		// w.repo.MarkFailed(ctx, n.Id)
		w.repo.MoveToDeadLetter(ctx, n, err.Error())
		return
	}

	backoff := time.Duration(attempts*2) * time.Second
	nextRetry := time.Now().Add(backoff)

	w.repo.ScheduleRetry(ctx, n.Id, nextRetry)
}
