package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Note struct {
	ID        string         `gorm:"primaryKey;type:varchar(255)" json:"id" db:"id"`
	Title     string         `gorm:"type:varchar(255);not null" json:"title" db:"title"`
	Content   string         `gorm:"type:text" json:"content" db:"content"`
	Category  string         `gorm:"type:varchar(100)" json:"category" db:"category"`
	Tags      pq.StringArray `gorm:"type:text[]" json:"tags" db:"tags"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at" db:"updated_at"`
}

// BeforeCreate はノート作成前にUUIDを自動生成する
func (n *Note) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = generateNoteID()
	}
	return nil
}

// generateNoteID はユニークなノートIDを生成する
func generateNoteID() string {
	return uuid.New().String()
}
