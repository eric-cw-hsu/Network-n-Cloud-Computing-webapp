package middleware

import (
	"go-template/internal/auth/domain/basic"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(basicService *basic.BasicService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		basicToken := strings.Split(authHeader, " ")
		if len(basicToken) != 2 || basicToken[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		user, err := basicService.Authenticate(basicToken[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
