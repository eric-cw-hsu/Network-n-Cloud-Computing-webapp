package domain

import (
	"context"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {

	// GetByID retrieves a user by their ID
	GetByID(ctx context.Context, id string) (*User, error)

	// Update updates an existing user in the repository
	Update(ctx context.Context, user *User) error
}
