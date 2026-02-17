package store

import (
	"context"
	"time"
)

func (r *PostGresRepo) ScheduleRetry(ctx context.Context, id string, nextRetry time.Time) error {
	query := `UPDATE notifications
	SET STATUS='PENDING,
	attempts=attepts+1,
	next_retry_at=$2,
	updated_at=now()
	where id=$1'
	`
	_, err := r.db.ExecContext(ctx, query, id, nextRetry)
	return err
}
