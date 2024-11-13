package basic

import (
	"context"
	"encoding/base64"
	"fmt"
	"go-template/internal/auth/domain"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type MockAuthRepository struct {
	mock.Mock
}

func (m *MockAuthRepository) FindUserByEmail(ctx context.Context, email string) (*domain.AuthUser, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.AuthUser), args.Error(1)
}
func (m *MockAuthRepository) Create(ctx context.Context, user *domain.AuthUser) error { return nil }
func (m *MockAuthRepository) FindUserByUsername(ctx context.Context, username string) (*domain.AuthUser, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*domain.AuthUser), args.Error(1)
}
func (m *MockAuthRepository) FindUserByID(ctx context.Context, id string) (*domain.AuthUser, error) {
	return nil, nil
}
func (m *MockAuthRepository) Update(ctx context.Context, user *domain.AuthUser) error { return nil }
func (m *MockAuthRepository) VerifyAccount(ctx context.Context, user *domain.AuthUser) error {
	return nil
}

func TestBasicService(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockAuthRepository := new(MockAuthRepository)
	baseService := NewBasicService(mockAuthRepository)
	t.Run("Test splitToken", func(t *testing.T) {
		token := "username:password"
		email, password, err := baseService.splitToken(token)
		if err != nil {
			t.Errorf("Error when split token: %v", err)
		}

		if email != "username" {
			t.Errorf("Email should be username, got %s", email)
		}

		if password != "password" {
			t.Errorf("Password should be password, got %s", password)
		}

		// test invalid token
		token = "username"
		email, password, err = baseService.splitToken(token)
		if err == nil {
			t.Errorf("Error should not be nil")
		}

		if email != "" {
			t.Errorf("Email should be empty, got %s", email)
		}

		if password != "" {
			t.Errorf("Password should be empty, got %s", password)
		}
	})

	t.Run("Test Authenticate", func(t *testing.T) {
		password := "iampassword"
		mockUser, _ := domain.NewAuthUser(
			"test@example.com",
			"First",
			"Last",
			password,
		)

		mockAuthRepository.On("FindUserByEmail", mock.Anything, mock.Anything).Return(mockUser, nil)

		base64Token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", mockUser.Email, password)))

		_, err := baseService.Authenticate(base64Token)
		if err != nil {
			t.Errorf("Error should be nil, got %v", err)
		}

		// test invalid password
		_, err = baseService.Authenticate(base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", mockUser.Email, "invalidpassword"))))
		if err == nil {
			t.Errorf("Error should not be nil")
		}
		if err != ErrInvalidToken {
			t.Errorf("Error should be invalid credentials, got %v", err)
		}
	})
}
