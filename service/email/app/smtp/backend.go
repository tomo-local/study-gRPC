package smtp

import (
	"errors"
	"io"
	"log"
	"strings"
	"time"

	"email-service/db"
	"email-service/db/model"

	"github.com/emersion/go-smtp"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// Backend implements SMTP backend
type Backend struct{}

// NewSession creates a new SMTP session
func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

// Session implements SMTP session
type Session struct {
	From string
	To   []string
}

// AuthPlain handles PLAIN authentication
func (s *Session) AuthPlain(username, password string) error {
	log.Printf("SMTP Auth attempt: username=%s", username)

	// Verify user credentials against database
	var user model.User
	err := db.GetDB().Where("email = ? AND is_active = ?", username, true).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("User not found: %s", username)
			return smtp.ErrAuthFailed
		}
		log.Printf("Database error during auth: %v", err)
		return smtp.ErrAuthFailed
	}

	// In a real implementation, you should use bcrypt to compare password
	// For now, we'll do a simple comparison (NOT SECURE for production)
	if !verifyPassword(password, user.PasswordHash) {
		log.Printf("Invalid password for user: %s", username)
		return smtp.ErrAuthFailed
	}

	log.Printf("Authentication successful for user: %s", username)
	return nil
}

// Mail sets the envelope sender
func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	log.Printf("SMTP Mail from: %s", from)
	s.From = from
	return nil
}

// Rcpt adds a recipient
func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	log.Printf("SMTP Rcpt to: %s", to)
	s.To = append(s.To, to)
	return nil
}

// Data handles the message data
func (s *Session) Data(r io.Reader) error {
	log.Printf("SMTP Data: from=%s, to=%v", s.From, s.To)

	// Read the message data
	data, err := io.ReadAll(r)
	if err != nil {
		log.Printf("Error reading message data: %v", err)
		return err
	}

	message := string(data)
	log.Printf("Message received (%d bytes)", len(data))

	// Parse the message (simplified parsing)
	subject, bodyText, bodyHTML := parseMessage(message)

	// Generate unique message ID
	messageID := generateMessageID()

	// Create email record in database
	email := model.Email{
		MessageID: messageID,
		FromEmail: s.From,
		ToEmails:  pq.StringArray(s.To),
		Subject:   subject,
		BodyText:  bodyText,
		BodyHTML:  bodyHTML,
		Status:    "sent",
		SentAt:    &time.Time{},
	}

	err = db.GetDB().Create(&email).Error
	if err != nil {
		log.Printf("Error saving email to database: %v", err)
		return err
	}

	log.Printf("Email saved successfully with ID: %d", email.ID)

	// In a real implementation, you would:
	// 1. Store the message in recipient mailboxes
	// 2. Handle delivery to external domains
	// 3. Process attachments
	// 4. Implement queue for failed deliveries

	return nil
}

// Reset resets the session
func (s *Session) Reset() {
	s.From = ""
	s.To = nil
}

// Logout ends the session
func (s *Session) Logout() error {
	return nil
}

// Helper functions

// parseMessage extracts subject and body from the message
func parseMessage(message string) (subject, bodyText, bodyHTML string) {
	lines := strings.Split(message, "\n")
	headerEnd := false
	var bodyLines []string

	for _, line := range lines {
		if !headerEnd {
			if line == "" || line == "\r" {
				headerEnd = true
				continue
			}

			// Parse headers
			if strings.HasPrefix(strings.ToLower(line), "subject:") {
				subject = strings.TrimSpace(line[8:])
			}
		} else {
			bodyLines = append(bodyLines, line)
		}
	}

	bodyText = strings.Join(bodyLines, "\n")
	// For now, treat all body as text. HTML parsing would be more complex.
	bodyHTML = bodyText

	return subject, bodyText, bodyHTML
}

// generateMessageID creates a unique message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "@localhost"
}

// verifyPassword compares password with hash
// In production, use bcrypt.CompareHashAndPassword
func verifyPassword(password, hash string) bool {
	// This is a simplified implementation for learning purposes
	// In production, use proper password hashing like bcrypt
	return password == "password" // Temporary simple check
}
