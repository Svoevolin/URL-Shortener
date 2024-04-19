package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/Svoevolin/url-shortener/internal/database"
	"github.com/lib/pq"
)

const (
	UniqueViolationError = pq.ErrorCode("23505")
)

type UrlDB struct {
	db *sql.DB
}

func NewUrlDB(db *sql.DB) *UrlDB {
	return &UrlDB{db: db}
}

type Url struct {
	ID    int    `db:"id"`
	Alias string `db:"alias"`
	Url   string `db:"url"`
}

func (db *UrlDB) SaveURL(urlToSave string, alias string) (int64, error) {
	const op = "database.postgres.models.SaveURL"

	stmt, err := db.db.Prepare("INSERT INTO url(url, alias) VALUES ($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRow(urlToSave, alias).Scan(&id)
	if err != nil {
		if IsErrorCode(err, UniqueViolationError) {
			return 0, fmt.Errorf("%s: %w", op, database.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (db *UrlDB) GetURL(alias string) (string, error) {
	const op = "database.postgres.models.GetURL"

	stmt, err := db.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", database.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: execute statement %w", op, err)
	}

	return url, nil
}

func IsErrorCode(err error, errorCode pq.ErrorCode) bool {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		return pgErr.Code == errorCode
	}
	return false
}
