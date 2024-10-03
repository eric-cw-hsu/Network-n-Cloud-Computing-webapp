package domain

import (
	"encoding/json"
	"errors"

	"github.com/golang-jwt/jwt"
)

type AuthUserInfo struct {
	ID    int64
	Email string
}

func (authUserInfo AuthUserInfo) GenerateClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"id":    authUserInfo.ID,
		"email": authUserInfo.Email,
	}
}

func (authUserInfo AuthUserInfo) MarshalBinary() ([]byte, error) {
	return json.Marshal(authUserInfo)
}

func (authUserInfo *AuthUserInfo) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, authUserInfo)
}

func FromClaims(claims jwt.MapClaims) (*AuthUserInfo, error) {
	id, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("missing key in claims")
	}

	for _, key := range []string{"email", "username", "role"} {
		if _, ok := claims[key]; !ok {
			return nil, errors.New("missing key in claims")
		}
	}

	return &AuthUserInfo{
		ID:    int64(id),
		Email: claims["email"].(string),
	}, nil
}

func NewAuthUserInfo(id int64, email string) AuthUserInfo {
	return AuthUserInfo{
		ID:    id,
		Email: email,
	}
}
