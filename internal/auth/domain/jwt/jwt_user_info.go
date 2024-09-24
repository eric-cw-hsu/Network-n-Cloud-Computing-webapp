package jwt

import (
	"go-template/pkg/apperrors"

	"github.com/golang-jwt/jwt"
)

type JWTUserInfo struct {
	ID       int64
	Email    string
	Username string
	Role     string
}

func (jwtUserInfo JWTUserInfo) GenerateClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"id":       jwtUserInfo.ID,
		"email":    jwtUserInfo.Email,
		"username": jwtUserInfo.Username,
		"role":     jwtUserInfo.Role,
	}
}

func FromClaims(claims jwt.MapClaims) (*JWTUserInfo, error) {
	id, ok := claims["id"].(float64)
	if !ok {
		return nil, apperrors.NewInvalidClaims("id")
	}

	for _, key := range []string{"email", "username", "role"} {
		if _, ok := claims[key]; !ok {
			return nil, apperrors.NewInvalidClaims(key)
		}
	}

	return &JWTUserInfo{
		ID:       int64(id),
		Email:    claims["email"].(string),
		Username: claims["username"].(string),
		Role:     claims["role"].(string),
	}, nil
}

func NewJWTUserInfo(
	id int64, email string, username string, role string,
) JWTUserInfo {
	return JWTUserInfo{
		ID:       id,
		Email:    email,
		Username: username,
		Role:     role,
	}
}
