package http

import (
	authDomain "go-template/internal/auth/domain"
	"go-template/internal/s3"
	"go-template/internal/user/application"
	"go-template/internal/user/domain"
	"go-template/internal/user/interfaces/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userApplicationService application.UserApplicationService
	s3Module               s3.S3Module
}

func NewUserHandler(userApplicationService application.UserApplicationService, s3Module s3.S3Module) *UserHandler {
	return &UserHandler{
		userApplicationService: userApplicationService,
		s3Module:               s3Module,
	}
}

func (h *UserHandler) UploadProfilePic(c *gin.Context) {
	authUser, _ := c.Get("user")
	user := domain.NewUser(authUser.(*authDomain.AuthUser).ID)

	// multipart/form-data
	profilePicFile, err := c.FormFile("profilePic")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "profilePic is required",
		})
		return
	}

	// Validate file extension
	if !h.userApplicationService.ValidateProfilePicExtension(profilePicFile.Filename) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Invalid file extension",
		})
		return
	}

	// check if profile pic already exists
	profilePic, apperr := h.userApplicationService.GetProfilePic(c, user)
	if apperr == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Profile pic already exists",
		})
		return
	}

	// Save the file
	profilePic, apperr = h.userApplicationService.UploadProfilePic(c, user, profilePicFile)
	if apperr != nil {
		c.JSON(apperr.Status(), gin.H{
			"error": apperr.Message,
		})
		return
	}

	c.JSON(http.StatusOK, dto.NewPicResponse(user, profilePic, h.s3Module.GetBucketName()))
}

func (h *UserHandler) GetProfilePic(c *gin.Context) {
	authUser, _ := c.Get("user")
	user := domain.NewUser(authUser.(*authDomain.AuthUser).ID)

	profilePic, err := h.userApplicationService.GetProfilePic(c, user)
	if err != nil {
		c.JSON(err.Status(), gin.H{
			"error": err.Message,
		})
		return
	}

	c.JSON(http.StatusOK, dto.NewPicResponse(user, profilePic, h.s3Module.GetBucketName()))
}

func (h *UserHandler) DeleteProfilePic(c *gin.Context) {
	authUser, _ := c.Get("user")
	user := domain.NewUser(authUser.(*authDomain.AuthUser).ID)

	apperr := h.userApplicationService.DeleteProfilePic(c, user)
	if apperr != nil {
		c.JSON(apperr.Status(), gin.H{
			"error": apperr.Message,
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
