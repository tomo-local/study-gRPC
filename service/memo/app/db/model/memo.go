package model

import "time"

type FileType string

const (
	FileTypeTxt  FileType = "txt"
	FileTypeMd   FileType = "md"
	FileTypeJson FileType = "json"
)

type Memo struct {
	ID         string    `json:"id"`
	FileType   FileType  `json:"file_type"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	ModifiedAt time.Time `json:"modified_at"`
}
