package cloudwatch

import (
	"context"
	"go-template/internal/shared/config"
	"go-template/internal/shared/infrastructure/logger"
	"log"
	"sync"
	"time"

	appConfig "go-template/internal/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

type CloudWatchModule interface {
	PublishMetric(namespace, metricName string, value float64, unit types.StandardUnit)
	Shutdown()
}

type module struct {
	client           *cloudwatch.Client
	cloudWatchConfig *CloudWatchConfig
	logger           logger.Logger
	metricDataBuffer map[string][]types.MetricDatum
	mu               sync.Mutex
	pushInterval     time.Duration
	bufferSize       int
	shutdownChan     chan struct{}
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

	mod := &module{
		client:           client,
		cloudWatchConfig: cloudWatchConfig,
		logger:           logger,
		metricDataBuffer: make(map[string][]types.MetricDatum),
		pushInterval:     time.Duration(cloudWatchConfig.AWS.CloudWatch.PushInterval) * time.Second,
		bufferSize:       cloudWatchConfig.AWS.CloudWatch.BufferSize,
		shutdownChan:     make(chan struct{}),
	}

	go mod.startAutoFlush()

	return mod
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
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.metricDataBuffer[namespace]; !exists {
		m.metricDataBuffer[namespace] = []types.MetricDatum{}
	}

	// append metric data to buffer
	m.metricDataBuffer[namespace] = append(m.metricDataBuffer[namespace], types.MetricDatum{
		MetricName: aws.String(metricName),
		Value:      aws.Float64(value),
		Unit:       unit,
		Timestamp:  aws.Time(time.Now()),
	})

	if len(m.metricDataBuffer[namespace]) >= m.bufferSize {
		m.flush(namespace)
	}
}

func (m *module) startAutoFlush() {
	ticker := time.NewTicker(m.pushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.Lock()
			for namespace := range m.metricDataBuffer {
				m.flush(namespace)
			}
			m.mu.Unlock()
		case <-m.shutdownChan:
			m.mu.Lock()
			for namespace := range m.metricDataBuffer {
				m.flush(namespace)
			}
			m.mu.Unlock()
			return
		}
	}
}

func (m *module) flush(namespace string) {
	if metrics, exists := m.metricDataBuffer[namespace]; exists && len(metrics) > 0 {

		if appConfig.App.Environment != "development" {
			_, err := m.client.PutMetricData(context.TODO(), &cloudwatch.PutMetricDataInput{
				Namespace:  aws.String(namespace),
				MetricData: metrics,
			})
			if err != nil {
				m.logger.Error("Failed to publish metrics: ", err)
			}
		}

		m.metricDataBuffer[namespace] = m.metricDataBuffer[namespace][:0]
	}
}

func (m *module) Shutdown() {
	close(m.shutdownChan)

	time.Sleep(2 * time.Second)
}
