package s3

import (
	"bytes"
	"context"
	"go-template/internal/cloudwatch"
	appConfig "go-template/internal/config"
	"go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Module interface {
	GetFile(key string) ([]byte, error)
	UploadFile(key string, file []byte) error
	DeleteFile(key string) error
	GetBucketName() string
}

type module struct {
	s3Config         *S3Config
	client           *s3.Client
	cloudWatchModule cloudwatch.CloudWatchModule
}

func NewModule(logger logger.Logger, cloudWatchModule cloudwatch.CloudWatchModule) S3Module {
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
		client:           client,
		s3Config:         s3Config,
		cloudWatchModule: cloudWatchModule,
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
	startTime := time.Now()

	downloader := manager.NewDownloader(m.client)
	buffer := manager.NewWriteAtBuffer([]byte{})

	numBytes, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
	})

	defer m.logLatencyMetric("get_file", float64(time.Since(startTime).Milliseconds()))

	if err != nil {
		return nil, err
	}

	return buffer.Bytes()[:numBytes], nil
}

func (m *module) UploadFile(key string, file []byte) error {
	startTime := time.Now()

	uploader := manager.NewUploader(m.client)
	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(file),
	})

	defer m.logLatencyMetric("upload_file", float64(time.Since(startTime).Milliseconds()))

	return err
}

func (m *module) DeleteFile(key string) error {
	startTime := time.Now()

	_, err := m.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(m.s3Config.AWS.BucketName),
		Key:    aws.String(key),
	})

	defer m.logLatencyMetric("delete_file", float64(time.Since(startTime).Milliseconds()))
	return err
}

func (m *module) GetBucketName() string {
	return m.s3Config.AWS.BucketName
}

func (m *module) logLatencyMetric(action string, latency float64) {
	m.cloudWatchModule.PublishMetric(
		appConfig.App.Name+"/S3",
		action+"_latency",
		latency,
		types.StandardUnitMilliseconds,
	)
}
