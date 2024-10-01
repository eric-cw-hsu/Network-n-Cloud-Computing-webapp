package http

import (
	"go-template/internal/auth/application"
	"go-template/internal/auth/interfaces/dto"
	"net/http"

	_ "go-template/docs"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService application.AuthApplicationService
}

func NewAuthHandler(authService application.AuthApplicationService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

type RegisterInput struct {
	Email     string `json:"email" example:"user@example.com"`
	FirstName string `json:"first_name" example:"John"`
	LastName  string `json:"last_name" example:"Doe"`
	Password  string `json:"password" example:"secretpassword"`
}

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body RegisterInput true "User registration details"
// @Success 201 {object} dto.UserResponse
// @Router /api/v1/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), input.Email, input.FirstName, input.LastName, input.Password)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

// @Summary Login
// @Description Login to the application
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.LoginInput true "User login details"
// @Success 200 {object} dto.LoginResponse
// @Router /api/v1/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ---- Login with JWT ----
	user, token, err := h.authService.Login(c.Request.Context(), input.Email, input.Password)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}
	// ---- [END] Login with JWT ----

	c.JSON(http.StatusOK, gin.H{
		"user":  dto.NewUserResponse(user),
		"token": token,
	})
}
