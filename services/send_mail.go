package services

import (
	"fmt"
	"net/mail"
	"net/smtp"
)

func SendEmail(to, subject, body string) error {
	// 1. Configure SMTP settings (use environment variables!)
	from := "karanidancan120@gmail.com"
	password := "uyct osxl wphg ymdg" // ⚠️ Never hardcode passwords!
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

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
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" + body + "\r\n",
	)

	// 5. Send the email
	err := smtp.SendMail(
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