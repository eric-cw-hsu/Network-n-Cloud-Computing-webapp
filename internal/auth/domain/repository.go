package domain

import (
	"context"
)

type AuthRepository interface {
	Create(ctx context.Context, user *AuthUser) error
	FindUserByEmail(ctx context.Context, email string) (*AuthUser, error)
	FindUserByUsername(ctx context.Context, username string) (*AuthUser, error)
	Update(ctx context.Context, user *AuthUser) error
}
