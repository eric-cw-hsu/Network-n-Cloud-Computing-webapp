package infrastructure

import (
	"context"
	"go-template/internal/shared/infrastructure/database"
	"go-template/internal/user/domain"
)

type postgresUserRepository struct {
	db database.BaseDatabase
}

func NewPostgresUserRepository(db database.BaseDatabase) domain.UserRepository {
	return &postgresUserRepository{
		db: db,
	}
}

func (r *postgresUserRepository) SaveProfilePic(ctx context.Context, user *domain.User, profilePic *domain.ProfilePic) error {
	query := `UPDATE users SET pic_filename = $1, pic_uploaded_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, profilePic.Filename, profilePic.UploadedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresUserRepository) GetProfilePic(ctx context.Context, user *domain.User) (*domain.ProfilePic, error) {
	query := `SELECT pic_filename, pic_uploaded_at FROM users WHERE id = $1`

	profilePic := domain.ProfilePic{}
	err := r.db.QueryRowContext(ctx, query, user.ID).Scan(&profilePic.Filename, &profilePic.UploadedAt)
	if err != nil {
		return nil, err
	}

	return &profilePic, nil
}

func (r *postgresUserRepository) DeleteProfilePic(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET pic_filename = '', pic_uploaded_at = NULL WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}

	return nil
}
