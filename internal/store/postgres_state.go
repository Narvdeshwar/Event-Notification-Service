package store

import "context"

func (r *PostgresRepo) MarkSent(ctx context.Context, id string) error {
	query := `UPDATE notifications
	set STATUS='SENT', updated_at=now()
	Where id=$1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
