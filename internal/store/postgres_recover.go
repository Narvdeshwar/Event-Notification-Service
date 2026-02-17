package store

import (
	"context"
	"time"
)

func (r *PostGresRepo) RecoverStuckJobs(ctx context.Context, timeout time.Duration) error {
	query := `UPDATE notifications
	SET status='PENDING',
	updated_at=now()
	WHERE status='PROCESSING'
	AND updated_at<now()-$1::interval`
	_, err := r.db.ExecContext(ctx, query, timeout.String())
	return err
}
