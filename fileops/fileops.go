package fileops

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// FileOperation represents a file system operation result
type FileOperation struct {
	Success bool
	Message string
	Error   error
}

// ListDirectory lists files and directories in the current working directory
func ListDirectory(path string) *FileOperation {
	if path == "" {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return &FileOperation{
				Success: false,
				Error:   err,
				Message: fmt.Sprintf("Failed to get current directory: %v", err),
			}
		}
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to read directory '%s': %v", path, err),
		}
	}

	var result strings.Builder
	result.WriteString(fmt.Sprintf("Contents of %s:\n", path))
	
	for _, entry := range entries {
		prefix := "📄"
		if entry.IsDir() {
			prefix = "📁"
		}
		result.WriteString(fmt.Sprintf("%s %s\n", prefix, entry.Name()))
	}

	return &FileOperation{
		Success: true,
		Message: result.String(),
	}
}

// CreateFile creates a new file with optional content
func CreateFile(filename string, content string) *FileOperation {
	if filename == "" {
		return &FileOperation{
			Success: false,
			Message: "Filename cannot be empty",
		}
	}

	// Create file with content
	err := os.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to create file '%s': %v", filename, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Created file '%s'", filename),
	}
}

// ReadFile reads and returns the content of a file
func ReadFile(filename string) *FileOperation {
	if filename == "" {
		return &FileOperation{
			Success: false,
			Message: "Filename cannot be empty",
		}
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to read file '%s': %v", filename, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Content of '%s':\n%s", filename, string(content)),
	}
}

// UpdateFile updates an existing file with new content
func UpdateFile(filename string, content string) *FileOperation {
	if filename == "" {
		return &FileOperation{
			Success: false,
			Message: "Filename cannot be empty",
		}
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("File '%s' does not exist", filename),
		}
	}

	err := os.WriteFile(filename, []byte(content), 0600)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to update file '%s': %v", filename, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Updated file '%s'", filename),
	}
}

// DeleteFile removes a file
func DeleteFile(filename string) *FileOperation {
	if filename == "" {
		return &FileOperation{
			Success: false,
			Message: "Filename cannot be empty",
		}
	}

	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("File '%s' does not exist", filename),
		}
	}

	err := os.Remove(filename)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to delete file '%s': %v", filename, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Deleted file '%s'", filename),
	}
}

// CreateDirectory creates a new directory
func CreateDirectory(dirname string) *FileOperation {
	if dirname == "" {
		return &FileOperation{
			Success: false,
			Message: "Directory name cannot be empty",
		}
	}

	err := os.MkdirAll(dirname, 0750)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to create directory '%s': %v", dirname, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Created directory '%s'", dirname),
	}
}

// DeleteDirectory removes a directory and its contents
func DeleteDirectory(dirname string) *FileOperation {
	if dirname == "" {
		return &FileOperation{
			Success: false,
			Message: "Directory name cannot be empty",
		}
	}

	// Check if directory exists
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("Directory '%s' does not exist", dirname),
		}
	}

	err := os.RemoveAll(dirname)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to delete directory '%s': %v", dirname, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Deleted directory '%s'", dirname),
	}
}

// CopyFile copies a file from source to destination
func CopyFile(src, dst string) *FileOperation {
	if src == "" || dst == "" {
		return &FileOperation{
			Success: false,
			Message: "Source and destination filenames cannot be empty",
		}
	}

	// Check if source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("Source file '%s' does not exist", src),
		}
	}

	sourceFile, err := os.Open(src)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to open source file '%s': %v", src, err),
		}
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to create destination file '%s': %v", dst, err),
		}
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to copy file: %v", err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Copied '%s' to '%s'", src, dst),
	}
}

// MoveFile moves a file from source to destination
func MoveFile(src, dst string) *FileOperation {
	if src == "" || dst == "" {
		return &FileOperation{
			Success: false,
			Message: "Source and destination filenames cannot be empty",
		}
	}

	// Check if source file exists
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("Source file '%s' does not exist", src),
		}
	}

	err := os.Rename(src, dst)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to move file from '%s' to '%s': %v", src, dst, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Moved '%s' to '%s'", src, dst),
	}
}

// GetWorkingDirectory returns the current working directory
func GetWorkingDirectory() *FileOperation {
	wd, err := os.Getwd()
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to get working directory: %v", err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Current working directory: %s", wd),
	}
}

// ChangeDirectory changes the current working directory
func ChangeDirectory(path string) *FileOperation {
	if path == "" {
		return &FileOperation{
			Success: false,
			Message: "Path cannot be empty",
		}
	}

	// Convert relative path to absolute
	absPath, err := filepath.Abs(path)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to resolve path '%s': %v", path, err),
		}
	}

	// Check if directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return &FileOperation{
			Success: false,
			Message: fmt.Sprintf("Directory '%s' does not exist", absPath),
		}
	}

	err = os.Chdir(absPath)
	if err != nil {
		return &FileOperation{
			Success: false,
			Error:   err,
			Message: fmt.Sprintf("Failed to change directory to '%s': %v", absPath, err),
		}
	}

	return &FileOperation{
		Success: true,
		Message: fmt.Sprintf("Changed directory to '%s'", absPath),
	}
}