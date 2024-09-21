package config

type AuthConfig struct {
	Auth struct {
		JWTSecret       string `mapstructure:"jwt_secret"`
		TokenExpiration int    `mapstructure:"token_expiration"`
	} `mapstructure:"auth"`
}
