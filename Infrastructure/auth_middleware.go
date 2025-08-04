package infrastructure

import (
	domain "blog-api/Domain"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authHeader = "Authorization"
)

type AuthMiddleware struct {
	JWTService domain.IJWTService
}

func NewAuthMiddleware(JWTService domain.IJWTService) *AuthMiddleware {
	return &AuthMiddleware{JWTService: JWTService}
}

func (a *AuthMiddleware) MiddleWare() gin.HandlerFunc {

	return func(c *gin.Context) {
		auth := c.GetHeader(authHeader)
		if auth == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
			return
		}

		authParts := strings.Split(auth, " ")
		if len(authParts) != 2 || strings.ToLower(authParts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization token"})
			return
		}

		claims, err := a.JWTService.ValidateAccessToken(authParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "unauthorized"})
			return
		}

		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
