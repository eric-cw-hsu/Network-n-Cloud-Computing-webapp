package jwt

type JWTConfig struct {
	JWTSecret       string
	TokenExpiration int
}

func NewJWTConfig(secret string, expiration int) *JWTConfig {
	return &JWTConfig{
		JWTSecret:       secret,
		TokenExpiration: expiration,
	}
}
