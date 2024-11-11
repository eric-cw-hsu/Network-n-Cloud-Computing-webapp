package middleware

import (
	"fmt"
	"go-template/internal/aws/cloudwatch"
	"go-template/internal/config"
	"go-template/internal/shared/infrastructure/logger"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/gin-gonic/gin"
)

type RequestLoggerMiddleware struct {
	logger           logger.Logger
	cloudWatchModule cloudwatch.CloudWatchModule
}

func NewRequestLoggerMiddleware(logger logger.Logger, cloudWatchModule cloudwatch.CloudWatchModule) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{
		logger:           logger,
		cloudWatchModule: cloudWatchModule,
	}
}

// Handler logs the request
// It logs request method, request path, request ip, latency, and response status code
func (m *RequestLoggerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		m.logger.Info(
			fmt.Sprintf(
				`{"method": "%s", "path": "%s", "client_ip": "%s", "latency": "%s", "status": %d, "errors": "%s"}`,
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				latency,
				c.Writer.Status(),
				c.Errors.String(),
			),
		)

		apiName := fmt.Sprintf("[%s]%s", c.Request.Method, c.Request.URL.Path)

		m.cloudWatchModule.PublishMetric(
			config.App.Name+"/API",
			apiName+"_latency",
			float64(latency.Milliseconds()),
			types.StandardUnitMilliseconds,
		)

		m.cloudWatchModule.PublishMetric(
			config.App.Name+"/API",
			apiName+"_count",
			1,
			types.StandardUnitCount,
		)
	}
}
