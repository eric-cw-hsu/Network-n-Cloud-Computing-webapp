package domain

import (
	"context"
	"fmt"
	"go-template/internal/auth/config"
	"go-template/internal/aws/sns"
	appConfig "go-template/internal/config"
	"go-template/internal/shared/infrastructure/logger"
	"time"

	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken        = fmt.Errorf("invalid token")
	ErrTokenExpired        = fmt.Errorf("token expired")
	ErrUserAlreadyVerified = fmt.Errorf("user already verified")
)

type AuthService interface {
	CreateUser(ctx context.Context, email, firstName, lastName, password string) (*AuthUser, error)
	CheckUserExists(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, user *AuthUser) error
	SendVerificationEmail(user *AuthUser) error
	VerifyVerificationEmailToken(token, userId, expiredAt string) error
	VerifiedUserAccountStatus(ctx context.Context, userId string) error
}

type authService struct {
	repository AuthRepository
	logger     logger.Logger
	authConfig *config.AuthConfig
	snsModule  sns.SNSModule
}

func NewAuthService(repo AuthRepository, logger logger.Logger, authConfig *config.AuthConfig, snsModule sns.SNSModule) AuthService {
	return &authService{
		repository: repo,
		logger:     logger,
		authConfig: authConfig,
		snsModule:  snsModule,
	}
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

/*
Send Verification Email to the user after registration.
- Generate a verification token
- Send an email with the verification token
  - The Email Service is be implemented in a micro service
  - The service will be called by Amazon Lambda
*/
func (s *authService) SendVerificationEmail(user *AuthUser) error {
	message, err := s.generateVerificationEmailMessage(user)
	if err != nil {
		return err
	}

	return s.snsModule.PublishMessage(s.authConfig.Auth.VerificationEmailTopicArn, message)
}

/*
Construct the verification email message
- The message will construct as JSON format
- And Serialize the message to string
*/
func (s *authService) generateVerificationEmailMessage(user *AuthUser) (string, error) {
	token, expiredAt, err := s.generateVerificationEmailToken(user, s.authConfig.Auth.VerifyEmailExpirationTime)
	if err != nil {
		return "", err
	}

	message := fmt.Sprintf(`{
		"to_name": "%s",
		"to_addr": "%s",
		"user_id": "%s",
		"token": "%s",
		"expiration": "%d"
	}`, user.FirstName, user.Email, user.ID, token, expiredAt.Unix())

	return message, nil
}

/*
Generate a token for verifying email
Using JWT to generate a token
*/
func (s *authService) generateVerificationEmailToken(user *AuthUser, expiredTime int) (string, time.Time, error) {
	expiredAt := time.Now().Add(time.Duration(expiredTime) * time.Second)
	claims := &jwt.StandardClaims{
		ExpiresAt: expiredAt.Unix(),
		Subject:   user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(appConfig.App.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiredAt, nil
}

func (s *authService) VerifyVerificationEmailToken(token, userId, expiredAt string) error {
	claims := &jwt.StandardClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(appConfig.App.SecretKey), nil
	})
	if err != nil {
		if (err.(*jwt.ValidationError)).Errors == jwt.ValidationErrorExpired {
			return ErrTokenExpired
		}
		return ErrInvalidToken
	}

	if claims.Subject != userId {
		return ErrInvalidToken
	}

	return nil
}

func (s *authService) VerifiedUserAccountStatus(ctx context.Context, userId string) error {
	user, err := s.repository.FindUserByID(ctx, userId)
	if err != nil {
		return err
	}

	if user.Verify {
		return ErrUserAlreadyVerified
	}

	return s.repository.VerifyAccount(ctx, user)
}
