package router

import (
	controllers "blog-api/Delivery/Controllers"
	infrastructure "blog-api/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(uc *controllers.UserController, ac *controllers.AuthController, bc *controllers.BlogController, authMiddleware *infrastructure.AuthMiddleware) *gin.Engine {
	router := gin.Default()

	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", uc.RegisterHandler)
		authRoutes.POST("/login", uc.LoginHandler)
		authRoutes.POST("/refresh", ac.RefreshTokenHandler)
		authRoutes.POST("/logout", authMiddleware.Middleware(), uc.LogoutHandler)
		authRoutes.POST("/promote", authMiddleware.Middleware(), uc.Promote)
		authRoutes.POST("/update", authMiddleware.Middleware(), uc.UpdateProfile)
	}

	router.POST("/password/request-reset", uc.RequestPasswordResetHandler)
	router.POST("/password/reset", uc.ResetPasswordHandler)

	blogRoutes := router.Group("/blogs")
	{
		blogRoutes.POST("/", authMiddleware.Middleware(), bc.CreateBlogHandler)
		blogRoutes.PUT("/:id", authMiddleware.Middleware(), bc.UpdateBlogHandler)
		blogRoutes.GET("/filter", authMiddleware.Middleware(), bc.FilterBlogsHandler)

	}

	return router
}
