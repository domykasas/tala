package fileops

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir creates a temporary directory for testing
func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "tala-test-")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

// cleanupTestDir removes the temporary directory
func cleanupTestDir(t *testing.T, dir string) {
	if err := os.RemoveAll(dir); err != nil {
		t.Errorf("Failed to cleanup temp dir: %v", err)
	}
}

func TestCreateFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	tests := []struct {
		name     string
		filename string
		content  string
		wantErr  bool
	}{
		{
			name:     "create simple file",
			filename: "test.txt",
			content:  "Hello World",
			wantErr:  false,
		},
		{
			name:     "create file with empty content",
			filename: "empty.txt",
			content:  "",
			wantErr:  false,
		},
		{
			name:     "create file with empty filename",
			filename: "",
			content:  "content",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateFile(tt.filename, tt.content)
			if result.Success == tt.wantErr {
				t.Errorf("CreateFile() success = %v, wantErr %v", result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				// Verify file exists and has correct content
				content, err := os.ReadFile(tt.filename)
				if err != nil {
					t.Errorf("Failed to read created file: %v", err)
				}
				if string(content) != tt.content {
					t.Errorf("File content = %v, want %v", string(content), tt.content)
				}
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create test file
	testContent := "Test file content"
	err := os.WriteFile("test.txt", []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "read existing file",
			filename: "test.txt",
			wantErr:  false,
		},
		{
			name:     "read non-existent file",
			filename: "nonexistent.txt",
			wantErr:  true,
		},
		{
			name:     "read with empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReadFile(tt.filename)
			if result.Success == tt.wantErr {
				t.Errorf("ReadFile() success = %v, wantErr %v", result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				if !contains(result.Message, testContent) {
					t.Errorf("ReadFile() message should contain file content")
				}
			}
		})
	}
}

func TestDeleteFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create test file
	err := os.WriteFile("test.txt", []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		wantErr  bool
	}{
		{
			name:     "delete existing file",
			filename: "test.txt",
			wantErr:  false,
		},
		{
			name:     "delete non-existent file",
			filename: "nonexistent.txt",
			wantErr:  true,
		},
		{
			name:     "delete with empty filename",
			filename: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DeleteFile(tt.filename)
			if (result.Error != nil || !result.Success) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, success = %v, wantErr %v", result.Error, result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				// Verify file no longer exists
				if _, err := os.Stat(tt.filename); !os.IsNotExist(err) {
					t.Errorf("File still exists after deletion")
				}
			}
		})
	}
}

func TestCreateDirectory(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	tests := []struct {
		name    string
		dirname string
		wantErr bool
	}{
		{
			name:    "create simple directory",
			dirname: "testdir",
			wantErr: false,
		},
		{
			name:    "create nested directory",
			dirname: "parent/child",
			wantErr: false,
		},
		{
			name:    "create with empty dirname",
			dirname: "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CreateDirectory(tt.dirname)
			if result.Success == tt.wantErr {
				t.Errorf("CreateDirectory() success = %v, wantErr %v", result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				// Verify directory exists
				if _, err := os.Stat(tt.dirname); os.IsNotExist(err) {
					t.Errorf("Directory was not created")
				}
			}
		})
	}
}

func TestListDirectory(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create test files and directories
	os.WriteFile("file1.txt", []byte("test"), 0644)
	os.WriteFile("file2.txt", []byte("test"), 0644)
	os.Mkdir("dir1", 0755)

	result := ListDirectory("")
	if result.Error != nil {
		t.Errorf("ListDirectory() error = %v", result.Error)
		return
	}

	if !result.Success {
		t.Errorf("ListDirectory() success = false")
		return
	}

	// Check that the listing contains our test files
	if !contains(result.Message, "file1.txt") {
		t.Errorf("ListDirectory() should contain file1.txt")
	}
	if !contains(result.Message, "file2.txt") {
		t.Errorf("ListDirectory() should contain file2.txt")
	}
	if !contains(result.Message, "dir1") {
		t.Errorf("ListDirectory() should contain dir1")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Change to temp directory
	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	// Create source file
	sourceContent := "Source file content"
	err := os.WriteFile("source.txt", []byte(sourceContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	tests := []struct {
		name    string
		src     string
		dst     string
		wantErr bool
	}{
		{
			name:    "copy existing file",
			src:     "source.txt",
			dst:     "destination.txt",
			wantErr: false,
		},
		{
			name:    "copy non-existent file",
			src:     "nonexistent.txt",
			dst:     "destination.txt",
			wantErr: true,
		},
		{
			name:    "copy with empty source",
			src:     "",
			dst:     "destination.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CopyFile(tt.src, tt.dst)
			if (result.Error != nil || !result.Success) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, success = %v, wantErr %v", result.Error, result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				// Verify destination file exists and has correct content
				content, err := os.ReadFile(tt.dst)
				if err != nil {
					t.Errorf("Failed to read destination file: %v", err)
				}
				if string(content) != sourceContent {
					t.Errorf("Destination file content = %v, want %v", string(content), sourceContent)
				}
			}
		})
	}
}

func TestGetWorkingDirectory(t *testing.T) {
	result := GetWorkingDirectory()
	if result.Error != nil {
		t.Errorf("GetWorkingDirectory() error = %v", result.Error)
		return
	}
	if !result.Success {
		t.Errorf("GetWorkingDirectory() success = false")
		return
	}
	if !contains(result.Message, "Current working directory:") {
		t.Errorf("GetWorkingDirectory() should contain working directory message")
	}
}

func TestChangeDirectory(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer cleanupTestDir(t, tmpDir)

	// Create test directory
	testDir := filepath.Join(tmpDir, "testdir")
	err := os.Mkdir(testDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	// Store original directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{
			name:    "change to existing directory",
			path:    testDir,
			wantErr: false,
		},
		{
			name:    "change to non-existent directory",
			path:    "/nonexistent/path",
			wantErr: true,
		},
		{
			name:    "change with empty path",
			path:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ChangeDirectory(tt.path)
			if (result.Error != nil || !result.Success) != tt.wantErr {
				t.Errorf("ChangeDirectory() error = %v, success = %v, wantErr %v", result.Error, result.Success, tt.wantErr)
				return
			}
			if !tt.wantErr && result.Success {
				// Verify we're in the correct directory
				currentDir, err := os.Getwd()
				if err != nil {
					t.Errorf("Failed to get current directory: %v", err)
				}
				expectedDir, _ := filepath.Abs(tt.path)
				if currentDir != expectedDir {
					t.Errorf("Current directory = %v, want %v", currentDir, expectedDir)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > len(substr) && (s[:len(substr)] == substr || 
		s[len(s)-len(substr):] == substr || 
		containsHelper(s, substr))))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}