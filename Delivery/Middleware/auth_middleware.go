package middleware

import (
	"blog-api/Domain"
	"blog-api/Infrastructure"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtService *infrastructure.JwtService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		// Expecting format "Bearer <token>"
		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenStr := fields[1]

		// Parse and validate token
		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil || claims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Step 4: Add user data to context
		ctx := context.WithValue(c.Request.Context(), "user", claims)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
