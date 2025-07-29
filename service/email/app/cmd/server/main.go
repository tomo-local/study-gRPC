package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"email-service/config"
	"email-service/db"
	"email-service/smtp"
)

func main() {
	log.Println("Starting Email Service...")

	// Load configuration
	config.Load()

	// Initialize database
	db.Init()
	defer db.Close()

	// Create SMTP server
	smtpServer := smtp.NewServer()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start SMTP server in a goroutine
	go func() {
		if err := smtpServer.Start(); err != nil {
			log.Printf("SMTP server error: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received")

	// Stop SMTP server
	if err := smtpServer.Stop(); err != nil {
		log.Printf("Error stopping SMTP server: %v", err)
	}

	log.Println("Email Service stopped")
}
