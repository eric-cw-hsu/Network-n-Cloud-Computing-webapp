package application

import (
	"context"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/shared/infrastructure/logger"
)

type AuthApplicationService interface {
	Register(ctx context.Context, email, username, password string) (*domain.AuthUser, error)
	Login(ctx context.Context, email, username, password string) (*domain.AuthUser, string, error)
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

func (s *authApplicationService) Register(ctx context.Context, email, username, password string) (*domain.AuthUser, error) {
	// 1. check if user already exists
	exists, err := s.authService.CheckUserExists(ctx, email, username)
	if err != nil {
		return nil, err
	}

	if exists {
		return nil, domain.ErrUserAlreadyExists
	}

	// 2. create user
	return s.authService.CreateUser(ctx, email, username, password)
}

func (s *authApplicationService) Login(ctx context.Context, email, username, password string) (*domain.AuthUser, string, error) {
	// 1. check username, email, password in the database
	user, err := s.authService.AuthenticateUser(ctx, email, username, password)
	if err != nil {
		return nil, "", err
	}

	// 2. generate jwt token
	jwtUserInfo := jwt.NewJWTUserInfo(user.ID, user.Email, user.Username, user.Role)

	token, err := s.jwtService.GenerateToken(jwtUserInfo)
	if err != nil {
		return nil, "", err
	}

	user.UpdateLastLogin()

	return user, token, nil
}
