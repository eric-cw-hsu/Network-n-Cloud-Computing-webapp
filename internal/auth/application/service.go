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
	Login(ctx context.Context, email, password string) (*domain.AuthUser, string, *apperrors.Error)
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
	// TODO: validate input

	// 1. check if user already exists
	exists, err := s.authService.CheckUserExists(ctx, email)
	if err != nil {
		s.logger.Error("Failed to check if user exists", err)
		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	if exists {
		s.logger.Error("User already exists", nil)
		return &domain.AuthUser{}, apperrors.NewConflict("user already exists")
	}

	// 2. create user
	authUser, err := s.authService.CreateUser(ctx, email, firstName, lastName, password)
	if err != nil {
		s.logger.Error("Failed to create user", err)
		return &domain.AuthUser{}, apperrors.NewInternal()
	}

	return authUser, nil
}

func (s *authApplicationService) Login(ctx context.Context, email, password string) (*domain.AuthUser, string, *apperrors.Error) {
	return s.loginWithJWT(ctx, email, password)
}

func (s *authApplicationService) loginWithJWT(ctx context.Context, email, password string) (*domain.AuthUser, string, *apperrors.Error) {
	// 1. check username, email, password in the database
	user, err := s.authService.AuthenticateUser(ctx, email, password)
	if err != nil {
		s.logger.Error("Failed to authenticate user", err)
		return nil, "", apperrors.NewAuthorization("invalid credentials")
	}

	authUserInfo := domain.NewAuthUserInfo(user.ID, user.Email)

	// 2. generate jwt token
	token, err := s.jwtService.GenerateToken(authUserInfo)
	if err != nil {
		s.logger.Error("Failed to generate token", err)
		return nil, "", apperrors.NewInternal()
	}

	// 3. update last login
	user.UpdateLastLogin()

	return user, token, nil
}
