package store

import "context"

func (r *PostgresRepo) MarkFailed(ctx context.Context, id string) error {
	query := `UPDATE notifications
    SET STATUS='FAILED',
    updated_at=now()
    Where id=$1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
