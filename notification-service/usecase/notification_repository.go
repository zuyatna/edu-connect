package usecase

import (
	"notification_service/model"
	"notification_service/repository"
	"notification_service/service"

	"github.com/sirupsen/logrus"
)

type INotificationUsecase interface {
	SendNotification(notification model.Notification) error
}

type notificationUsecase struct {
	repo   repository.INotificationRepository
	logger *logrus.Logger
}

func NewNotificationUsecase(
	repo repository.INotificationRepository,
	logger *logrus.Logger,
) INotificationUsecase {
	return &notificationUsecase{
		repo:   repo,
		logger: logger,
	}
}

func (u *notificationUsecase) SendNotification(notification model.Notification) error {

	err := u.repo.Create(&notification)
	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"email": notification.Email,
			"error": err.Error(),
		}).Error("Failed to save notification")
		return err
	}

	err = service.SendEmail(notification.Email, notification.Subject, notification.Message)
	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"email": notification.Email,
			"error": err.Error(),
		}).Error("Failed to send email")
		return err
	}

	err = u.repo.MarkAsSent(notification.NotificationID)
	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"id":    notification.NotificationID,
			"error": err.Error(),
		}).Error("Failed to update notification status")
		return err
	}

	u.logger.WithFields(logrus.Fields{
		"id":    notification.NotificationID,
		"email": notification.Email,
	}).Info("Notification processed successfully")

	return nil
}
