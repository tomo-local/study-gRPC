package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	config "memo/config/server"
	"memo/db/model"
)

// FileService defines the interface for file-based memo operations.
type FileService interface {
	CreateFile(memo *model.Memo) (*model.Memo, error)
	UpdateFile(memo *model.Memo) (*model.Memo, error)
	DeleteFile(id string) error
}

type fileService struct {
	folderPath string
}

// GetService creates a new FileService.
func GetService(config *config.Config) (FileService, error) {
	return &fileService{
		folderPath: config.FolderPath,
	}, nil
}

// CreateFile creates a new file for the given memo.
func (f *fileService) CreateFile(memo *model.Memo) (*model.Memo, error) {
	if err := os.MkdirAll(f.folderPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}
	filename := fmt.Sprintf("%s_%s.%s", memo.ID, memo.Title, memo.FileType)
	filePath := filepath.Join(f.folderPath, filename)
	content := fmt.Sprintf("Title: %s\nContent: %s\nCreatedAt: %s\nUpdatedAt: %s\n", memo.Title, memo.Content, memo.CreatedAt, memo.UpdatedAt)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	return memo, nil
}

// UpdateFile updates the file for the given memo.
// It currently reuses CreateFile, which overwrites the existing file.
func (f *fileService) UpdateFile(memo *model.Memo) (*model.Memo, error) {
	return f.CreateFile(memo)
}

// DeleteFile deletes a memo file by its ID.
func (f *fileService) DeleteFile(id string) error {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file with id %s not found (directory does not exist)", id)
		}
		return fmt.Errorf("failed to read directory: %w", err)
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), id+"_") {
			filePath := filepath.Join(f.folderPath, entry.Name())
			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("failed to remove file %s: %w", filePath, err)
			}
			return nil
		}
	}
	return fmt.Errorf("file with id %s not found", id)
}
