package router

import (
	controllers "blog-api/Delivery/Controllers"

	"github.com/gin-gonic/gin"
)

func SetupRouter(uc *controllers.UserController, ac *controllers.AuthController) *gin.Engine {
	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", uc.RegisterHandler)
		authRoutes.POST("/login", uc.LoginHandler)
		authRoutes.POST("/refresh", ac.RefreshTokenHandler)
	}

	return router
}
