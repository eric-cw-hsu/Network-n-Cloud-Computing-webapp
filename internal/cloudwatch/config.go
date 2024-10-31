package cloudwatch

type CloudWatchConfig struct {
	AWS struct {
		Region string `mapstructure:"region"`
	} `mapstructure:"aws"`
}
