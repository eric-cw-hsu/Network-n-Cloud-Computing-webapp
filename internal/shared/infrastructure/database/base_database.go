package database

import (
	"context"
	"database/sql"
)

type BaseDatabase interface {
	CheckDBConnection() error
	GetConnection() *sql.DB
	Close()
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}
