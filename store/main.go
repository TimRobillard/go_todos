package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	user := os.Getenv("PG_USER")
	dbname := os.Getenv("PG_DBNAME")
	pw := os.Getenv("PG_PASSDWORD")
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", user, dbname, pw)
	if db, err := sql.Open("postgres", connStr); err != nil {
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
