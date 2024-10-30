package s3

import (
	"bytes"
	"context"
	"go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Module interface {
	GetFile(key string) ([]byte, error)
	UploadFile(key string, file []byte) error
	DeleteFile(key string) error
	GetBucketName() string
}

type module struct {
	s3Config *S3Config
	client   *s3.Client
}

func NewModule(logger logger.Logger) S3Module {
	s3Config := loadConfig()

	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(s3Config.AWS.Region),
	)

	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	return &module{
		client:   client,
		s3Config: s3Config,
	}
}

func loadConfig() *S3Config {
	// load s3 config with viper
	var s3Config *S3Config
	if err := config.Load(&s3Config); err != nil {
		log.Fatalf("Failed to load s3 config: %v", err)
	}

	return s3Config
}

func (m *module) GetFile(key string) ([]byte, error) {
	downloader := manager.NewDownloader(m.client)
	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return buffer.Bytes()[:numBytes], nil
}

func (m *module) UploadFile(key string, file []byte) error {
	uploader := manager.NewUploader(m.client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})

	return err
}

func (m *module) DeleteFile(key string) error {
	_, err := m.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
	})
	return err
}

func (m *module) GetBucketName() string {
	return m.s3Config.AWS.BucketName
}
