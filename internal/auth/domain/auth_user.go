package domain

import (
	"errors"
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
		return &AuthUser{}, errors.New("email and username cannot be empty")
	}

	if password == "" {
		return &AuthUser{}, errors.New("password cannot be empty")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return &AuthUser{}, err
	}

	return &AuthUser{
		Email:        email,
		Username:     username,
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
