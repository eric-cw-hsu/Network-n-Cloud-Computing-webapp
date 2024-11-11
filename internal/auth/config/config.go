package config

type AuthConfig struct {
	Auth struct {
		VerifyEmailExpirationTime int    `mapstructure:"verify_email_expiration_time"`
		VerificationEmailTopicArn string `mapstructure:"verification_email_topic_arn"`
	} `mapstructure:"auth"`
}
