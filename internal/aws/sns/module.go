package sns

import (
	"context"
	"fmt"
	"go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

type SNSModule interface {
	PublishMessage(topicArn string, message string) error
}

type module struct {
	snsConfig *SNSConfig
	client    *sns.Client
	logger    logger.Logger
}

func NewModule(logger logger.Logger) SNSModule {
	snsConfig := loadConfig()

	ctx := context.Background()
	sdkConfig, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		fmt.Println(err)
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}

	snsClient := sns.NewFromConfig(sdkConfig)

	return &module{
		snsConfig: snsConfig,
		client:    snsClient,
		logger:    logger,
	}
}

func loadConfig() *SNSConfig {
	// load s3 config with viper
	var snsConfig *SNSConfig
	if err := config.Load(&snsConfig); err != nil {
		log.Fatalf("Failed to load sns config: %v", err)
	}

	return snsConfig
}

func (m *module) PublishMessage(topicArn string, message string) error {
	input := &sns.PublishInput{
		Message:  &message,
		TopicArn: &topicArn,
	}

	_, err := m.client.Publish(context.Background(), input)
	if err != nil {
		m.logger.Error("Failed to publish message to SNS topic ", err)
		return err
	}

	return nil
}
