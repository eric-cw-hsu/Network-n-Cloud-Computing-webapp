package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"go-template/internal/auth/domain"
	"go-template/pkg/apperrors"

	"github.com/lib/pq"
)

type postgresAuthRepository struct {
	db *sql.DB
}

func NewPostgresAuthRepository(db *sql.DB) domain.AuthRepository {
	return &postgresAuthRepository{db: db}
}

func (r *postgresAuthRepository) Create(ctx context.Context, user *domain.AuthUser) error {
	query := `INSERT INTO users (email, username, password, created_at, updated_at, last_login_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.Username, user.PasswordHash, user.CreatedAt, user.UpdatedAt, user.LastLoginAt)
	if err != nil {
		fmt.Println(err)
		return apperrors.NewInternal()
	}

	return nil
}

func (r *postgresAuthRepository) FindUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	query := `SELECT id, email, username, password, last_login_at FROM users WHERE email = $1`
	return r.findUser(ctx, query, email)
}

func (r *postgresAuthRepository) FindUserByUsername(ctx context.Context, username string) (*domain.AuthUser, error) {
	query := `SELECT id, email, username, password, last_login_at FROM users WHERE username = $1`
	return r.findUser(ctx, query, username)
}

func (r *postgresAuthRepository) findUser(ctx context.Context, query string, arg interface{}) (*domain.AuthUser, error) {
	var user domain.AuthUser
	err := r.db.QueryRowContext(ctx, query, arg).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.LastLoginAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *postgresAuthRepository) Update(ctx context.Context, user *domain.AuthUser) error {
	query := `
			UPDATE users 
			SET email = $2, username = $3, password = $4, last_login_at = $5
			WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.LastLoginAt,
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
