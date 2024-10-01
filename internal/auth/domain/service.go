package domain

import (
	"context"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type AuthService interface {
	CreateUser(ctx context.Context, email, username, password string) (*AuthUser, error)
	AuthenticateUser(ctx context.Context, email, username, password string) (*AuthUser, error)
	CheckUserExists(ctx context.Context, email, username string) (bool, error)
}

type authService struct {
	repository AuthRepository
	logger     logger.Logger
}

func NewAuthService(repo AuthRepository, logger logger.Logger) AuthService {
	return &authService{repository: repo, logger: logger}
}

func (s *authService) CreateUser(ctx context.Context, email, username, password string) (*AuthUser, error) {
	user, err := NewAuthUser(email, username, password)
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

func (s *authService) AuthenticateUser(ctx context.Context, email, username, password string) (*AuthUser, error) {
	var user *AuthUser
	var err error

	if email != "" {
		user, err = s.repository.FindUserByEmail(ctx, email)
	} else if username != "" {
		user, err = s.repository.FindUserByUsername(ctx, username)
	} else {
		return nil, apperrors.NewUnprocessableEntity("either email or username must be provided")
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
func (s *authService) CheckUserExists(ctx context.Context, email, username string) (bool, error) {
	user, err := s.repository.FindUserByEmail(ctx, email)
	if err != nil && err != ErrUserNotFound {
		return false, err
	}

	if user != nil {
		return true, nil
	}

	user, err = s.repository.FindUserByUsername(ctx, username)
	if err != nil && err != ErrUserNotFound {
		return false, err
	}

	if user != nil {
		return true, nil
	}

	return false, nil
}
