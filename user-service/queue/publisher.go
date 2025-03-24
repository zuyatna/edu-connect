package queue

import (
	"encoding/json"
	"os"

	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type IEmailPublisher interface {
	PublishVerificationToken(email, token string) error
	PublishResetPasswordToken(email, token string) error
}

type EmailPublisher struct {
	channel *amqp091.Channel
	queue   amqp091.Queue
}

func NewEmailPublisher(channel *amqp091.Channel, queueName string) (*EmailPublisher, error) {
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EmailPublisher{
		channel: channel,
		queue:   queue,
	}, nil
}

func (p *EmailPublisher) PublishVerificationToken(email, token string) error {

	appURL := os.Getenv("APP_URL")

	verifyLink := appURL + "/v1/verify?token=" + token

	htmlMessage := `
		<p>Halo,</p>
		<p>Silakan klik tombol di bawah ini untuk melakukan verifikasi email Anda:</p>
		<a href="` + verifyLink + `" style="display:inline-block;padding:10px 20px;background-color:#4CAF50;color:#fff;text-decoration:none;border-radius:5px;">Verifikasi Email</a>
		<p>Link: <a href="` + verifyLink + `">` + verifyLink + `</a></p>
	`

	payload := map[string]interface{}{
		"email":   email,
		"subject": "Verifikasi Email",
		"message": htmlMessage,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"", p.queue.Name, false, false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	return err
}

func (p *EmailPublisher) PublishResetPasswordToken(email, token string) error {

	appURL := os.Getenv("APP_URL")

	resetLink := appURL + "/v1/reset-password?token=" + token

	htmlMessage := `
		<p>Halo,</p>
		<p>Silakan klik tombol di bawah ini untuk mengatur ulang password Anda:</p>
		<a href="` + resetLink + `" style="display:inline-block;padding:10px 20px;background-color:#f44336;color:#fff;text-decoration:none;border-radius:5px;">Reset Password</a>
		<p>Atau salin link ini ke browser:</p>
		<p><a href="` + resetLink + `">` + resetLink + `</a></p>
	`

	payload := map[string]interface{}{
		"email":   email,
		"subject": "Reset Password",
		"message": htmlMessage,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = p.channel.Publish(
		"",
		p.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		logrus.WithError(err).Error("Failed to publish reset password email")
		return err
	}

	logrus.WithField("email", email).Info("Reset password email published")
	return nil
}
