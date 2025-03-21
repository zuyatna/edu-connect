package service

import (
	"fmt"
	"notification_service/config"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "your-email@gmail.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	d := gomail.NewDialer(
		config.EmailSettings.SMTPHost,
		config.EmailSettings.SMTPPort,
		config.EmailSettings.Username,
		config.EmailSettings.Password,
	)

	// dialer := gomail.NewDialer("smtp.gmail.com", 587, "your-email@gmail.com", "your-email-password")

	if err := d.DialAndSend(mailer); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}
