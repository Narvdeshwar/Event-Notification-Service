package store

import (
	"context"
	"event-driven-notification-service/internal/model"
)

func (r *PostgresRepo) Insert(ctx context.Context, n *model.Notification) error {
	query := `INSERT INTO notifications (id, type, recipients, payload, status, attempts, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8) ON CONFLICT (id) DO NOTHING`

	_, err := r.db.ExecContext(
		ctx,
		query,
		n.Id,
		n.Type,
		n.Recipient,
		n.Payload,
		n.Status,
		n.Attempts,
		n.CreatedAt,
		n.UpdatedAt,
	)

	return err
}
