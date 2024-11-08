package domain

import (
	"errors"
	"go-template/internal/s3"
	"mime/multipart"
)

var ErrInvalidProfilePicContent = errors.New("invalid profile pic content")

type UserService interface {
	ParseProfilePic(profilePic *multipart.FileHeader) ([]byte, error)
	UploadProfilePic(userId, filename string, fileBytes []byte) (*ProfilePic, error)
	DeleteProfilePic(key string) error
	GetProfilePic(key string) ([]byte, error)
}

type userService struct {
	s3Module   s3.S3Module
	repository UserRepository
}

func NewUserService(s3Module s3.S3Module) UserService {
	return &userService{
		s3Module: s3Module,
	}
}

func (s *userService) ParseProfilePic(profilePic *multipart.FileHeader) ([]byte, error) {
	file, err := profilePic.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileBytes := make([]byte, profilePic.Size)
	if _, err = file.Read(fileBytes); err != nil {
		return nil, ErrInvalidProfilePicContent
	}

	return fileBytes, nil
}

func (s *userService) UploadProfilePic(userId, filename string, fileBytes []byte) (*ProfilePic, error) {
	uniqueKey := userId + "/" + filename
	uploadResult, err := s.s3Module.UploadFile(uniqueKey, fileBytes)
	if err != nil {
		return nil, err
	}

	return NewProfilePic(
		filename,
		uploadResult.Location,
		*uploadResult.Key,
		*uploadResult.ETag,
		uploadResult.ServerSideEncryption,
		*uploadResult.SSEKMSKeyId,
	), nil
}

func (s *userService) DeleteProfilePic(key string) error {
	return s.s3Module.DeleteFile(key)
}

func (s *userService) GetProfilePic(key string) ([]byte, error) {
	return s.s3Module.GetFile(key)
}
