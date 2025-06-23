package db

import (
	"encoding/json"
	"fmt"
	"memo/db/model"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

type FileTimestamps struct {
	CreatedAt time.Time
	UpdatedAt time.Time
}

func getFileTimestamps(filePath string) (FileTimestamps, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return FileTimestamps{}, err
	}

	sys := stat.Sys()
	statT, ok := sys.(*syscall.Stat_t)
	if !ok {
		return FileTimestamps{}, fmt.Errorf("failed to cast to *syscall.Stat_t")
	}

	return FileTimestamps{
		CreatedAt: time.Unix(statT.Birthtimespec.Sec, statT.Birthtimespec.Nsec),
		UpdatedAt: time.Unix(statT.Mtimespec.Sec, statT.Mtimespec.Nsec),
	}, nil
}

func changeFileType(fileName string) model.FileType {
	ext := filepath.Ext(fileName)
	fileType := model.FileType(strings.TrimPrefix(ext, "."))

	return fileType
}

func getTitle(fileName string) string {
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	title := strings.Split(baseName, "_")[0]

	return title
}

func generateContent(fileType model.FileType, content string) string {
	switch fileType {
	case model.FileTypeJson:
		return checkJsonContent(content)
	default:
		return content
	}
}

func checkJsonContent(content string) string {
	var jsonContent map[string]interface{}
	err := json.Unmarshal([]byte(content), &jsonContent)
	if err != nil {
		return content
	}

	return jsonContent["content"].(string)
}
