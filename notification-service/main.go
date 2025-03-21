package main

import (
	"log"
	"net"
	"notification_service/config"
	"notification_service/handler"
	"notification_service/proto/notification"
	"notification_service/queue"
	"notification_service/repository"
	"notification_service/usecase"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {

	config.InitDB()
	config.InitRabbitMQ()
	defer config.CloseRabbitMQ()

	logger := logrus.New()

	repo := repository.NewNotificationRepository(config.DB, logger)
	publisher := queue.NewRabbitMQPublisher()
	uc := usecase.NewNotificationUsecase(repo, publisher, logger)
	grpcHandler := handler.NewNotificationGRPCHandler(uc, logger)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	notification.RegisterNotificationServiceServer(s, grpcHandler)

	logger.Info("gRPC server started on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
