package store

import (
	"database/sql"
)

type PostGresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostGresRepo {
	return &PostGresRepo{
		db: db,
	}
}
