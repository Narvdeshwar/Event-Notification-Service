package store

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Connect establishes a connection to the database with retries.
func Connect(dbURL string) (*sql.DB, error) {
	var db *sql.DB
	var err error
	maxRetries := 10

	for i := 0; i < maxRetries; i++ {
		db, err = sql.Open("postgres", dbURL)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			return db, nil
		}

		log.Printf("Failed to connect to db (retry %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
