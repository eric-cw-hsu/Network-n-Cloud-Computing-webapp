package application

import (
	"context"
	"go-template/internal/user/domain"
)

type UserApplicationService interface {
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateUserEmail(ctx context.Context, id, email string) (*domain.User, error)
}

type userApplicationService struct {
	userService domain.UserService
}

func NewUserApplicationService(userService domain.UserService) UserApplicationService {
	return &userApplicationService{userService: userService}
}

func (s *userApplicationService) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return s.userService.GetUserByID(ctx, id)
}

func (s *userApplicationService) UpdateUserEmail(ctx context.Context, id, email string) (*domain.User, error) {
	return s.userService.UpdateUser(ctx, id, email)
}
