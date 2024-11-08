package database

import (
	"context"
	"database/sql"
	"errors"
)

type BaseDatabase interface {
	CheckDBConnection() error
	GetConnection() *sql.DB
	Close()
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	AutoMigrate() error
}

var ErrNoRows = sql.ErrNoRows
var ErrDatabaseError = errors.New("database error")
