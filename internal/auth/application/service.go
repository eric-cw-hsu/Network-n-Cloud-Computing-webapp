package application

import (
	"context"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type AuthApplicationService interface {
	Register(ctx context.Context, email, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error)
	UpdateUser(ctx context.Context, user *domain.AuthUser, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error)
}

type authApplicationService struct {
	authService domain.AuthService
	jwtService  *jwt.JWTService
	logger      logger.Logger
}

func NewAuthApplicationService(
	authService domain.AuthService,
	jwtService *jwt.JWTService,
	logger logger.Logger,
) AuthApplicationService {
	return &authApplicationService{
		authService: authService,
		jwtService:  jwtService,
		logger:      logger,
	}
}

func (s *authApplicationService) Register(ctx context.Context, email, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error) {
	// 1. check if user already exists
	exists, err := s.authService.CheckUserExists(ctx, email)
	if err != nil {
		s.logger.Error("Failed to check if user exists", err)

		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	if exists {
		s.logger.Debug("User already exists", nil)
		return &domain.AuthUser{}, apperrors.NewBadRequest("user already exists")
	}

	// 2. create user
	authUser, err := s.authService.CreateUser(ctx, email, firstName, lastName, password)
	if err != nil {
		s.logger.Error("Failed to create user", err)
		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	return authUser, nil
}

func (s *authApplicationService) UpdateUser(ctx context.Context, user *domain.AuthUser, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error) {
	// 2. update user
	err := user.Update(firstName, lastName, password)
	if err != nil {
		s.logger.Error("Failed to update user", err)
		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	err = s.authService.UpdateUser(ctx, user)
	if err != nil {
		s.logger.Error("Failed to update user", err)
		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	return user, nil
}
