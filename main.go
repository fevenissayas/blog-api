package main

import (
	"time"

	controllers "blog-api/Delivery/Controllers"
	router "blog-api/Delivery/Router"
	infrastructure "blog-api/Infrastructure"
	repositories "blog-api/Repositories"
	usecases "blog-api/Usecases"
)

func main() {

	infrastructure.LoadEnv()

	client := repositories.NewMongoClient()
	db := client.Database(infrastructure.Env.DB_NAME)

	jwtService := infrastructure.NewJWTService()
	passwordService := infrastructure.NewPasswordService()
	authMiddleware := infrastructure.NewAuthMiddleware(jwtService)

	
	smtpPort := infrastructure.ParsePort(infrastructure.Env.EMAIL_PORT, 465)


	emailService := infrastructure.NewSMTPEmailService(
		infrastructure.Env.EMAIL_FROM,
		infrastructure.Env.EMAIL_HOST,
		smtpPort,
		infrastructure.Env.EMAIL_USERNAME,
		infrastructure.Env.EMAIL_PASSWORD,
	)

	userRepository := repositories.NewUserRepository(db)
	refreshRepository := repositories.NewRefreshTokenRepository(db)
	blogRepository := repositories.NewBlogRepository(db)
	resetPasswordRepo := repositories.NewPasswordResetTokenRepo(db)
	likeRepo :=repositories.NewLikeRepository(db)
	commentRepository := repositories.NewCommentRepository(db)

    Aiservice:= infrastructure.NewAiService()

	userUsecase := usecases.NewUserUseCase(userRepository, refreshRepository, resetPasswordRepo, jwtService, passwordService, emailService, 3*time.Second)
	authUsecase := usecases.NewAuthUsecase(jwtService, userRepository, refreshRepository, 3*time.Second)
	blogUsecase := usecases.NewBlogUseCase(blogRepository,Aiservice)
	likeUsecase := usecases.NewLikeUsecase(likeRepo)
	commentUsecase := usecases.NewCommentUsecase(commentRepository)

	userController := controllers.NewUserController(userUsecase)
	authController := controllers.NewAuthController(authUsecase)
	blogController := controllers.NewBlogController(blogUsecase)
	likeController := controllers.NewLikeController(likeUsecase)
	commentController := controllers.NewCommentController(commentUsecase)

	r := router.SetupRouter(userController, authController, blogController, likeController,authMiddleware, commentController)

	port := infrastructure.Env.PORT
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
