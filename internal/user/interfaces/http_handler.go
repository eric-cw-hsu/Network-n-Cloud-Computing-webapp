package interfaces

import (
	"go-template/internal/user/application"
	"go-template/internal/user/interfaces/dto"
	"go-template/pkg/apperrors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService application.UserApplicationService
}

func NewUserHandler(userService application.UserApplicationService) *UserHandler {
	return &UserHandler{userService: userService}
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "User ID"
// @Success 200 {object} dto.UserResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	user, err := h.userService.GetUser(c.Request.Context(), id)
	if err != nil {
		appErr := err.(*apperrors.Error)
		c.JSON(appErr.Status(), gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateUserEmail(c.Request.Context(), id, input.Email)
	if err != nil {
		appErr := err.(*apperrors.Error)
		c.JSON(appErr.Status(), gin.H{"error": appErr.Message})
		return
	}

	c.JSON(http.StatusOK, dto.NewUserResponse(user))
}
