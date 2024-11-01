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
	query := `INSERT INTO user_pic(user_id, filename, uploaded_at, url, s3_key, etag, encryption, encryption_key) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.ExecContext(ctx, query, user.ID, profilePic.Filename, profilePic.UploadedAt, profilePic.Url, profilePic.S3Key, profilePic.ETag, profilePic.Encryption, profilePic.EncryptionKey)
	if err != nil {
		return err
	}

	return nil
}

func (r *postgresUserRepository) GetProfilePic(ctx context.Context, user *domain.User) (*domain.ProfilePic, error) {
	query := `SELECT filename, uploaded_at, url, s3_key, etag, encryption, encryption_key FROM user_pic WHERE user_id = $1`

	profilePic := domain.ProfilePic{}
	err := r.db.QueryRowContext(ctx, query, user.ID).Scan(&profilePic.Filename, &profilePic.UploadedAt, &profilePic.Url, &profilePic.S3Key, &profilePic.ETag, &profilePic.Encryption, &profilePic.EncryptionKey)
	if err != nil {
		return nil, err
	}

	return &profilePic, nil
}

func (r *postgresUserRepository) DeleteProfilePic(ctx context.Context, user *domain.User) error {
	query := `DELETE FROM user_pic WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, user.ID)
	if err != nil {
		return err
	}

	return nil
}
