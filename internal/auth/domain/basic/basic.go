package basic

import (
	"context"
	"encoding/base64"
	"errors"
	"go-template/internal/auth/domain"
	"strings"
	"time"
)

type BasicService struct {
	authRepository domain.AuthRepository
}

func NewBasicService(authRepository domain.AuthRepository) *BasicService {
	return &BasicService{
		authRepository: authRepository,
	}
}

func (bs *BasicService) Authenticate(token string) (*domain.AuthUser, error) {
	// 0. decode token from base64 format
	decodedByte64Token, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return &domain.AuthUser{}, err
	}
	token = string(decodedByte64Token)

	// 1. split token into username and password with ":"
	email, password, err := bs.splitToken(token)
	if err != nil {
		return &domain.AuthUser{}, err
	}

	// 2. check if username and password are valid
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	user, err := bs.authRepository.FindUserByEmail(ctx, email)
	if err != nil {
		return &domain.AuthUser{}, err
	}

	if !domain.VerifyPassword(user, password) {
		return &domain.AuthUser{}, errors.New("invalid credentials")
	}

	return user, nil
}

func (bs *BasicService) splitToken(token string) (email, password string, err error) {
	slices := strings.Split(token, ":")
	if len(slices) != 2 {
		return "", "", errors.New("invalid token")
	}

	return slices[0], slices[1], nil
}
