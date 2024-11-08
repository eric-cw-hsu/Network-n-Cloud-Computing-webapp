package application

import (
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/basic"
	"go-template/internal/auth/domain/jwt"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type Authenticator interface {
	JWTAuthenticate(token string) (*domain.AuthUserInfo, *apperrors.Error)
	BasicAuthenticate(token string) (*domain.AuthUser, *apperrors.Error)
}

type authenticatorService struct {
	jwtService   *jwt.JWTService
	basicService *basic.BasicService
	logger       logger.Logger
}

func NewAuthenticatorService(
	jwtService *jwt.JWTService,
	basicService *basic.BasicService,
	logger logger.Logger,
) Authenticator {
	return &authenticatorService{
		jwtService:   jwtService,
		basicService: basicService,
		logger:       logger,
	}
}

func (s *authenticatorService) JWTAuthenticate(token string) (*domain.AuthUserInfo, *apperrors.Error) {
	authUserInfo, err := s.jwtService.Authenticate(token)
	if err != nil {
		s.logger.Debug("Failed to authenticate token", err)
		return &domain.AuthUserInfo{}, apperrors.NewAuthorization("invalid token")
	}

	return authUserInfo, nil
}

func (s *authenticatorService) BasicAuthenticate(token string) (*domain.AuthUser, *apperrors.Error) {
	user, err := s.basicService.Authenticate(token)
	if err != nil {
		if err == basic.ErrInvalidToken {
			s.logger.Debug("Failed to authenticate user", err)
		} else {
			s.logger.Error("Failed to authenticate user", err)
		}
		return &domain.AuthUser{}, apperrors.NewAuthorization("invalid credentials")
	}

	return user, nil
}
