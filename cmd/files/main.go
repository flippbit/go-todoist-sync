package files

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func FindFilesWithExtension(dir, ext string) ([]FileDetails, error) {
	var files []FileDetails

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == "."+ext {
			files = append(files, FileDetails{
				Name: strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())),
				Path: path,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func ParseFiles(files []FileDetails) []ParsedDailyNote {
	var notes []ParsedDailyNote

	for _, file := range files {
		date, err := time.Parse("2006-01-02", file.Name)
		if err == nil {
			notes = append(notes, ParsedDailyNote{
				Path:      file.Path,
				CreatedAt: date,
			})
		}
	}

	return notes
}

func CreateFile(path string, content string) error {
	byteContent := []byte(content)

	err := os.WriteFile(path, byteContent, 0644)
	if err != nil {
		log.Printf("Error creating/writing file: %v", err)
		return err
	}

	log.Printf("File created successfully: %s", path)
	return nil
}
