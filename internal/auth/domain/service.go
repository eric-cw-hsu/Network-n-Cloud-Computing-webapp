package domain

import (
	"context"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type AuthService interface {
	CreateUser(ctx context.Context, email, firstName, lastName, password string) (*AuthUser, error)
	AuthenticateUser(ctx context.Context, email, password string) (*AuthUser, error)
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
		s.logger.Error(err)
		return nil, err
	}

	user, err = s.repository.FindUserByEmail(ctx, email)
	if err != nil {
		s.logger.Error(err)
		return nil, apperrors.NewInternal()
	}

	return user, nil
}

func (s *authService) AuthenticateUser(ctx context.Context, email, password string) (*AuthUser, error) {
	var user *AuthUser
	var err error

	if email != "" {
		user, err = s.repository.FindUserByEmail(ctx, email)
	} else {
		return nil, apperrors.NewUnprocessableEntity("email must be provided")
	}
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	if !VerifyPassword(user, password) {
		return nil, apperrors.NewAuthorization("invalid credentials")
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
