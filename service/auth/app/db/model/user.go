package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User はユーザー情報を表すモデル
type User struct {
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email         string    `gorm:"uniqueIndex;not null;size:320" json:"email"`
	Name          string    `gorm:"not null;size:255" json:"name"`
	PasswordHash  string    `gorm:"not null;size:255" json:"-"`
	EmailVerified bool      `gorm:"default:false" json:"email_verified"`
	CreatedAt     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"not null" json:"updated_at"`
}

// BeforeCreate は新しいユーザーを作成する前に実行される
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return nil
}

// PasswordResetToken はパスワードリセットトークンを表すモデル
type PasswordResetToken struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// BeforeCreate は新しいパスワードリセットトークンを作成する前に実行される
func (p *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return nil
}

// IsExpired はトークンが期限切れかどうかを確認する
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}

// EmailVerificationToken はメールアドレス認証トークンを表すモデル
type EmailVerificationToken struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// BeforeCreate は新しいメールアドレス認証トークンを作成する前に実行される
func (e *EmailVerificationToken) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return nil
}

// IsExpired はトークンが期限切れかどうかを確認する
func (e *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}
