package cloudwatch

type CloudWatchConfig struct {
	AWS struct {
		Region string `mapstructure:"region"`

		CloudWatch struct {
			PushInterval int `mapstructure:"push_interval"`
			BufferSize   int `mapstructure:"buffer_size"`
		} `mapstructure:"cloudwatch"`
	} `mapstructure:"aws"`
}
