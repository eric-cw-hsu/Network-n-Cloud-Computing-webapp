package dto

import (
	"go-template/internal/auth/domain"
)

type UserResponse struct {
	ID          int64  `json:"id"`
	Email       string `json:"email"`
	Username    string `json:"username"`
	LastLoginAt string `json:"last_login_at"`
}

func NewUserResponse(user *domain.AuthUser) *UserResponse {
	return &UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		Username:    user.Username,
		LastLoginAt: user.LastLoginAt.Format("2006-01-02 15:04:05"),
	}
}
