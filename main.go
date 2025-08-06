package main

import (
	"log"
	"strconv"
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

	smtpPort, err := strconv.Atoi(infrastructure.Env.EMAIL_PORT)
if err != nil {
	log.Fatalf("Invalid EMAIL_PORT: %v", err)
}

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

	userUsecase := usecases.NewUserUseCase(userRepository, refreshRepository, resetPasswordRepo, jwtService, passwordService, emailService, 3*time.Second)
	authUsecase := usecases.NewAuthUsecase(jwtService, userRepository, refreshRepository, 3*time.Second)
	blogUsecase := usecases.NewBlogUseCase(blogRepository)

	userController := controllers.NewUserController(userUsecase)
	authController := controllers.NewAuthController(authUsecase)
	blogController := controllers.NewBlogController(blogUsecase)

	r := router.SetupRouter(userController, authController, blogController, authMiddleware)

	port := infrastructure.Env.PORT
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
