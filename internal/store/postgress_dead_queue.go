package store

import (
	"context"
	"event-driven-notification-service/internal/model"
	"github.com/google/uuid"
)

func (r *PostGresRepo) MoveToDeadLetter(
	ctx context.Context,
	n model.Notification,
	errMsg string,
) error {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	insertQuery := `
	INSERT INTO dead_letter_notifications
	(id, notification_id, payload, error, attempts, failed_at)
	VALUES ($1, $2, $3, $4, $5, now())
	`

	_, err = tx.ExecContext(
		ctx,
		insertQuery,
		uuid.NewString(),
		n.Id,
		n.Payload,
		errMsg,
		n.Attempts,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE notifications SET status = 'FAILED' WHERE id = $1`,
		n.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
