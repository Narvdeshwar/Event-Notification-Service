package store

import (
	"context"
	"database/sql"
	"event-driven-notification-service/internal/model"
)

type PostGresRepo struct {
	db *sql.DB
}

func(r *PostGresRepo) Insert(ctx context.Context,n *model.Notification) error {
	query:=`Insert into notifications(id,type,recipients,payload,status,attempts,created_at,updated_at)
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8) on conflict (id) do nothing`;
	_,err:=r.db.ExecContext(ctx,query,n.Id,n.Type,n.Recipient,n.Payload,n.Status,n.CreatedAt,n.UpdatedAt)
	return err
}
func (r *PostGresRepo) FetchAndMarkProcessing(
	ctx context.Context,
	limit int,
) ([]model.Notification, error) {

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query := `
	SELECT id, type, recipient, payload, status,
	       attempts, next_retry_at, created_at, updated_at
	FROM notifications
	WHERE status = 'PENDING'
	AND (next_retry_at IS NULL OR next_retry_at <= now())
	ORDER BY created_at
	LIMIT $1
	FOR UPDATE SKIP LOCKED
	`

	rows, err := tx.QueryContext(ctx, query, limit)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Notification

	for rows.Next() {
		var n model.Notification
		err := rows.Scan(
			&n.Id,
			&n.Type,
			&n.Recipient,
			&n.Payload,
			&n.Status,
			&n.Attempts,
			&n.NextRetry,
			&n.CreatedAt,
			&n.UpdatedAt,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		jobs = append(jobs, n)
	}

	// 🔥 Immediately mark them as PROCESSING
	for _, job := range jobs {
		_, err := tx.ExecContext(
			ctx,
			`UPDATE notifications 
			 SET status = 'PROCESSING',
			     updated_at = now()
			 WHERE id = $1`,
			job.Id,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return jobs, nil
}

