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
	ListFiles() ([]*model.Memo, error)
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
			fileName := formatFileName(entry.Name())
			memo := &model.Memo{
				ID:       id,
				FileType: fileName.FileType,
				Title:    fileName.Title,
			}

			filePath := memo.GetFilePath(f.folderPath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}

			generatedContent := generateContent(fileName.FileType, string(content))
			timestamps, err := getFileTimestamps(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to get file timestamps: %w", err)
			}

			return &model.Memo{
				ID:        id,
				Title:     fileName.Title,
				FileType:  fileName.FileType,
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
			fileName := formatFileName(entry.Name())
			content := generateContent(fileName.FileType, targetMemo.Content)

			if content == "" {
				return nil, fmt.Errorf("content is empty")
			}

			updatedMemo := &model.Memo{
				ID:       targetMemo.ID,
				FileType: fileName.FileType,
				Title:    fileName.Title,
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
				Title:     fileName.Title,
				FileType:  fileName.FileType,
				Content:   content,
				CreatedAt: timestamps.CreatedAt,
				UpdatedAt: timestamps.UpdatedAt,
			}, nil
		}
	}

	return nil, fmt.Errorf("file with id %s not found", targetMemo.ID)
}

// ListFiles lists all memo files in the folder.
func (f *fileService) ListFiles() ([]*model.Memo, error) {
	dirEntries, err := os.ReadDir(f.folderPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	files := make([]*model.Memo, 0, len(dirEntries))
	for _, entry := range dirEntries {
		if !entry.IsDir() {
			fileName := formatFileName(entry.Name())

			memo := &model.Memo{
				ID:       fileName.ID,
				FileType: fileName.FileType,
				Title:    fileName.Title,
			}

			filePath := memo.GetFilePath(f.folderPath)
			content, err := os.ReadFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to read file: %w", err)
			}

			generatedContent := generateContent(fileName.FileType, string(content))

			files = append(files, &model.Memo{
				ID:       fileName.ID,
				Title:    fileName.Title,
				FileType: fileName.FileType,
				Content:  generatedContent,
			})
		}
	}

	return files, nil
}
