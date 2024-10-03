package dto

import (
	"go-template/internal/auth/domain"
)

type UserResponse struct {
	ID             int64  `json:"id"`
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	AccountCreated string `json:"account_created"`
	AccountUpdated string `json:"account_updated"`
}

func NewUserResponse(user *domain.AuthUser) *UserResponse {
	return &UserResponse{
		ID:             user.ID,
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		AccountCreated: user.CreatedAt.Format("2006-01-02 15:04:05"),
		AccountUpdated: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
