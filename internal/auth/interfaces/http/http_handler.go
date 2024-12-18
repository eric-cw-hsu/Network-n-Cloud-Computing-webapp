package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-template/internal/auth/application"
	"go-template/internal/auth/domain"
	"go-template/internal/auth/interfaces/dto"
	"io"
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

// @Summary Register a new user
// @Description Create a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param input body dto.RegisterInput true "User registration details"
// @Success 201 {object} dto.UserResponse
// @Router /v1/user [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var input dto.RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), input.Email, input.FirstName, input.LastName, input.Password)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusCreated, dto.NewUserResponse(user))
}

// @Summary Get user profile
// @Description Get user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.UserResponse
// @Router /v1/user [get]
func (h *AuthHandler) GetUser(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, dto.NewUserResponse(user.(*domain.AuthUser)))
}

// @Summary Update user profile
// @Description Update user profile
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body dto.UpdateUserInput true "User update details"
// @Success 200 {object} dto.UserResponse
// @Router /v1/user [put]
func (h *AuthHandler) UpdateUser(c *gin.Context) {

	rawBody, parseErr := c.GetRawData()
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": parseErr.Error()})
		return
	}

	// check if the input contains invalid data without dto.UpdateUserInput
	if err := checkFieldsIsValid(rawBody, []string{"first_name", "last_name", "password"}); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(rawBody))

	var input dto.UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, _ := c.Get("user")
	_, err := h.authService.UpdateUser(c.Request.Context(), user.(*domain.AuthUser), input.FirstName, input.LastName, input.Password)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) VerifyAccount(c *gin.Context) {
	token := c.Query("token")
	userId := c.Query("user_id")

	if token == "" || userId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token, and user_id are required"})
		return
	}

	err := h.authService.VerifyAccount(c.Request.Context(), token, userId)
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *AuthHandler) ResendVerification(c *gin.Context) {
	user, _ := c.Get("user")
	err := h.authService.ResendVerification(c.Request.Context(), user.(*domain.AuthUser))
	if err != nil {
		c.JSON(err.Status(), gin.H{"error": err.Message})
		return
	}

	c.Status(http.StatusNoContent)
}

func checkFieldsIsValid(rawBody []byte, expectedFields []string) error {
	var data map[string]interface{}
	if err := json.Unmarshal(rawBody, &data); err != nil {
		return err
	}

	// Check for extra fields
	for key := range data {
		if !contains(expectedFields, key) {
			return fmt.Errorf("Extra field: %s", key)
		}
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
