package middleware

import (
	"github.com/gin-gonic/gin"
)

func RemovePayloadForMethodNotAllowed() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if c.Writer.Status() == 405 {
			// remove default payload
			c.Writer.Write([]byte(""))

			// // add no-cache header
			c.Writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		}
	}
}
