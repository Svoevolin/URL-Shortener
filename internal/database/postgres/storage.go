package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func New(dsn string) (*sql.DB, error) {
	const op = "database.postgres.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}
