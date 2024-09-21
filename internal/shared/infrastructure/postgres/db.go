package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CheckDBConnection(dataSourceName string) error {
	db, err := NewDB(dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- db.PingContext(ctx)
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
