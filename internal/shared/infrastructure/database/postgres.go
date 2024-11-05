package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-template/internal/cloudwatch"
	"go-template/internal/config"
	"go-template/internal/utils"
	"log"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

type PostgresDatabase struct {
	conn             *sql.DB
	cloudWatchModule cloudwatch.CloudWatchModule
}

func NewPostgresDatabase(databaseSourceString string, cloudWatchModule cloudwatch.CloudWatchModule) BaseDatabase {
	conn, err := sql.Open("postgres", databaseSourceString)
	if err != nil {
		log.Fatalf("Failed to created database instance: %v", err)
	}

	initPostgresParameters(conn)

	return &PostgresDatabase{conn, cloudWatchModule}
}

func initPostgresParameters(conn *sql.DB) {
	// TODO: Set these values in the config file
	conn.SetMaxOpenConns(config.App.Database.MaxOpenConns)
	conn.SetMaxIdleConns(config.App.Database.MaxIdleConns)
	conn.SetConnMaxLifetime(time.Minute * 5)
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
	start := time.Now()

	result, err := db.conn.ExecContext(ctx, query, args...)

	defer db.logLatencyMetric(ctx, query, float64(time.Since(start).Milliseconds()))

	return result, err
}

func (db *PostgresDatabase) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()

	row := db.conn.QueryRowContext(ctx, query, args...)

	defer db.logLatencyMetric(ctx, query, float64(time.Since(start).Milliseconds()))

	return row
}

func (db *PostgresDatabase) AutoMigrate() error {
	driver, err := postgres.WithInstance(db.conn, &postgres.Config{})
	if err != nil {
		return err
	}

	rootPath, _ := utils.GetProjectRootPath()
	fileURL := fmt.Sprintf("file://%s", filepath.Join(rootPath, "migrations"))

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

func (db *PostgresDatabase) logLatencyMetric(ctx context.Context, query string, latency float64) {
	db.cloudWatchModule.PublishMetric(
		config.App.Name+"/Postgres",
		query,
		latency,
		types.StandardUnitMilliseconds,
	)
}
