package main

import (
	"notification_service/config"
	"notification_service/queue"
	"notification_service/repository"
	"notification_service/usecase"

	"github.com/sirupsen/logrus"
)

func main() {

	db := config.InitDB()

	logger := logrus.New()

	// migration.Migration(db)

	config.InitRabbitMQ()
	defer config.CloseRabbitMQ()

	notificationRepo := repository.NewNotificationRepository(db, logger)
	notificationUseCase := usecase.NewNotificationUsecase(notificationRepo, logger)

	go queue.StartConsumer(config.RabbitMQConn, notificationUseCase, logger)

	logger.Info("Notification Service is running...")
	select {}
}
