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
	GetFile(id string) (*model.Memo, error)
	ListFiles(maxResults int64) ([]*model.Memo, error)
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
	filename := fmt.Sprintf("%s_%s.%s", memo.Title, memo.ID, memo.FileType)
	filePath := filepath.Join(f.folderPath, filename)
	var content string
	switch memo.FileType {
	case model.FileTypeMd:
		content = generateMarkdownContent(memo)
	case model.FileTypeJson:
		content = generateJsonContent(memo)
	default:
		content = generateTextContent(memo)
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}
	return memo, nil
}

// GetFile retrieves a memo file by its ID.
func (f *fileService) GetFile(id string) (*model.Memo, error) {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasPrefix(entry.Name(), id+"_") {
			filePath := filepath.Join(f.folderPath, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			memo := &model.Memo{
				ID:       id,
				FileType: model.FileType(strings.TrimSuffix(filepath.Ext(entry.Name()), ".")),
				Content:  string(content),
			}
			// Extract title and timestamps from the content
			// This is a simplified example; actual parsing logic may vary
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				memo.Title = strings.TrimPrefix(lines[0], "Title: ")
			}
			if len(lines) > 1 {
				memo.CreatedAt = strings.TrimPrefix(lines[1], "CreatedAt: ")
			}
			if len(lines) > 2 {
				memo.UpdatedAt = strings.TrimPrefix(lines[2], "UpdatedAt: ")
			}
			return memo, nil
		}
	}
	return nil, fmt.Errorf("file with id %s not found", id)
}

// ValidPath checks if the provided path is valid.
func (f *fileService) ListFiles(maxResults int64) ([]*model.Memo, error) {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var memos []*model.Memo
	// If maxResults is less than or equal to 0, return all files
	if maxResults > 0 && int64(len(dirEntries)) > maxResults {
		dirEntries = dirEntries[:maxResults]
	}

	// Iterate through directory entries and filter valid memo files
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.Contains(entry.Name(), "_") {
			filePath := filepath.Join(f.folderPath, entry.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file %s: %w", filePath, err)
			}

			memo := &model.Memo{
				ID:       strings.Split(entry.Name(), "_")[1],
				FileType: model.FileType(strings.TrimSuffix(filepath.Ext(entry.Name()), ".")),
				Content:  string(content),
			}
			// Extract title and timestamps from the content
			lines := strings.Split(string(content), "\n")
			if len(lines) > 0 {
				memo.Title = strings.TrimPrefix(lines[0], "Title: ")
			}
			if len(lines) > 1 {
				memo.CreatedAt = strings.TrimPrefix(lines[1], "CreatedAt: ")
			}
			if len(lines) > 2 {
				memo.UpdatedAt = strings.TrimPrefix(lines[2], "UpdatedAt: ")
			}
			memos = append(memos, memo)
		}
	}
	return memos, nil
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
