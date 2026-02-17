package worker

import (
	"context"
	"event-driven-notification-service/internal/model"
	"event-driven-notification-service/internal/store"
	"log"
	"time"
)

type Poller struct {
	repo      store.NotificationRepository
	jobQueue  chan model.Notification
	batchSize int
	interval  time.Duration
}

func (p *Poller) Start(ctx context.Context) {
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.repo.RecoverStuckJob(ctx, 5*time.Minute)
			p.fetch()
		}
	}
}

func (p *Poller) fetch() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	jobs, err := p.repo.fetchAndMarkProcessing(ctx, p.batchSize)
	if err != nil {
		log.Println("fetch error", err)
		return
	}
	for _, job := range jobs {
		p.jobQueue <- job
	}
}
