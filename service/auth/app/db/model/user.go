package model

// User はユーザー情報を表すモデル
type User struct {
	Timestamp
	ID            string `gorm:"primaryKey" json:"id"`
	Email         string `gorm:"uniqueIndex;not null;size:320" json:"email"`
	Name          string `gorm:"not null;size:255" json:"name"`
	PasswordHash  string `gorm:"not null;size:255" json:"-"`
	EmailVerified bool   `gorm:"default:false" json:"email_verified"`
}
