package application

import (
	"context"
	"database/sql"
	"go-template/internal/shared/infrastructure/logger"
	"go-template/internal/user/domain"
	"go-template/pkg/apperrors"
	"mime/multipart"
	"strings"
)

type UserApplicationService interface {
	UploadProfilePic(ctx context.Context, user *domain.User, profilePicFile *multipart.FileHeader) (*domain.ProfilePic, *apperrors.Error)
	DeleteProfilePic(ctx context.Context, user *domain.User) *apperrors.Error
	GetProfilePic(ctx context.Context, user *domain.User) (*domain.ProfilePic, *apperrors.Error)
	ValidateProfilePicExtension(filename string) bool
}

type userApplicationService struct {
	logger         logger.Logger
	userService    domain.UserService
	userRepository domain.UserRepository
}

func NewUserApplicationService(logger logger.Logger, userService domain.UserService, userRepository domain.UserRepository) UserApplicationService {
	return &userApplicationService{
		logger:         logger,
		userService:    userService,
		userRepository: userRepository,
	}
}

func (s *userApplicationService) UploadProfilePic(ctx context.Context, user *domain.User, profilePicFile *multipart.FileHeader) (*domain.ProfilePic, *apperrors.Error) {

	fileBytes, err := s.userService.ParseProfilePic(profilePicFile)
	if err != nil {
		s.logger.Error("Failed to parse profile pic", err)
		return nil, apperrors.NewInternal()
	}

	profilePic, err := s.userService.UploadProfilePic(user.ID, profilePicFile.Filename, fileBytes)
	if err != nil {
		s.logger.Error("Failed to upload profile pic", err)
		return nil, apperrors.NewInternal()
	}

	err = s.userRepository.SaveProfilePic(ctx, user, profilePic)
	if err != nil {
		s.logger.Error("Failed to save profile pic", err)
		return nil, apperrors.NewInternal()
	}

	return profilePic, nil
}

func (s *userApplicationService) DeleteProfilePic(ctx context.Context, user *domain.User) *apperrors.Error {
	// Get the profile pic
	profilePic, err := s.userRepository.GetProfilePic(ctx, user)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("profile pic not found", err)
			return apperrors.NewNotFound("profile pic not found")
		}

		s.logger.Error("Failed to get profile pic", err)
		return apperrors.NewInternal()
	}

	// Delete the profile pic from S3
	err = s.userService.DeleteProfilePic(user.ID + "/" + profilePic.Filename)
	if err != nil {
		s.logger.Error("Failed to delete profile pic", err)
		return apperrors.NewInternal()
	}

	// Delete the profile pic from the database
	err = s.userRepository.DeleteProfilePic(ctx, user)
	if err != nil {
		s.logger.Error("Failed to delete profile pic from database", err)
		return apperrors.NewInternal()
	}

	return nil
}

func (s *userApplicationService) ValidateProfilePicExtension(filename string) bool {
	allowedExtensions := []string{"jpg", "jpeg", "png"}

	extension := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension {
			return true
		}
	}

	return false
}

func (s *userApplicationService) GetProfilePic(ctx context.Context, user *domain.User) (*domain.ProfilePic, *apperrors.Error) {
	profilePic, err := s.userRepository.GetProfilePic(ctx, user)

	if err != nil {
		if err == sql.ErrNoRows {
			s.logger.Error("profile pic not found", err)
			return nil, apperrors.NewNotFound("profile pic not found")
		}

		s.logger.Error("Failed to get profile pic", err)
		return nil, apperrors.NewInternal()
	}

	return profilePic, nil
}
