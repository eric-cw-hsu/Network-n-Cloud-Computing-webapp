package domain

import "context"

type UserRepository interface {
	SaveProfilePic(ctx context.Context, user *User, profilePic *ProfilePic) error
	GetProfilePic(ctx context.Context, user *User) (*ProfilePic, error)
	DeleteProfilePic(ctx context.Context, user *User) error
}
