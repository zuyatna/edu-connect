package handler

import (
	"context"
	"notification_service/proto/notification"
	"notification_service/usecase"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationGRPCHandler struct {
	notification.UnimplementedNotificationServiceServer
	usecase usecase.INotificationUsecase
	logger  *logrus.Logger
}

func NewNotificationGRPCHandler(usecase usecase.INotificationUsecase, logger *logrus.Logger) *NotificationGRPCHandler {
	return &NotificationGRPCHandler{
		usecase: usecase,
		logger:  logger,
	}
}

func (h *NotificationGRPCHandler) SendNotification(ctx context.Context, req *notification.SendNotificationRequest) (*notification.SendNotificationResponse, error) {
	err := h.usecase.SendNotification(req.Email, req.Subject, req.Message)
	if err != nil {
		h.logger.WithError(err).Error("Failed to send notification via gRPC")
		return nil, status.Errorf(codes.Internal, "failed to send notification")
	}

	h.logger.WithField("email", req.Email).Info("Notification sent via gRPC")

	return &notification.SendNotificationResponse{
		Status: "sent",
	}, nil
}
