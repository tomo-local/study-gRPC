package service

import (
	"errors"
	"time"

	"email-service/db"
	"email-service/db/model"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// EmailService provides email-related operations
type EmailService struct {
	db *gorm.DB
}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{
		db: db.GetDB(),
	}
}

// CreateEmail creates a new email record
func (s *EmailService) CreateEmail(fromEmail string, toEmails []string, subject, bodyText, bodyHTML string) (*model.Email, error) {
	email := &model.Email{
		MessageID: generateMessageID(),
		FromEmail: fromEmail,
		ToEmails:  pq.StringArray(toEmails),
		Subject:   subject,
		BodyText:  bodyText,
		BodyHTML:  bodyHTML,
		Status:    "pending",
	}

	err := s.db.Create(email).Error
	if err != nil {
		return nil, err
	}

	return email, nil
}

// GetEmailByID retrieves an email by ID
func (s *EmailService) GetEmailByID(id uint) (*model.Email, error) {
	var email model.Email
	err := s.db.First(&email, id).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}

// GetEmailsByUser retrieves emails for a specific user
func (s *EmailService) GetEmailsByUser(userEmail string) ([]model.Email, error) {
	var emails []model.Email
	err := s.db.Where("from_email = ? OR ? = ANY(to_emails)", userEmail, userEmail).Find(&emails).Error
	if err != nil {
		return nil, err
	}
	return emails, nil
}

// UpdateEmailStatus updates the status of an email
func (s *EmailService) UpdateEmailStatus(id uint, status string) error {
	return s.db.Model(&model.Email{}).Where("id = ?", id).Update("status", status).Error
}

// DeleteEmail deletes an email
func (s *EmailService) DeleteEmail(id uint) error {
	return s.db.Delete(&model.Email{}, id).Error
}

// UserService provides user-related operations
type UserService struct {
	db *gorm.DB
}

// NewUserService creates a new user service
func NewUserService() *UserService {
	return &UserService{
		db: db.GetDB(),
	}
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := s.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new user
func (s *UserService) CreateUser(email, passwordHash, fullName string) (*model.User, error) {
	user := &model.User{
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
		IsActive:     true,
	}

	err := s.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Helper function to generate message ID
func generateMessageID() string {
	return time.Now().Format("20060102150405") + "@localhost"
}
