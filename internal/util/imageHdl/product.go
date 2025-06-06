package imageHdl

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func GenerateFilename(filename string) string {
    return strings.ReplaceAll(filename, " ", "_") + ".jpg"
}

func EnsureUploadDirectoryExists() (string, error) {
	rootDir, err := os.Getwd()
    if err != nil {
        log.Fatalf("Error getting current working directory: %v", err)
    }

	uploadPath := filepath.Join(rootDir, "internal", "assets", "products")
	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}
	return uploadPath, nil
}

func SaveImage(src io.Reader, destinationPath string) error {
	dst, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create image file on server: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to save image file: %w", err)
	}
	return nil
}

func DeleteImage(imagePath string) error {
	err := os.Remove(imagePath)
	if err != nil {
		return fmt.Errorf("failed to delete image file: %w", err)
	}
	return nil
}