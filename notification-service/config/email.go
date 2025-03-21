package config

type EmailConfig struct {
	SMTPHost string
	SMTPPort int
	Username string
	Password string
	From     string
}

var EmailSettings = EmailConfig{
	SMTPHost: "smtp.gmail.com",
	SMTPPort: 587,
	Username: "your-email@gmail.com",
	Password: "your-app-password",
	From:     "your-email@gmail.com",
}
