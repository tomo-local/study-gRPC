package model

import (
	"time"
)

// PasswordResetToken はパスワードリセットトークンを表すモデル
type PasswordResetToken struct {
	Timestamp
	ID        int64     `gorm:"primarykey"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Token     string    `gorm:"uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
}

// IsExpired はトークンが期限切れかどうかを確認する
func (p *PasswordResetToken) IsExpired() bool {
	return time.Now().After(p.ExpiresAt)
}
