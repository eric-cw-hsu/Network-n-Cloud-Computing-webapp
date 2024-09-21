package middleware

import (
	"fmt"
	"go-template/internal/shared/infrastructure/logger"
	"time"

	"github.com/gin-gonic/gin"
)

type RequestLoggerMiddleware struct {
	logger logger.Logger
}

func NewRequestLoggerMiddleware(logger logger.Logger) *RequestLoggerMiddleware {
	return &RequestLoggerMiddleware{logger: logger}
}

// Handler logs the request
// It logs request method, request path, request ip, latency, and response status code
func (m *RequestLoggerMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {

		startTime := time.Now()

		c.Next()

		latency := time.Since(startTime)

		m.logger.Info(fmt.Sprintf("%s %s %s %s %d %s",
			c.Request.Method,
			c.Request.URL.Path,
			c.ClientIP(),
			latency,
			c.Writer.Status(),
			c.Errors.String(),
		))
	}
}
