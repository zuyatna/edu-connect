package queue

import (
	"encoding/json"
	"notification_service/model"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type NotificationProcessor interface {
	SendNotification(notification model.Notification) error
}

func StartConsumer(conn *amqp091.Connection, uc NotificationProcessor, logger *logrus.Logger) {
	ch, err := conn.Channel()
	if err != nil {
		logger.Fatal("Failed to open RabbitMQ channel:", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"email",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal("Failed to declare queue:", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Fatal("Failed to register consumer:", err)
	}

	logger.Info("Waiting for messages from RabbitMQ...")

	for msg := range msgs {
		var notification model.Notification
		err := json.Unmarshal(msg.Body, &notification)
		if err != nil {
			logger.Error("Failed to unmarshal message:", err)
			continue
		}

		err = uc.SendNotification(notification)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"email": notification.Email,
				"error": err.Error(),
			}).Error("Failed to process notification")
		}
	}
}
