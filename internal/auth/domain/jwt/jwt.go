package jwt

import (
	"go-template/internal/auth/domain"
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

func (service *JWTService) Authenticate(tokenString string) (*domain.AuthUserInfo, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.jwtConfig.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return &domain.AuthUserInfo{}, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return &domain.AuthUserInfo{}, err
	}

	userInfo, err := domain.FromClaims(claims)
	if err != nil {
		return &domain.AuthUserInfo{}, err
	}

	return userInfo, nil
}

func (service *JWTService) GenerateToken(userInfo domain.AuthUserInfo) (string, error) {
	claims := userInfo.GenerateClaims()
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(service.jwtConfig.TokenExpiration)).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(service.jwtConfig.JWTSecret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
