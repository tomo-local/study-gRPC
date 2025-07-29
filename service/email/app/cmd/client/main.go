package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func main() {
	// SMTP server configuration
	smtpHost := "localhost"
	smtpPort := "2525"

	// Sender info
	from := "test@example.com"
	password := "password"

	// Recipients
	to := []string{"admin@example.com"}

	// Message
	subject := "Test Email from Go SMTP Client"
	body := "This is a test email sent from the Go SMTP client.\n\nRegards,\nEmail Service"

	// Create message
	message := fmt.Sprintf("From: %s\r\n", from)
	message += fmt.Sprintf("To: %s\r\n", to[0])
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n"
	message += body

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Send email
	addr := smtpHost + ":" + smtpPort
	err := smtp.SendMail(addr, auth, from, to, []byte(message))
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return
	}

	fmt.Println("Email sent successfully!")
	fmt.Printf("From: %s\n", from)
	fmt.Printf("To: %s\n", to[0])
	fmt.Printf("Subject: %s\n", subject)
}
