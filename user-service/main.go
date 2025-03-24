package main

import (
	"os"
	"userService/config"
	"userService/handler"
	"userService/queue"
	"userService/repository"
	"userService/route"
	"userService/usecase"

	"github.com/labstack/echo/v4"
)

func main() {

	db := config.InitDB()

	rabbitConn, channel := config.InitRabbitMQ()
	defer rabbitConn.Close()
	defer channel.Close()

	// migration.Migration(db)

	emailPublisher, err := queue.NewEmailPublisher(channel, "email")
	if err != nil {
		panic("Failed to initialize email publisher: " + err.Error())
	}

	userRepo := repository.NewUserRepository(db)
	verificationRepo := repository.NewVerificationRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	verificationUC := usecase.NewVerificationUseCase(userRepo, verificationRepo, emailPublisher)
	userUC := usecase.NewUserUseCase(userRepo, verificationUC)
	passwordResetUC := usecase.NewPasswordResetUseCase(userRepo, passwordResetRepo, emailPublisher)

	userHandler := handler.NewUserHandler(userUC)
	verificationHandler := handler.NewVerificationHandler(verificationUC)
	passwordResetHandler := handler.NewPasswordResetHandler(passwordResetUC)

	e := echo.New()
	route.Init(e, userHandler, *verificationHandler, *passwordResetHandler)
	e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

}
