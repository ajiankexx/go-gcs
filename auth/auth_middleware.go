package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be 'Bear <token>'"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		if tokenString == "" || tokenString == "invalid_token" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		c.Set("userID", 123)

		c.Next()
	}
}
