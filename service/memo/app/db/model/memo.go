package model

type FileType string

const (
	FileTypeTxt  FileType = "txt"
	FileTypeMd   FileType = "md"
	FileTypeJson FileType = "json"
)

type Memo struct {
	ID        string   `json:"id"`
	FileType  FileType `json:"file_type"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}
