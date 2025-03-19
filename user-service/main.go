package main

import (
	"os"
	"userService/config"
	"userService/handler"
	"userService/repository"
	"userService/route"
	"userService/usecase"

	"github.com/labstack/echo/v4"
)

func main() {

	db := config.InitDB()

	// migration.Migration(db)

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)

	e := echo.New()
	route.Init(e, userHandler)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
