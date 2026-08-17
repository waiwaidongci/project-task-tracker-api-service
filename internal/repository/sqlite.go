package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	db.SetMaxOpenConns(1)
	return &SQLiteRepository{db: db}
}

func Open(dsn string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	return db, nil
}

func closeRows(rows *sql.Rows) {
	if rows != nil {
		_ = rows.Close()
	}
}

func (r *SQLiteRepository) DB() *sql.DB {
	return r.db
}

func buildPagination(page, pageSize int) (limit, offset int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize, (page - 1) * pageSize
}

func appendWhere(base string, conditions []string, args []any) (string, []any) {
	if len(conditions) == 0 {
		return base, args
	}
	return base + " WHERE " + strings.Join(conditions, " AND "), args
}

func queryExists(ctx context.Context, db *sql.DB, query string, args ...any) (bool, error) {
	var exists int
	if err := db.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return false, err
	}
	return exists == 1, nil
}
