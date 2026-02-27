package store

import (
	"context"
	"time"
)

func (r *PostgresRepo) ScheduleRetry(ctx context.Context, id string, nextRetry time.Time) error {
	query := `UPDATE notifications
	SET status = 'PENDING',
	attempts = attempts + 1,
	next_retry_at = $2,
	updated_at = now()
	WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id, nextRetry)
	return err
}
