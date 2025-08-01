package main

import (
	controllers "blog-api/Delivery/Controllers"
	router "blog-api/Delivery/Router"
	infrastructure "blog-api/Infrastructure"
	repositories "blog-api/Repositories"
	usecases "blog-api/Usecases"
)

func main() {
	infrastructure.LoadEnv()
	infrastructure.InitMongo()
	userRepo := repositories.NewUserRepo(infrastructure.UserCollection())
    JwtService := infrastructure.NewJwtService()
    userUsecase := usecases.NewUserUsecase(userRepo, JwtService)

    userController := controllers.NewUserController(userUsecase)

    router := router.SetUpRouter(userController)

    router.Run(":8080")
}
