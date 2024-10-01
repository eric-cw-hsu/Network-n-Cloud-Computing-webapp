package application

import (
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type Authenticator interface {
	JWTAuthenticate(token string) (*domain.AuthUserInfo, *apperrors.Error)
}

type authenticatorService struct {
	jwtService *jwt.JWTService
	logger     logger.Logger
}

func NewAuthenticatorService(
	jwtService *jwt.JWTService,
	logger logger.Logger,
) Authenticator {
	return &authenticatorService{
		jwtService: jwtService,
		logger:     logger,
	}
}

func (s *authenticatorService) JWTAuthenticate(token string) (*domain.AuthUserInfo, *apperrors.Error) {
	authUserInfo, err := s.jwtService.Authenticate(token)
	if err != nil {
		s.logger.Error("Failed to authenticate token", err)
		return &domain.AuthUserInfo{}, apperrors.NewAuthorization("invalid token")
	}

	return authUserInfo, nil
}
