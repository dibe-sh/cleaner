package main

import (
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
)

func TestShouldRemoveDir(t *testing.T) {
	dirsToRemove := []string{"node_modules", "dist", "build"}

	tests := []struct {
		dirName string
		want    bool
	}{
		{"node_modules", true},
		{"build", true},
		{"src", false},
	}

	for _, tt := range tests {
		got := shouldRemoveDir(tt.dirName, dirsToRemove)
		if got != tt.want {
			t.Errorf("shouldRemoveDir(%q) = %v; want %v", tt.dirName, got, tt.want)
		}
	}
}

func TestShouldRemoveFile(t *testing.T) {
	tests := []struct {
		fileName   string
		extensions []string
		matchRegex bool
		want       bool
	}{
		{"program.exe", []string{".exe", ".dll", ".tmp"}, false, true},
		{"tempfile.tmp", []string{".exe", ".dll", ".tmp"}, false, true},
		{"document.txt", []string{".exe", ".dll", ".tmp"}, false, false},
		{"tempfile.log", []string{`temp.*\.log`}, true, true},
		{"logfile.log", []string{`temp.*\.log`}, true, false},
		{"debug_info.txt", []string{`debug_.*\.txt`}, true, true},
		{"info.txt", []string{`debug_.*\.txt`}, true, false},
	}

	for _, tt := range tests {
		got := shouldRemoveFile(tt.fileName, tt.extensions, tt.matchRegex)
		if got != tt.want {
			t.Errorf("shouldRemoveFile(%q, %v, %v) = %v; want %v", tt.fileName, tt.extensions, tt.matchRegex, got, tt.want)
		}
	}
}

func TestIsExcluded(t *testing.T) {
	excludeList := []string{".git", ".svn", ".env"}

	tests := []struct {
		name string
		want bool
	}{
		{".git", true},
		{".env", true},
		{"README.md", false},
	}

	for _, tt := range tests {
		got := isExcluded(tt.name, excludeList)
		if got != tt.want {
			t.Errorf("isExcluded(%q) = %v; want %v", tt.name, got, tt.want)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	config := getDefaultConfig()

	if len(config.DirectoriesToRemove) == 0 {
		t.Error("Expected default DirectoriesToRemove to be non-empty")
	}
	if len(config.FileExtensionsToRemove) == 0 {
		t.Error("Expected default FileExtensionsToRemove to be non-empty")
	}
}

func TestCalculateCleanableSizeDryRun(t *testing.T) {
	rootDir := "test_dir"

	// Remove existing test directory if it exists to ensure a clean setup
	if _, err := os.Stat(rootDir); err == nil {
		err := os.RemoveAll(rootDir)
		if err != nil {
			t.Fatalf("Failed to remove existing test directory: %v", err)
		}
	}

	// Create a fresh test directory
	err := os.Mkdir(rootDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create test root directory: %v", err)
	}
	defer os.RemoveAll(rootDir) // Clean up after the test

	// List of paths to create and check for size
	filesToRemove := []string{
		filepath.Join(rootDir, "node_modules", "module.js"),
		filepath.Join(rootDir, "dist", "bundle.js"),
		filepath.Join(rootDir, "build", "output.o"),
		filepath.Join(rootDir, "tempfile.log"),
	}

	// Create files and directories for the test
	for _, path := range filesToRemove {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", path, err)
		}
		err = os.WriteFile(path, []byte("sample data"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", path, err)
		}
	}

	config := &Config{
		DirectoriesToRemove:    []string{"node_modules", "dist", "build"},
		FileExtensionsToRemove: []string{".log"},
		ExcludeDirectories:     []string{".git"},
		ExcludeFiles:           []string{},
		MatchRegex:             true,
	}

	var wg sync.WaitGroup
	var paths []string
	var totalSize int64
	var totalSizeMutex sync.Mutex
	logMutex := &sync.Mutex{}
	consoleMutex := &sync.Mutex{}
	semaphore := make(chan struct{}, runtime.NumCPU())

	wg.Add(1)
	go walkDir(rootDir, config, true, true, "text", &wg, semaphore, nil, logMutex, consoleMutex, &totalSize, &totalSizeMutex, &paths)
	wg.Wait()

	// Expected paths to delete
	expectedPaths := []string{
		filepath.Join(rootDir, "node_modules"),
		filepath.Join(rootDir, "dist"),
		filepath.Join(rootDir, "build"),
		filepath.Join(rootDir, "tempfile.log"),
	}

	// Compare actual paths with expected paths
	if len(paths) != len(expectedPaths) {
		t.Fatalf("Expected %d paths, but got %d paths", len(expectedPaths), len(paths))
	}

	for _, expectedPath := range expectedPaths {
		found := false
		for _, path := range paths {
			if path == expectedPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected path %s was not found in the paths to delete", expectedPath)
		}
	}
}
