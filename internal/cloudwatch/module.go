package cloudwatch

import (
	"context"
	"go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudWatchModule interface {
	PublishMetric(namespace, metricName string, value float64, unit types.StandardUnit)
}

type module struct {
	client           *cloudwatch.Client
	cloudWatchConfig *CloudWatchConfig
	logger           logger.Logger
}

func NewModule(logger logger.Logger) CloudWatchModule {
	cloudWatchConfig := loadConfig()
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(cloudWatchConfig.AWS.Region),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := cloudwatch.NewFromConfig(cfg)

	return &module{
		client:           client,
		cloudWatchConfig: cloudWatchConfig,
		logger:           logger,
	}
}

func loadConfig() *CloudWatchConfig {
	// load cloudwatch config with viper
	var cloudwatchConfig *CloudWatchConfig
	if err := config.Load(&cloudwatchConfig); err != nil {
		log.Fatalf("Failed to load cloudwatch config: %v", err)
	}

	return cloudwatchConfig
}

func (m *module) PublishMetric(namespace, metricName string, value float64, unit types.StandardUnit) {
	_, err := m.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
		Namespace: aws.String(namespace),
		MetricData: []types.MetricDatum{
			{
				MetricName: aws.String(metricName),
				Value:      aws.Float64(value),
				Unit:       unit,
			},
		},
	})
	if err != nil {
		m.logger.Error("Failed to publish metric: ", err)
		return
	}
}
