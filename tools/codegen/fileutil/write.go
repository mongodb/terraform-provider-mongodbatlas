package fileutil

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteFile ensures the parent directory exists before writing the file.
func WriteFile(filePath string, content []byte) error {
	// read/write/execute for owner, and read/execute for group and others
	const filePermission = 0o755

	// Create directories if they don't exist
	dir := filepath.Dir(filePath)
	dirPermission := os.FileMode(filePermission)
	if err := os.MkdirAll(dir, dirPermission); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write content to file (will override content if file exists)
	if err := os.WriteFile(filePath, content, filePermission); err != nil {
		return fmt.Errorf("failed to write to file %s: %w", filePath, err)
	}
	return nil
}
