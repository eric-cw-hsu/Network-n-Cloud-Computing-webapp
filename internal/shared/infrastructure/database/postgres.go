package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	conn *sql.DB
}

func NewPostgresDatabase(databaseSourceString string) BaseDatabase {
	conn, err := sql.Open("postgres", databaseSourceString)
	if err != nil {
		log.Fatalf("Failed to created database instance: %v", err)
	}

	return &PostgresDatabase{conn}
}

func (db *PostgresDatabase) CheckDBConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- db.conn.PingContext(ctx)
	}()

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
		return nil
	case <-ctx.Done():
		return errors.New("db connection timeout")
	}
}

func (db *PostgresDatabase) GetConnection() *sql.DB {
	return db.conn
}

func (db *PostgresDatabase) Close() {
	db.conn.Close()
}

func (db *PostgresDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.conn.ExecContext(ctx, query, args...)
}

func (db *PostgresDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}
