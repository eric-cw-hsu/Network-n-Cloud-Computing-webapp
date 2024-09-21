package jwt

import (
	"go-template/pkg/apperrors"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	jwtConfig *JWTConfig
}

func NewJWTService(jwtConfig *JWTConfig) *JWTService {
	return &JWTService{
		jwtConfig: jwtConfig,
	}
}

func (service *JWTService) Authenticate(tokenString string) (*JWTUserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.jwtConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return &JWTUserInfo{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &JWTUserInfo{}, apperrors.NewAuthorization("invalid token claims")
	}

	userInfo, err := FromClaims(claims)
	if err != nil {
		return &JWTUserInfo{}, apperrors.NewAuthorization("failed to parse claims")
	}

	return userInfo, nil

}

func (service *JWTService) GenerateToken(userInfo JWTUserInfo) (string, error) {
	claims := userInfo.GenerateClaims()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(service.jwtConfig.TokenExpiration)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(service.jwtConfig.JWTSecret))
	if err != nil {
		return "", apperrors.NewAuthorization("could not sign token")
	}

	return signedToken, nil
}
