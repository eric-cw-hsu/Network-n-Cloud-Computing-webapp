package domain

import (
	"errors"
	"time"

	"github.com/samborkent/uuidv7"
)

type AuthUser struct {
	ID           string
	Email        string
	FirstName    string
	LastName     string
	PasswordHash string

	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrDuplicateEntry    = errors.New("duplicate entry")
	ErrUserAlreadyExists = errors.New("user already exists")
)

func NewAuthUser(
	email string, firstName string,
	lastName string, password string,
) (*AuthUser, error) {
	if email == "" {
		return &AuthUser{}, errors.New("email cannot be empty")
	}

	if firstName == "" || lastName == "" {
		return &AuthUser{}, errors.New("first name and last name cannot be empty")
	}

	if password == "" {
		return &AuthUser{}, errors.New("password cannot be empty")
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return &AuthUser{}, err
	}

	return &AuthUser{
		ID:           uuidv7.New().String(),
		Email:        email,
		FirstName:    firstName,
		LastName:     lastName,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}, nil
}

func (u *AuthUser) UpdateLastLogin() {
	u.UpdatedAt = time.Now()
}

func (u *AuthUser) Update(email, firstName, lastName, password string) error {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	u.Email = email
	u.FirstName = firstName
	u.LastName = lastName
	u.PasswordHash = hashedPassword
	u.UpdatedAt = time.Now()

	return nil
}
