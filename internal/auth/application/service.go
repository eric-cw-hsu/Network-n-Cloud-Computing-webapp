package application

import (
	"context"
	"go-template/internal/auth/domain"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/pkg/apperrors"
)

type AuthApplicationService interface {
	Register(ctx context.Context, email, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error)
	UpdateUser(ctx context.Context, user *domain.AuthUser, firstName, lastName, password string) (*domain.AuthUser, *apperrors.Error)
	VerifyAccount(ctx context.Context, token, userId string) *apperrors.Error
	ResendVerification(ctx context.Context, user *domain.AuthUser) *apperrors.Error
}

type authApplicationService struct {
	authService domain.AuthService
	logger      logger.Logger
}

func NewAuthApplicationService(
	authService domain.AuthService,
	logger logger.Logger,
) AuthApplicationService {
	return &authApplicationService{
		authService: authService,
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

	// 3. send verification email
	err = s.authService.SendVerificationEmail(authUser)
	if err != nil {
		s.logger.Error("Failed to send verification email", err)
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

func (s *authApplicationService) VerifyAccount(ctx context.Context, token, userId string) *apperrors.Error {
	// 1. verify account
	err := s.authService.VerifyVerificationEmailToken(token, userId)
	if err != nil {
		return apperrors.NewBadRequest(err.Error())
	}

	// 2. update user account status in database
	err = s.authService.VerifiedUserAccountStatus(ctx, userId)
	if err != nil {
		if err == domain.ErrUserNotFound {
			return apperrors.NewBadRequest("Invalid token")
		}

		if err == domain.ErrUserAlreadyVerified {
			return apperrors.NewBadRequest("User already verified")
		}

		s.logger.Error("Failed to update user account status", err)
		return apperrors.NewInternal()
	}

	return nil
}

func (s *authApplicationService) ResendVerification(ctx context.Context, user *domain.AuthUser) *apperrors.Error {
	// 1. check if user is already verified
	if user.Verify {
		return apperrors.NewBadRequest("User already verified")
	}

	// 2. send verification email
	err := s.authService.SendVerificationEmail(user)
	if err != nil {
		s.logger.Error("Failed to send verification email", err)
		return apperrors.NewInternal()
	}

	return nil
}
