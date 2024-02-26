package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	if db, err := sql.Open("postgres", "user=tim dbname=todos password=password sslmode=disable"); err != nil {
		return nil, err
	} else {
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return &PostgresStore{
			db: db,
		}, nil
	}
}
