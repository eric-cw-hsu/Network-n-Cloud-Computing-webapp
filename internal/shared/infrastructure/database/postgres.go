package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

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

func (db *PostgresDatabase) AutoMigrate() error {
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return err
	}

	pwd, _ := os.Getwd()
	fileURL := fmt.Sprintf("file://%s", filepath.Join(pwd, "migrations"))

	fmt.Print("Migrating database...\n")
	migrateInstance, err := migrate.NewWithDatabaseInstance(fileURL, "postgres", driver)
	if err != nil {
		return err
	}

	// check if current database is up to date
	if err := migrateInstance.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}
