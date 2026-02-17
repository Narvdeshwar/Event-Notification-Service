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
