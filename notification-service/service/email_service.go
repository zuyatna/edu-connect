package service

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {

	goenvload := godotenv.Load()
	if goenvload != nil {
		log.Println(goenvload.Error())
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "no-reply@educonnect.com")
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(os.Getenv("EMAIL_HOST"), 587, os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"))

	if err := dialer.DialAndSend(mailer); err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	return nil
}
