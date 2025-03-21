package queue

import (
	"encoding/json"
	"fmt"
	"notification_service/config"
	"notification_service/model"
	"notification_service/repository"
	"notification_service/service"
)

func StartConsumer(repo repository.INotificationRepository) {
	ch, _ := config.RabbitMQConn.Channel()

	msgs, _ := ch.Consume(
		"email_queue", "", true, false, false, false, nil,
	)

	for msg := range msgs {
		var notification model.Notification
		json.Unmarshal(msg.Body, &notification)

		err := service.SendEmail(notification.Email, notification.Subject, notification.Message)
		if err == nil {
			repo.MarkAsSent(uint(notification.ID))
			fmt.Println("Email sent to:", notification.Email)
		} else {
			fmt.Println("Failed to send email:", err)
		}
	}
}
