package middleware

import (
	"go-template/internal/auth/domain"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AccountVerificationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, _ := c.Get("user")

		if !user.(*domain.AuthUser).Verify {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Account not verified"})
			c.Abort()
			return
		}

		c.Next()
	}
}
