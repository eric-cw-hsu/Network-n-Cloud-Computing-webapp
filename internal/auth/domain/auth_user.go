package domain

import (
	"errors"
	"go-template/pkg/apperrors"
	"time"
)

type AuthUser struct {
	ID           int64
	Email        string
	Username     string
	PasswordHash string
	Role         string
	LastLoginAt  time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func NewAuthUser(email string, username string, password string) (*AuthUser, error) {
	if email == "" || username == "" {
		return nil, apperrors.NewBadRequest("email or username cannot be empty for both")
	}

	if password == "" {
		return nil, apperrors.NewBadRequest("password cannot be empty")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	return &AuthUser{
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastLoginAt:  time.Now(),
	}, nil
}

func (u *AuthUser) UpdateLastLogin() {
	u.LastLoginAt = time.Now()
	u.UpdatedAt = time.Now()
}
