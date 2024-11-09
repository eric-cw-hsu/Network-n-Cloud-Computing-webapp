package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"go-template/internal/auth/domain"
	"go-template/internal/shared/infrastructure/database"

	"github.com/lib/pq"
)

type postgresAuthRepository struct {
	db database.BaseDatabase
}

func NewPostgresAuthRepository(db database.BaseDatabase) domain.AuthRepository {
	return &postgresAuthRepository{db: db}
}

func (r *postgresAuthRepository) Create(ctx context.Context, user *domain.AuthUser) error {
	query := `INSERT INTO users (id, email, first_name, last_name, password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.FirstName, user.LastName, user.PasswordHash, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return database.ErrDatabaseError
	}

	return nil
}

func (r *postgresAuthRepository) FindUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE email = $1`
	return r.findUser(ctx, query, email)
}

func (r *postgresAuthRepository) FindUserByUsername(ctx context.Context, username string) (*domain.AuthUser, error) {
	query := `SELECT id, email, first_name, last_name, password, created_at, updated_at FROM users WHERE username = $1`
	return r.findUser(ctx, query, username)
}

func (r *postgresAuthRepository) findUser(ctx context.Context, query string, arg interface{}) (*domain.AuthUser, error) {
	var user domain.AuthUser
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, database.ErrDatabaseError
	}

	return &user, nil
}

func (r *postgresAuthRepository) Update(ctx context.Context, user *domain.AuthUser) error {
	query := `
			UPDATE users 
			SET first_name = $2, last_name = $3,
			password = $4, updated_at = $5
			WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.FirstName,
		user.LastName,
		user.PasswordHash,
		user.UpdatedAt,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return domain.ErrDuplicateEntry
			}
		}
		return err
	}
	return nil
}
