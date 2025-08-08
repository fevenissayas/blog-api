package router

import (
	controllers "blog-api/Delivery/Controllers"
	infrastructure "blog-api/Infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(uc *controllers.UserController, ac *controllers.AuthController, bc *controllers.BlogController,likeCtrl *controllers.LikeController ,authMiddleware *infrastructure.AuthMiddleware, commentsController *controllers.CommentController) *gin.Engine {
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
		blogRoutes.GET("/", bc.GetBlogsHandler)           // Get paginated blogs
		blogRoutes.GET("/:id", bc.GetBlogByIDHandler)     // Get single blog

		blogRoutes.POST("/",authMiddleware.Middleware(),bc.CreateBlogHandler)
		blogRoutes.PUT("/:id",authMiddleware.Middleware(),bc.UpdateBlogHandler)
		blogRoutes.DELETE("/:id",authMiddleware.Middleware(),bc.DeleteBlog)
		blogRoutes.GET("/filter", authMiddleware.Middleware(), bc.FilterBlogsHandler)
		blogRoutes.POST("/aisuggestion", authMiddleware.Middleware(), bc.AiSuggestion)
		blogRoutes.GET("/search", authMiddleware.Middleware(), bc.SearchBlogs) // this shall be corrected 

		likes := blogRoutes.Group("/:id/likes",authMiddleware.Middleware())
		{
			likes.POST("/", likeCtrl.LikeBlogHandler)
			likes.DELETE("/", likeCtrl.UnlikeBlogHandler)
			likes.GET("/", likeCtrl.GetLikeCountHandler)
			likes.GET("/is-liked", likeCtrl.IsBlogLikedHandler)
		}

		// Comment endpoints
		comments := blogRoutes.Group(":blogID/comments", authMiddleware.Middleware())
		{
			comments.POST("/", commentsController.CreateComment)
			comments.GET("/", commentsController.GetComments)
		}
		router.DELETE("/comments/:commentID", authMiddleware.Middleware(), commentsController.DeleteComment)
	}

	return router
}
