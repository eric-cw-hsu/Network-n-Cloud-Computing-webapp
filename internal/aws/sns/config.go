package sns

type SNSConfig struct {
	AWS struct {
		Region string `mapstructure:"region"`
	} `mapstructure:"aws"`
}
