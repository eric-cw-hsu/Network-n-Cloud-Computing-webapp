package middleware

import (
	"github.com/gin-gonic/gin"
)

func RemovePayloadForMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == 405 {
			// add no-cache header
			c.Writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")

			// remove default payload
			c.Writer.Write([]byte(""))
		}
	}
}
