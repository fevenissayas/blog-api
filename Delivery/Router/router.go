package router

import (
	controller "blog-api/Delivery/Controllers"
	"blog-api/Infrastructure"
	"blog-api/Middleware"

	"github.com/gin-gonic/gin"
)

func SetUpRouter(c *controller.UserController, jwtService *infrastructure.JwtService) *gin.Engine {
	router := gin.Default()

	// Public routes
	router.POST("/register", c.Register)
	router.POST("/login", c.Login)

	// Protected routes
	authMiddleware := middleware.AuthMiddleware(jwtService)
	protected := router.Group("/api")
	protected.Use(authMiddleware)

	return router
}
