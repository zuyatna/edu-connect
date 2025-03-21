package queue

import (
	"encoding/json"
	"log"
	"notification_service/config"
	"notification_service/model"

	"github.com/streadway/amqp"
)

type RabbitMQPublisher interface {
	Publish(notification *model.Notification) error
}

type rabbitMQPublisher struct {
	channel *amqp.Channel
}

func NewRabbitMQPublisher() RabbitMQPublisher {
	ch, err := config.RabbitMQConn.Channel()
	if err != nil {
		log.Fatal("Failed to open channel:", err)
	}

	return &rabbitMQPublisher{channel: ch}
}

func (p *rabbitMQPublisher) Publish(notification *model.Notification) error {
	body, _ := json.Marshal(notification)

	return p.channel.Publish(
		"", "email_queue",
		false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}
