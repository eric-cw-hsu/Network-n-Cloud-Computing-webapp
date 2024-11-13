package application

import (
	"go-template/internal/auth/domain"
	"go-template/internal/auth/domain/basic"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type Authenticator interface {
	BasicAuthenticate(token string) (*domain.AuthUser, *apperrors.Error)
}

type authenticatorService struct {
	basicService *basic.BasicService
	logger       logger.Logger
}

func NewAuthenticatorService(
	basicService *basic.BasicService,
	logger logger.Logger,
) Authenticator {
	return &authenticatorService{
		basicService: basicService,
		logger:       logger,
	}
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
