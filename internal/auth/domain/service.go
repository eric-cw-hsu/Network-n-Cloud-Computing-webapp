package domain

import (
	"context"
	"go-template/internal/shared/infrastructure/logger"
)

type AuthService interface {
	CreateUser(ctx context.Context, email, firstName, lastName, password string) (*AuthUser, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, user *AuthUser) error
}

type authService struct {
	repository AuthRepository
	logger     logger.Logger
}

func NewAuthService(repo AuthRepository, logger logger.Logger) AuthService {
	return &authService{repository: repo, logger: logger}
}

func (s *authService) CreateUser(ctx context.Context, email, firstName, lastName, password string) (*AuthUser, error) {
	user, err := NewAuthUser(email, firstName, lastName, password)
	if err != nil {
		return nil, err
	}

	err = s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user, err = s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// CheckUserExists checks if a user with the given email or username already exists in the database
func (s *authService) CheckUserExists(ctx context.Context, email string) (bool, error) {
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil && err != ErrUserNotFound {
		return false, err
	}

	if user != nil {
		return true, nil
	}

	return false, nil
}

func (s *authService) UpdateUser(ctx context.Context, user *AuthUser) error {
	return s.repository.Update(ctx, user)
}
