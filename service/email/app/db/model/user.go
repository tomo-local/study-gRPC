package model

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// User represents a user in the email system
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"`
	FullName     string    `json:"full_name"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName overrides the table name for GORM
func (User) TableName() string {
	return "users"
}

// Validate validates the user data
func (u *User) Validate() error {
	if strings.TrimSpace(u.Email) == "" {
		return errors.New("email is required")
	}

	if strings.TrimSpace(u.PasswordHash) == "" {
		return errors.New("password hash is required")
	}

	// Basic email format validation
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}

	return nil
}

// BeforeCreate hook is called before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.Validate(); err != nil {
		return err
	}

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	if u.IsActive == false && tx.Statement.Changed("IsActive") == false {
		u.IsActive = true // Default value
	}

	return nil
}

// BeforeUpdate hook is called before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if err := u.Validate(); err != nil {
		return err
	}

	u.UpdatedAt = time.Now()
	return nil
}

// SetPassword sets the password hash for the user
func (u *User) SetPassword(passwordHash string) {
	u.PasswordHash = passwordHash
}

// IsEmailValid checks if the email format is valid
func (u *User) IsEmailValid() bool {
	email := strings.TrimSpace(u.Email)
	return email != "" && strings.Contains(email, "@") && strings.Contains(email, ".")
}
