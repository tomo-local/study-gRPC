package model

import (
	"fmt"
	"path/filepath"
	"time"
)

type FileType string

const (
	FileTypeTxt  FileType = "txt"
	FileTypeMd   FileType = "md"
	FileTypeJson FileType = "json"
)

type Memo struct {
	ID        string    `json:"id"`
	FileType  FileType  `json:"file_type"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (m *Memo) GetFilePath(folderPath string) string {
	return filepath.Join(folderPath, fmt.Sprintf("%s_%s.%s", m.Title, m.ID, m.FileType))
}
