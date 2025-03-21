package repository

import (
	"notification_service/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	Create(notification *model.Notification) error
	MarkAsSent(id uint) error
}

type notificationRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewNotificationRepository(db *gorm.DB, logger *logrus.Logger) INotificationRepository {
	return &notificationRepository{
		db:     db,
		logger: logger,
	}
}

func (r *notificationRepository) Create(notification *model.Notification) error {
	if err := r.db.Create(notification).Error; err != nil {
		r.logger.WithFields(logrus.Fields{
			"email": notification.Email,
			"error": err,
		}).Error("Failed to create notification")
		return err
	}
	r.logger.WithFields(logrus.Fields{
		"id":    notification.ID,
		"email": notification.Email,
	}).Info("Notification created")
	return nil
}

func (r *notificationRepository) MarkAsSent(id uint) error {
	if err := r.db.Model(&model.Notification{}).Where("id = ?", id).Update("status", "sent").Error; err != nil {
		r.logger.WithField("id", id).Error("Failed to mark notification as sent")
		return err
	}
	r.logger.WithField("id", id).Info("Notification marked as sent")
	return nil
}
