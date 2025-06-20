package services

import (
	"fmt"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func SendEmail(to, subject, body string) error {
	// 1. Get email configuration
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("failed to load .env file: %v", err)
	}
	from := os.Getenv("EMAIL")
	password := os.Getenv("SMTP_PASSWORD")
	smtpHost := "smtp.gmail.com"
	smtpPort := os.Getenv("SMTP_PORT")

	// 2. Validate email format
	if _, err := mail.ParseAddress(to); err != nil {
		return fmt.Errorf("invalid recipient email: %v", err)
	}

	// 3. Set up authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// 4. Compose the email (MIME format)
	msg := []byte(
		"From: " + from + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=UTF-8\r\n" +
			"\r\n" + body + "\r\n",
	)

	// 5. Send the email
	err = smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		msg,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}