package domain

import (
	"go-template/pkg/apperrors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewUser(email, password string) (*User, error) {
	if email == "" {
		return nil, apperrors.NewBadRequest("email cannot be empty")
	}
	if password == "" {
		return nil, apperrors.NewBadRequest("password cannot be empty")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.NewInternal()
	}

	return &User{
		Email:     email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (u *User) Update(email string) {
	u.Email = email
	u.UpdatedAt = time.Now()
}
