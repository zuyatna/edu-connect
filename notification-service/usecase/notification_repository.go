package usecase

import (
	"notification_service/model"
	"notification_service/repository"

	"github.com/sirupsen/logrus"
)

type INotificationUsecase interface {
	SendNotification(email, subject, message string) error
}

type notificationUsecase struct {
	repo   repository.INotificationRepository
	queue  queue.IRabbitMQPublisher
	logger *logrus.Logger
}

func NewNotificationUsecase(
	repo repository.INotificationRepository,
	queue queue.IRabbitMQPublisher,
	logger *logrus.Logger,
) INotificationUsecase {
	return &notificationUsecase{
		repo:   repo,
		queue:  queue,
		logger: logger,
	}
}

func (u *notificationUsecase) SendNotification(email, subject, message string) error {
	notification := &model.Notification{
		Email:   email,
		Subject: subject,
		Message: message,
		Status:  "pending",
	}

	if err := u.repo.Create(notification); err != nil {
		u.logger.WithError(err).Error("Failed to create notification in usecase")
		return err
	}

	if err := u.queue.Publish(notification); err != nil {
		u.logger.WithError(err).Error("Failed to publish notification to RabbitMQ")
		return err
	}

	u.logger.WithFields(logrus.Fields{
		"id":    notification.ID,
		"email": email,
	}).Info("Notification sent to queue")

	return nil
}
