package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * EmptyQueryParameterChecker is a middleware to check if the request has query parameter.
 * If the request has query parameter, it will return 400 Bad Request.
 */
func EmptyQueryParameterChecker() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.RawQuery != "" {
			c.Status(http.StatusBadRequest)
			c.Abort()
			return
		}
		c.Next()
	}
}
