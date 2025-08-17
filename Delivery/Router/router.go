package router

import (
	controllers "blog-api/Delivery/Controllers"
	infrastructure "blog-api/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	uc *controllers.UserController,
	ac *controllers.AuthController,
	bc *controllers.BlogController,
	likeCtrl *controllers.LikeController,
	authMiddleware *infrastructure.AuthMiddleware,
	commentsController *controllers.CommentController,
) *gin.Engine {
	router := gin.Default()

	// --- Global CORS middleware ---
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:3000") // TODO: change to frontend domain in production
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// --- Auth routes ---
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/register", uc.RegisterHandler)
		authRoutes.POST("/verify-email", uc.VerifyEmailHandler)
		authRoutes.POST("/login", uc.LoginHandler)
		authRoutes.POST("/refresh", ac.RefreshTokenHandler)
		authRoutes.POST("/logout", authMiddleware.Middleware(), uc.LogoutHandler)
		authRoutes.POST("/update", authMiddleware.Middleware(), uc.UpdateProfile)

		// Admin-only subroutes
		adminRoutes := authRoutes.Group("/admin")
		adminRoutes.Use(authMiddleware.Middleware())
		{
			adminRoutes.POST("/register", uc.RegisterHandler) // Admin can register other admins
			adminRoutes.POST("/promote", uc.Promote)
		}
	}

	// --- Password reset routes ---
	router.POST("/password/request-reset", uc.RequestPasswordResetHandler)
	router.POST("/password/reset", uc.ResetPasswordHandler)

	// --- Blog routes ---
	blogRoutes := router.Group("/blogs")
	{
		blogRoutes.GET("/", bc.GetBlogsHandler)       // Paginated blogs
		blogRoutes.GET("/:id", bc.GetBlogByIDHandler) // Single blog

		blogRoutes.POST("/", authMiddleware.Middleware(), bc.CreateBlogHandler)
		blogRoutes.PUT("/:id", authMiddleware.Middleware(), bc.UpdateBlogHandler)
		blogRoutes.DELETE("/:id", authMiddleware.Middleware(), bc.DeleteBlog)
		blogRoutes.GET("/filter", authMiddleware.Middleware(), bc.FilterBlogsHandler)
		blogRoutes.POST("/aisuggestion", authMiddleware.Middleware(), bc.AiSuggestion)
		blogRoutes.GET("/search", authMiddleware.Middleware(), bc.SearchBlogs)

		// Likes
		likes := blogRoutes.Group("/:id/likes", authMiddleware.Middleware())
		{
			likes.POST("/", likeCtrl.LikeBlogHandler)
			likes.DELETE("/", likeCtrl.RemoveLikeBlogHandler)
			likes.GET("/", likeCtrl.GetLikeCountHandler)
			likes.GET("/is-liked", likeCtrl.IsBlogLikedHandler)
		}

		// Comments
		comments := blogRoutes.Group("/:id/comments", authMiddleware.Middleware())
		{
			comments.POST("/", commentsController.CreateComment)
			comments.GET("/", commentsController.GetComments)
		}
	}

	// Comment deletion (separate for direct access)
	router.DELETE("/comments/:commentID", authMiddleware.Middleware(), commentsController.DeleteComment)

	return router
}
