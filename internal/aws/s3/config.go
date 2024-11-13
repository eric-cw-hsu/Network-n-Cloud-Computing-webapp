package s3

type S3Config struct {
	AWS struct {
		Region string `mapstructure:"region"`
		S3     struct {
			BucketName string `mapstructure:"bucket_name"`
		} `mapstructure:"s3"`
	} `mapstructure:"aws"`
}
