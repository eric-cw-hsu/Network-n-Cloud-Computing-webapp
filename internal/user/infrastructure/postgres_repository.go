package infrastructure

import (
	"context"
	"database/sql"
	"go-template/internal/user/domain"
	"go-template/pkg/apperrors"
)

type postgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) domain.UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, password, created_at, updated_at) VALUES ($1, $2, $3, $4)`
	res, err := r.db.ExecContext(ctx, query, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return apperrors.NewInternal()
	}

	user.ID, _ = res.LastInsertId()

	return nil
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, apperrors.NewNotFound("user")
	}
	if err != nil {
		return nil, apperrors.NewInternal()
	}
	return user, nil
}

func (r *postgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return apperrors.NewInternal()
	}
	return nil
}
