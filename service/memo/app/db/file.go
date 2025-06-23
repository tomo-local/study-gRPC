package db

import (
	"fmt"
	"os"
	"strings"

	config "memo/config/server"
	"memo/db/model"
)

// FileService defines the interface for file-based memo operations.
type FileService interface {
	CreateFile(memo *model.Memo) (*model.Memo, error)
	GetFile(id string) (*model.Memo, error)
	UpdateFile(memo *model.Memo) (*model.Memo, error)
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

	filePath := memo.GetFilePath(f.folderPath)
	content := generateContent(memo.FileType, memo.Content)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return nil, fmt.Errorf("failed to write file: %w", err)
	}

	timestamps, err := getFileTimestamps(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get file timestamps: %w", err)
	}

	return &model.Memo{
		ID:        memo.ID,
		Title:     memo.Title,
		FileType:  memo.FileType,
		Content:   content,
		CreatedAt: timestamps.CreatedAt,
		UpdatedAt: timestamps.UpdatedAt,
	}, nil
}

// GetFile retrieves a memo file by its ID.
func (f *fileService) GetFile(id string) (*model.Memo, error) {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.Contains(entry.Name(), id) {
			fileType := changeFileType(entry.Name())
			title := getTitle(entry.Name())
			memo := &model.Memo{
				ID:       id,
				FileType: fileType,
				Title:    title,
			}

			filePath := memo.GetFilePath(f.folderPath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}

			generatedContent := generateContent(fileType, string(content))
			timestamps, err := getFileTimestamps(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to get file timestamps: %w", err)
			}

			return &model.Memo{
				ID:        id,
				Title:     title,
				FileType:  fileType,
				Content:   generatedContent,
				CreatedAt: timestamps.CreatedAt,
				UpdatedAt: timestamps.UpdatedAt,
			}, nil
		}
	}

	return nil, fmt.Errorf("file with id %s not found", id)
}

// UpdateFile updates the file for the given memo.
// It currently reuses CreateFile, which overwrites the existing file.
func (f *fileService) UpdateFile(targetMemo *model.Memo) (*model.Memo, error) {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.Contains(entry.Name(), targetMemo.ID) {
			fileType := changeFileType(entry.Name())
			title := getTitle(entry.Name())
			content := generateContent(fileType, targetMemo.Content)

			if content == "" {
				return nil, fmt.Errorf("content is empty")
			}

			updatedMemo := &model.Memo{
				ID:       targetMemo.ID,
				FileType: fileType,
				Title:    title,
				Content:  content,
			}

			filePath := updatedMemo.GetFilePath(f.folderPath)

			if err := os.WriteFile(filePath, []byte(updatedMemo.Content), 0644); err != nil {
				return nil, fmt.Errorf("failed to write file: %w", err)
			}

			timestamps, err := getFileTimestamps(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to get file timestamps: %w", err)
			}

			return &model.Memo{
				ID:        updatedMemo.ID,
				Title:     title,
				FileType:  fileType,
				Content:   content,
				CreatedAt: timestamps.CreatedAt,
				UpdatedAt: timestamps.UpdatedAt,
			}, nil
		}
	}

	return nil, fmt.Errorf("file with id %s not found", targetMemo.ID)
}
