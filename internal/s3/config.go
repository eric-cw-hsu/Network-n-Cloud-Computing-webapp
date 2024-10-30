package s3

type S3Config struct {
	AWS struct {
		Region     string `mapstructure:"region"`
		BucketName string `mapstructure:"bucket_name"`
	} `mapstructure:"aws"`
}
