package repository

import (
	"notification_service/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type INotificationRepository interface {
	Create(notification *model.Notification) error
	MarkAsSent(id int) error
}

type notificationRepository struct {
	db     *gorm.DB
	logger *logrus.Logger
}

const (
	StatusPending = "pending"
	StatusSent    = "sent"
)

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
		"notification_id": notification.NotificationID,
		"email":           notification.Email,
		"status":          notification.Status,
	}).Info("Notification created")
	return nil
}

func (r *notificationRepository) MarkAsSent(id int) error {
	tx := r.db.Begin()
	if err := tx.Model(&model.Notification{}).Where("notification_id = ?", id).Update("status", StatusSent).Error; err != nil {
		tx.Rollback()
		r.logger.WithFields(logrus.Fields{
			"notification_id": id,
			"error":           err.Error(),
		}).Error("Failed to mark notification as sent")
		return err
	}
	tx.Commit()

	r.logger.WithField("notification_id", id).Info("Notification marked as sent")
	return nil
}
