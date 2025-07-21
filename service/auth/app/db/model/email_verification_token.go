package model

import (
	"time"
)

// EmailVerificationToken はメールアドレス認証トークンを表すモデル
type EmailVerificationToken struct {
	Timestamp
	ID        int64     `gorm:"primarykey"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// IsExpired はトークンが期限切れかどうかを確認する
func (e *EmailVerificationToken) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}
