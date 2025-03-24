package main

import (
	"os"
	"sync"
	"userService/config"
	"userService/grpc"
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

	grpcPort := os.Getenv("GRPC_PORT")

	if grpcPort == "" {
		grpcPort = "50051"
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		grpc.StartGRPCServer(userRepo, grpcPort)
	}()

	go func() {
		defer wg.Done()

		e := echo.New()
		route.Init(e, userHandler, *verificationHandler, *passwordResetHandler)

		e.Logger.Fatal(e.Start(":" + os.Getenv("PORT")))

	}()

	wg.Wait()

}
