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
	userRepository := repositories.NewUserRepository(db)
	refreshRepository := repositories.NewRefreshTokenRepository(db)
	userUsecase := usecases.NewUserUseCase(userRepository,refreshRepository, jwtService, passwordService, 3*time.Second)
	authUsecase :=usecases.NewAuthUsecase(jwtService,userRepository,refreshRepository,3*time.Second)
	userController := controllers.NewUserController(userUsecase)
	authController := controllers.NewAuthController(authUsecase)
	r := router.SetupRouter(userController,authController)

	port := infrastructure.Env.PORT
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
