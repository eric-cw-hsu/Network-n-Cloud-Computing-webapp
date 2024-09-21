package domain

import (
	"context"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*User, error)
	UpdateUser(ctx context.Context, id, email string) (*User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id, email string) (*User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Update(email)

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}
