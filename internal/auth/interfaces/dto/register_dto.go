package dto

type RegisterInput struct {
	Email     string `json:"email" example:"user@example.com" binding:"required,email"`
	FirstName string `json:"first_name" example:"John" binding:"required"`
	LastName  string `json:"last_name" example:"Doe" binding:"required"`
	Password  string `json:"password" example:"secretpassword" binding:"required"`
}
